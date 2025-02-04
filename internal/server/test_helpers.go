package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type mockService struct {
	userService    userService.IUserService
	libraryService libraryService.ILibraryService
}

type mockUserService struct {
	returningModel *model.User
	returningError error
}

type mockLibraryService struct {
	returningModel *model.Library
	returningError error
}

func (ms mockService) UserService() userService.IUserService {
	return ms.userService
}

func (mus mockUserService) CreateUser(username, password string) (*model.User, error) {
	return mus.returningModel, mus.returningError
}

func (mus mockUserService) ValidateUser(username, password string) (*model.User, error) {
	return mus.returningModel, mus.returningError
}

func (ms mockService) LibraryService() libraryService.ILibraryService {
	return ms.libraryService
}

func (ls mockLibraryService) CreateLibrary(actual model.Library) (*model.Library, error) {
	return ls.returningModel, ls.returningError
}

const SET_COOKIE_URL = "/set"
const OK = "ok"

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

func setupCookies(req *http.Request, r *gin.Engine) {
	res := httptest.NewRecorder()
	cookieReq, _ := http.NewRequest("GET", SET_COOKIE_URL, body(`{"value": "val"}`))
	r.ServeHTTP(res, cookieReq)

	req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}

func body(body string) *bytes.Reader {
	return bytes.NewReader([]byte(body))
}
