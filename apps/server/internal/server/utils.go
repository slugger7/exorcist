package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/dto"
)

type ApiError = string

func createError(e ApiError) map[string]any {
	return gin.H{"error": e}
}

func (s *server) getUserId(c *gin.Context) (*uuid.UUID, error) {
	session := sessions.Default(c)
	userString := session.Get(userKey).(string)
	userId, err := uuid.Parse(userString)
	if err != nil {
		s.logger.Errorf("could not parse userId from string: %v\n%v", userString, err.Error())
		return nil, err
	}
	return &userId, err
}

var TRUE bool = true
var FALSE bool = false

var MEDIA_SEARCH_DEFAULT dto.MediaSearchDTO = dto.MediaSearchDTO{
	Deleted: &FALSE,
	Exists:  &TRUE,
	PageRequestDTO: dto.PageRequestDTO{
		Limit: 50,
	},
}
