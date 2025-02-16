package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/mocks/mservice"
)

const SET_COOKIE_URL = "/set"
const OK = "ok"

type TestServer struct {
	server      *Server
	mockService *mservice.MockServices
	engine      *gin.Engine
	request     *http.Request
}

func setupEngine() *gin.Engine {
	r := gin.New()
	r.Use(sessions.Sessions("exorcist", cookie.NewStore([]byte("cookieSecret"))))

	r.GET(SET_COOKIE_URL, func(c *gin.Context) {
		session := sessions.Default(c)

		var cookieBody struct {
			value string
		}

		_ = c.BindJSON(&cookieBody)

		session.Set(userKey, cookieBody.value)
		_ = session.Save()
		c.String(http.StatusOK, OK)
	})
	return r
}

func (s *TestServer) withGetEndpoint(f gin.HandlerFunc, extraPathParams string) *TestServer {
	s.engine.GET(fmt.Sprintf("/%v", extraPathParams), f)
	return s
}

func (s *TestServer) withPostEndpoint(f gin.HandlerFunc) *TestServer {
	s.engine.POST("/", f)
	return s
}

func (s *TestServer) withGetRequest(body io.Reader, params string) *TestServer {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", params), body)
	s.request = req
	return s
}

func (s *TestServer) withPostRequest(body io.Reader) *TestServer {
	req, _ := http.NewRequest("POST", "/", body)
	s.request = req
	return s
}

func (s *TestServer) exec() *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.engine.ServeHTTP(rr, s.request)
	return rr
}

func setupServer() *TestServer {
	svc, mSvc := mservice.SetupMockService()
	server := &Server{logger: logger.New(&environment.EnvironmentVariables{LogLevel: "none"}), service: svc}
	engine := setupEngine()
	return &TestServer{server: server, mockService: mSvc, engine: engine}
}

func setupCookies(req *http.Request, r *gin.Engine) {
	res := httptest.NewRecorder()
	cookieReq, _ := http.NewRequest("GET", SET_COOKIE_URL, body(`{"value": "val"}`))
	r.ServeHTTP(res, cookieReq)

	req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}

func body(body string, args ...any) *bytes.Reader {
	message := body
	if args != nil {
		message = fmt.Sprintf(body, args...)
	}
	return bytes.NewReader([]byte(message))
}
