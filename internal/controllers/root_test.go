package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/services"
	"div-dash/util/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setUpDb() (sqlmock.Sqlmock, func()) {
	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)

	return mock, func() {
		sdb.Close()
	}
}

func NewApi() (sqlmock.Sqlmock, func(), *gin.Engine) {
	mock, cleanup := setUpDb()

	testutil.SetupConfig()

	router := gin.Default()
	RegisterRoutes(router)
	return mock, cleanup, router
}

func PerformRequest(router *gin.Engine, method, path string) *httptest.ResponseRecorder {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(method, path, nil)

	router.ServeHTTP(w, req)
	return w
}

func PerformRequestWithBody(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(method, path, strings.NewReader(body))

	router.ServeHTTP(w, req)
	return w
}

func PerformAuthenticatedRequest(router *gin.Engine, method, path string) *httptest.ResponseRecorder {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(method, path, nil)
	token, _ := services.TokenService().GenerateToken(0)
	req.Header.Add("Authorization", "Bearer "+token)

	router.ServeHTTP(w, req)
	return w
}
func PerformAuthenticatedRequestWithBody(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	token, _ := services.TokenService().GenerateToken(0)
	req.Header.Add("Authorization", "Bearer "+token)

	router.ServeHTTP(w, req)
	return w
}

func TestPing(t *testing.T) {

	router := gin.Default()
	RegisterRoutes(router)

	w := PerformRequest(router, "GET", "/api/ping")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message": "pong"}`, w.Body.String())
}
