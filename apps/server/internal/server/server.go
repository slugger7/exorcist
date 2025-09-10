package server

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

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

type server struct {
	env       *environment.EnvironmentVariables
	repo      repository.Repository
	service   service.Service
	logger    logger.Logger
	jobCh     chan bool
	wsService websockets.Websockets
}

func (s *server) withJobRunner(ctx context.Context, wg *sync.WaitGroup, ws websockets.Websockets) *server {
	ch := job.New(s.env, s.service, s.logger, ctx, wg, ws)
	s.jobCh = ch

	ch <- true // start if any jobs exist

	return s
}

func findLibPathByFilePath(p string, libPaths []model.LibraryPath) *model.LibraryPath {
	for _, l := range libPaths {
		if strings.HasPrefix(p, l.Path) {
			return &l
		}
	}
	return nil
}

func (s *server) withDirectoryWatcher(ctx context.Context, wg *sync.WaitGroup) *server {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.logger.Errorf("could not set up file watcher: %v", err.Error())
		return s
	}
	defer watcher.Close()

	libPaths, err := s.repo.LibraryPath().GetAll()
	if err != nil {
		s.logger.Errorf("could not get all library paths to add to watcher: %v", err.Error())
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
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
			case <-ctx.Done():
				return
			}
		}
	}()

	return s
}

func New(env *environment.EnvironmentVariables, wg *sync.WaitGroup) *http.Server {
	lg := logger.New(env)
	shutdownCtx, cancel := context.WithCancel(context.Background())

	repo := repository.New(env, shutdownCtx)

	newServer := &server{
		repo:      repo,
		env:       env,
		logger:    lg,
		wsService: websockets.New(env),
	}

	err := newServer.repo.Job().CancelInprogress()
	if err != nil {
		lg.Errorf("clearing in progress jobs on startup: %v", err.Error())
	}

	if env.JobRunner {
		newServer.withJobRunner(shutdownCtx, wg, newServer.wsService)
	}
	newServer.service = service.New(repo, env, newServer.jobCh, shutdownCtx)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server.RegisterOnShutdown(func() {
		newServer.logger.Info("Shutting down server. Stopping job runner.")
		cancel()
		close(newServer.jobCh)

		newServer.logger.Debug("Cancelled and closed")

	})

	return server
}
