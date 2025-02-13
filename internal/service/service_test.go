package service

import (
	"testing"

	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/mocks/mservice"
)

var service = &Service{
	env:                nil,
	logger:             logger.New(&environment.EnvironmentVariables{LogLevel: "none"}),
	userService:        mservice.MockUserService{},
	libraryService:     mservice.MockLibraryService{},
	libraryPathService: mservice.MockLibaryPathService{},
}

func Test_UserService(t *testing.T) {
	actualUserService := service.UserService()
	if actualUserService == nil {
		t.Error("user service was nil")
	}
}

func Test_LibraryService(t *testing.T) {
	actualLibraryService := service.LibraryService()
	if actualLibraryService == nil {
		t.Error("library service was nil")
	}
}

func Test_LibraryPathService(t *testing.T) {
	actualLibraryPathService := service.LibraryPathService()
	if actualLibraryPathService == nil {
		t.Errorf("library path service was nil")
	}
}
