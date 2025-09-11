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
	"github.com/slugger7/exorcist/internal/websockets"
)

type watcherService struct {
	logger    logger.Logger
	repo      repository.Repository
	wsService websockets.Websockets
	env       *environment.EnvironmentVariables
	ctx       context.Context
	wg        *sync.WaitGroup
}

type WatcherService interface {
	WithDirectoryWatcher()
}

var watcherServiceInstance *watcherService

func New(
	env environment.EnvironmentVariables,
	ctx context.Context,
	wg *sync.WaitGroup,
	repo repository.Repository,
	wsService websockets.Websockets,
) WatcherService {
	if watcherServiceInstance == nil {
		watcherServiceInstance = &watcherService{
			logger:    logger.New(&env),
			env:       &env,
			repo:      repo,
			wsService: wsService,
			ctx:       ctx,
			wg:        wg,
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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.logger.Errorf("could not set up file watcher: %v", err.Error())
		return
	}
	defer watcher.Close()

	libPaths, err := s.repo.LibraryPath().GetAll()
	if err != nil {
		s.logger.Errorf("could not get all library paths to add to watcher: %v", err.Error())
	}

	for _, lp := range libPaths {
		watcher.Add(lp.Path)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					ext := filepath.Ext(event.Name)
					if slices.Contains(constants.VideoExtensions[:], ext) {
						libPath := findLibPathByFilePath(event.Name, libPaths)

						if libPath == nil {
							continue
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

						continue
					}

					if slices.Contains(constants.ImageExtensions[:], ext) {
						// TODO: handle images
					}
				}

				if event.Has(fsnotify.Remove) {
					m, err := s.repo.Media().GetByPath(event.Name)
					if err != nil {
						s.logger.Errorf("remove event triggered but could not find media by path: %v", event.Name)
						continue
					}

					if m == nil {
						continue
					}

					m.Exists = false

					if err := s.repo.Media().UpdateExists(*m); err != nil {
						s.logger.Errorf("could not mark media as removed in file watcher: %v", err.Error())
						continue
					}

					s.wsService.MediaDelete(dto.MediaOverviewDTO{Id: m.ID, Deleted: true})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				s.logger.Errorf("error in fsnotify watcher: %v", err.Error())
			case <-s.ctx.Done():
				return
			}
		}
	}()
}
