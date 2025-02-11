package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_CreateLibrary_InvalidBody(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	r.POST("/", s.CreateLibrary)

	req, err := http.NewRequest("POST", "/", body(`{invalid json}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"could not read body of request"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_CreateLibrary_NoNameSpecified_ShouldThrowError(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	r.POST("/", s.CreateLibrary)

	req, err := http.NewRequest("POST", "/", body(`{"name": ""}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"no value for name"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_CreateLibrary_ErrorByService(t *testing.T) {
	r := setupEngine()
	s := setupServer()
	svc, mSvc := setupService()
	s.service = svc

	expectedErrorMessage := "expected error message"
	mSvc.libraryService.mockErrors[0] = errors.New(expectedErrorMessage)
	r.POST("/", s.CreateLibrary)

	expectedName := "expectedLibraryName"
	req, err := http.NewRequest("POST", "/", body(fmt.Sprintf(`{"name":"%v"}`, expectedName)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"could not create new library"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_CreateLibrary_Success(t *testing.T) {
	r := setupEngine()
	s := setupServer()
	svc, mSvc := setupService()
	s.service = svc

	expectedId, _ := uuid.NewRandom()
	expectedLibraryName := "some expected library name"
	mSvc.libraryService.mockModel[0] = &model.Library{
		ID:   expectedId,
		Name: expectedLibraryName,
	}

	r.POST("/", s.CreateLibrary)

	req, err := http.NewRequest("POST", "/", body(fmt.Sprintf(`{"name":"%v"}`, expectedLibraryName)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusCreated
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"id":"%v"}`, expectedId.String())
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_GetLibraries_ServiceReturnsError(t *testing.T) {
	r := setupEngine()
	ms, services := setupService()
	s := setupServer()
	s.service = ms
	expectedError := errors.New("expected error")
	services.libraryService.mockErrors[0] = expectedError

	r.GET("/", s.GetLibraries)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusInternalServerError
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"could not fetch libraries"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}
