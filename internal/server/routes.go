package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/environment"
)

const (
	root         string = "/"
	userRoute    string = "/users"
	libraryRoute string = "/libraries"
	videoRoute   string = "/videos"
)

func (s *Server) RegisterRoutes() http.Handler {
	if s.env.AppEnv == environment.AppEnvEnum.Prod {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Register authentication routes
	s.withAuthLogin(&r.RouterGroup, root).
		withAuthLogout(&r.RouterGroup, root)

	authenticated := r.Group("/api")
	authenticated.Use(s.AuthRequired)
	// Register user controller routes
	s.withUserCreate(authenticated, userRoute).
		withUserUpdatePassword(authenticated, userRoute)

	// Register library controller routes
	s.withLibraryGet(authenticated, libraryRoute).
		withLibraryGetAction(authenticated, libraryRoute).
		withLibraryPost(authenticated, libraryRoute)

	s.WithLibraryPathRoutes(authenticated).
		WithJobRoutes(authenticated)

	r.GET("/health", s.HealthHandler)
	return r
}

func (s *Server) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.repo.Health())
}
