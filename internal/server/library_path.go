package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

const libraryPathRoute = "/library"

func (s *Server) RegisterLibraryPathRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST(libraryPathRoute, s.CreateLibraryPath)

	return r
}

func (s *Server) CreateLibraryPath(c *gin.Context) {
	var body struct {
		LibraryId uuid.UUID
		Path      string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body of request"})
		return
	}
	if body.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no path in body"})
		return
	}

	library, err := s.service.LibraryService().GetLibraryById(body.LibraryId)
	if err != nil {
		s.logger.Errorf("Error finding library: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch library by id"})
		return
	}
	if library == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no library found with that id"})
		return
	}

	libPath := &model.LibraryPath{LibraryID: library.ID, Path: body.Path}
	libPath, err = s.service.LibraryPathService().Create(libPath)
	if err != nil {
		s.logger.Errorf("Erorr creating library path: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "colud not create new library path"})
		return
	}

	c.JSON(http.StatusCreated, libPath)
}
