package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"go.uber.org/mock/gomock"
)

func Test_GetVideo_InvalidId(t *testing.T) {
	s := setupServer(t)

	id := "some invalid uuid"

	s.server.withVideoGetById(&s.engine.RouterGroup, "/")
	rr := s.withGetRequest(id).
		exec()

	assert.StatusCode(t, http.StatusBadRequest, rr.Code)
	assert.Body(t, errBody(ErrInvalidIdFormat), rr.Body.String())
}

func Test_GetVideo_ServiceReturnsError(t *testing.T) {
	s := setupServer(t).
		withVideoService()

	id, _ := uuid.NewRandom()

	s.mockVideoService.EXPECT().
		GetById(gomock.Eq(id)).
		DoAndReturn(func(uuid.UUID) (*model.Video, error) {
			return nil, fmt.Errorf("some error")
		}).
		Times(1)

	s.server.withVideoGetById(&s.engine.RouterGroup, "/")
	rr := s.withGetRequest(id.String()).
		exec()

	assert.StatusCode(t, http.StatusInternalServerError, rr.Code)
	assert.Body(t, errBody(ErrGetVideoService), rr.Body.String())
}

func Test_GetVideo_VideoServiceNil(t *testing.T) {
	s := setupServer(t).
		withVideoService()

	id, _ := uuid.NewRandom()

	s.mockVideoService.EXPECT().
		GetById(gomock.Eq(id)).
		DoAndReturn(func(uuid.UUID) (*model.Video, error) {
			return nil, nil
		}).
		Times(1)

	s.server.withVideoGetById(&s.engine.RouterGroup, "/")
	rr := s.withGetRequest(id.String()).
		exec()

	assert.StatusCode(t, http.StatusNotFound, rr.Code)
	assert.Body(t, errBody(ErrVideoNotFound), rr.Body.String())
}

func Test_GetVideo_Success(t *testing.T) {
	s := setupServer(t).
		withVideoService()

	id, _ := uuid.NewRandom()
	video := &model.Video{ID: id}

	s.mockVideoService.EXPECT().
		GetById(gomock.Eq(id)).
		DoAndReturn(func(uuid.UUID) (*model.Video, error) {
			return video, nil
		}).
		Times(1)

	s.server.withVideoGetById(&s.engine.RouterGroup, "/")
	rr := s.withGetRequest(id.String()).
		exec()

	assert.StatusCode(t, http.StatusMovedPermanently, rr.Code)
	// TODO: figure out how to test that a video is returned
}
