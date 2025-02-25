package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const imageRoute string = "/images"

func (s *Server) withImageGetById(r *gin.RouterGroup, route string) *Server {
	r.GET(fmt.Sprintf("%v/:id", route), s.GetImage)
	return s
}

func (s *Server) GetImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"messag": "not implemented"})
}
