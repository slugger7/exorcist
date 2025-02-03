package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_Login_IncorrectBody(t *testing.T) {
	r := setupEngine()
	s := &Server{}

	r.POST("/", s.Login)

	req, err := http.NewRequest("POST", "/", body(`{invalid json}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code returned\nexpected %v but got %v", http.StatusBadRequest, status)
	}
	expectedBody := `{"error":"could not read body of request"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Login_InvalidParametersInBody(t *testing.T) {
	r := setupEngine()
	s := &Server{}

	r.POST("/", s.Login)

	req, err := http.NewRequest("POST", "/", body(`{"username": " ", "password": " "}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code returned\nexpected %v but got %v", http.StatusBadRequest, status)
	}
	expectedBody := `{"error":"parameters can't be empty"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Login_NoUserFromValidateUser(t *testing.T) {
	r := setupEngine()
	s := &Server{}
	s.service = mockService{mockUserService{validateUserModel: nil}}

	r.POST("/", s.Login)

	req, err := http.NewRequest("POST", "/", body(`{"username": "admin", "password": "admin"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("wrong status code returned\nexpected %v but got %v", http.StatusBadRequest, status)
	}
	expectedBody := `{"error":"could not authenticate with credentials"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Login_Success(t *testing.T) {
	r := setupEngine()
	s := &Server{}
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("could not generate random uuid %v", err)
	}
	s.service = mockService{userService: mockUserService{
		validateUserModel: &model.User{Username: "admin", ID: id},
	}}

	r.POST("/", s.Login)

	req, err := http.NewRequest("POST", "/", body(`{"username": "admin", "password": "admin"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Got an error %v %v", status, err)
	}

	expectedBody := `{"message":"successfully authenticated user"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}

	cookie := strings.Trim(rr.Header().Get("Set-Cookie"), " ")
	if cookie == "" {
		t.Errorf("No header is being set for exorcist")
	}
	if !strings.Contains(cookie, "exorcist") {
		t.Errorf("cookie was not set up correctly: %v", cookie)
	}
}

func Test_Logout_InvalidSessionToken(t *testing.T) {
	r := setupEngine()
	s := &Server{}
	s.service = mockService{mockUserService{validateUserModel: nil}}

	r.GET("/", s.Logout)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code returned\nexpected %v but got %v", http.StatusBadRequest, status)
	}
	expectedBody := `{"error":"invalid session token"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Logout_Success(t *testing.T) {
	s := &Server{}
	r := setupEngine()
	s.service = mockService{mockUserService{validateUserModel: nil}}

	r.GET("/", s.Logout)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	setupCookies(req, r)

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code returned\nexpected %v but got %v", http.StatusOK, status)
	}
	expectedBody := `{"message":"successfully logged out"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}
