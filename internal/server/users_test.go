package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/models"
	"go.uber.org/mock/gomock"
)

func Test_Create_InvalidBody(t *testing.T) {
	s := setupServer(t)

	rr := s.withPostEndpoint(s.server.CreateUser).
		withPostRequest(body("{invalid body}")).
		exec()
	assert.StatusCode(t, http.StatusUnprocessableEntity, rr.Code)
	expectedBody := `{"error":"invalid character 'i' looking for beginning of object key string"}`
	assert.Body(t, expectedBody, rr.Body.String())
}

func Test_Create_ServiceReturnsError(t *testing.T) {
	s := setupServer(t).
		withUserService()

	u := models.CreateUserModel{
		Username: "someUsername",
		Password: "somePassword",
	}
	s.mockUserService.EXPECT().
		Create(gomock.Eq(u.Username), gomock.Eq(u.Password)).
		DoAndReturn(func(string, string) (*model.User, error) {
			return nil, fmt.Errorf("some error")
		}).
		Times(1)

	rr := s.withPostEndpoint(s.server.CreateUser).
		withPostRequest(bodyM(u)).
		exec()

	assert.StatusCode(t, http.StatusBadRequest, rr.Code)

	body, _ := json.Marshal(gin.H{"error": ErrCreateUser})
	assert.Body(t, string(body), rr.Body.String())
}

func Test_Create_Success(t *testing.T) {
	s := setupServer(t).
		withUserService()

	nu := &models.CreateUserModel{
		Username: "expectedUsername",
		Password: "somePassword",
	}

	m := &model.User{
		Username: nu.Username,
		Password: nu.Password,
	}

	s.mockUserService.EXPECT().
		Create(gomock.Eq(nu.Username), gomock.Eq(nu.Password)).
		DoAndReturn(func(string, string) (*model.User, error) {
			return m, nil
		}).
		Times(1)

	rr := s.withPostEndpoint(s.server.CreateUser).
		withPostRequest(bodyM(nu)).
		exec()
	body, _ := json.Marshal(m)
	assert.StatusCode(t, http.StatusCreated, rr.Code)
	assert.Body(t, string(body), rr.Body.String())
}

func Test_UpdatePassword_InvalidBody(t *testing.T) {
	s := setupServer(t).
		withAuth()

	rr := s.withAuthPutEndpoint(s.server.UpdatePassword, "").
		withAuthPutRequest(body("{invalid json body}"), "").
		withCookie(TestCookie{}).
		exec()

	expectedStatus := http.StatusUnprocessableEntity
	if rr.Code != expectedStatus {
		t.Errorf("Expected status: %v\nGot status: %v", expectedStatus, rr.Code)
	}
}

func Test_UpdatePassword_ServiceReturnsError(t *testing.T) {
	s := setupServer(t).
		withUserService().
		withAuth()

	rpm := models.ResetPasswordModel{
		OldPassword:    "good old boy",
		NewPassword:    "sparkly new",
		RepeatPassword: "sparkly new",
	}
	id, _ := uuid.NewRandom()

	s.mockUserService.EXPECT().
		UpdatePassword(gomock.Eq(id), gomock.Eq(rpm)).
		DoAndReturn(func(uuid.UUID, models.ResetPasswordModel) error {
			return fmt.Errorf("some error")
		}).
		Times(1)

	rr := s.withAuthPutEndpoint(s.server.UpdatePassword, "").
		withAuthPutRequest(bodyM(rpm), "").
		withCookie(TestCookie{Value: id}).
		exec()

	assert.StatusCode(t, http.StatusInternalServerError, rr.Code)
	assert.Body(t, fmt.Sprintf(`{"error":"%v"}`, ErrUpdatePassword), rr.Body.String())
}

func Test_UpdatePasswrod_ServiceSucceeds(t *testing.T) {
	s := setupServer(t).
		withUserService().
		withAuth()

	rpm := models.ResetPasswordModel{
		OldPassword:    "good old boy",
		NewPassword:    "sparkly new",
		RepeatPassword: "sparkly new",
	}

	id, _ := uuid.NewRandom()

	s.mockUserService.EXPECT().
		UpdatePassword(gomock.Eq(id), gomock.Eq(rpm)).
		DoAndReturn(func(uuid.UUID, models.ResetPasswordModel) error {
			return nil
		}).
		Times(1)

	rr := s.withAuthPutEndpoint(s.server.UpdatePassword, "").
		withAuthPutRequest(bodyM(rpm), "").
		withCookie(TestCookie{Value: id}).
		exec()

	assert.StatusCode(t, http.StatusOK, rr.Code)
	assert.Body(t, fmt.Sprintf(`{"message":"%v"}`, OkPasswordUpdate), rr.Body.String())
}
