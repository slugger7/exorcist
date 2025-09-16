package filewatcher

import (
	"context"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/slugger7/exorcist/internal/constants"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/job"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/media"
	"github.com/slugger7/exorcist/internal/repository"
	"github.com/slugger7/exorcist/internal/service"
	"github.com/slugger7/exorcist/internal/websockets"
)

type watcherService struct {
	logger    logger.Logger
	repo      repository.Repository
	wsService websockets.Websockets
	service   service.Service
	env       *environment.EnvironmentVariables
	ctx       context.Context
	wg        *sync.WaitGroup
	watcher   *fsnotify.Watcher
	libPaths  []model.LibraryPath
}

type WatcherService interface {
	WithDirectoryWatcher()
	Add(libPath model.LibraryPath)
	Close()
}

var watcherServiceInstance *watcherService

func New(
	env environment.EnvironmentVariables,
	ctx context.Context,
	wg *sync.WaitGroup,
	repo repository.Repository,
	wsService websockets.Websockets,
	service service.Service,
) WatcherService {
	if watcherServiceInstance == nil {
		logger := logger.New(&env)
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			logger.Errorf("could not set up file watcher: %v", err.Error())
			return nil
		}

		watcherServiceInstance = &watcherService{
			logger:    logger,
			env:       &env,
			repo:      repo,
			wsService: wsService,
			ctx:       ctx,
			wg:        wg,
			watcher:   watcher,
			service:   service,
		}
		watcherServiceInstance.logger.Info("created file watcher instance")

		libPaths, err := watcherServiceInstance.repo.LibraryPath().GetAll()
		if err != nil {
			logger.Errorf("could not get all library paths to add to watcher: %v", err.Error())
		}

		watcherServiceInstance.libPaths = libPaths

		for _, lp := range libPaths {
			logger.Infof("watching %v", lp.Path)
			watcher.Add(lp.Path)
		}
	}

	return watcherServiceInstance
}

func findLibPathByFilePath(p string, libPaths []model.LibraryPath) *model.LibraryPath {
	for _, l := range libPaths {
		if strings.HasPrefix(p, l.Path) {
			return &l
		}
	}
	return nil
}

func (s *watcherService) WithDirectoryWatcher() {
	s.logger.Info("starting directory watcher")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					return
				}

				s.logger.Debug("Got a file event in watched directory")

				if event.Has(fsnotify.Create) {
					ext := filepath.Ext(event.Name)
					if slices.Contains(constants.VideoExtensions[:], ext) {
						s.logger.Infof("new file created: %v", event.Name)

						libPath := findLibPathByFilePath(event.Name, s.libPaths)

						if libPath == nil {
							continue
						}

						m, err := s.repo.Media().GetByPath(event.Name)
						if err != nil {
							s.logger.Errorf("could not get media by path(%v): %v", event.Name, err.Error())
							continue
						}

						if m != nil {
							if !m.Deleted && m.Exists {
								continue
							}
						}

						f, err := media.GetFileInformation(event.Name)
						if err != nil {
							s.logger.Errorf("could not successfully get file information: %v", err.Error())
							continue
						}

						if err := job.CreateNewMedia(libPath, nil, *f, *s.env, s.repo, s.logger, s.wsService); err != nil {
							s.logger.Errorf("could not create new media from watcher: %v", err.Error())
							continue
						}

						s.service.Job().StartJobRunner()

						continue
					}

					if slices.Contains(constants.ImageExtensions[:], ext) {
						// TODO: handle images
					}
				}

				if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					m, err := s.repo.Media().GetByPath(event.Name)
					if err != nil {
						s.logger.Errorf("remove event triggered but could not find media by path: %v", event.Name)
						continue
					}

					if m == nil {
						continue
					}

					s.logger.Infof("file removed or renamed: %v", event.Name)

					m.Exists = false

					if err := s.repo.Media().UpdateExists(*m); err != nil {
						s.logger.Errorf("could not mark media as removed in file watcher: %v", err.Error())
						continue
					}

					s.logger.Infof("media has been marked as not exsisting any more: %v", event.Name)

					s.wsService.MediaDelete(dto.MediaOverviewDTO{Id: m.ID, Deleted: true})
				}
			case err, ok := <-s.watcher.Errors:
				if !ok {
					return
				}

				s.logger.Errorf("error in fsnotify watcher: %v", err.Error())
			case <-s.ctx.Done():
				s.logger.Info("shutting down file watcher service due to shutdown")
				return
			}
		}
	}()
}

func (s *watcherService) Add(libPath model.LibraryPath) {
	s.libPaths = append(s.libPaths, libPath)
	s.watcher.Add(libPath.Path)
}

func (s *watcherService) Close() {
	s.watcher.Close()
}
