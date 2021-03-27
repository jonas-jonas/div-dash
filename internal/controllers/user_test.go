package controllers

import (
	"bytes"
	"database/sql/driver"
	"div-dash/internal/config"
	"div-dash/internal/services"
	"div-dash/internal/user"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
func TestUser(t *testing.T) {

	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)
	defer sdb.Close()

	router := gin.Default()
	RegisterRoutes(router)
	token, _ := services.TokenService().GenerateToken(0)

	t.Run("GET /user with valid id", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(1, "email@email.de", "password")

		mock.ExpectQuery("^-- name: GetUser :one .*$").WillReturnRows(rows)

		user, _ := services.UserService().CreateUser(user.CreateUserParams{
			Email:    "email@email.de",
			Password: "password",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user/"+strconv.Itoa(int(user.ID)), nil)
		req.Header.Add("Authorization", token)

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"id":1,"email":"email@email.de"}`, w.Body.String())
	})

	t.Run("GET /user with invalid id", func(t *testing.T) {
		mock.ExpectQuery("^-- name: GetUser :one .*$").WillReturnError(fmt.Errorf("sql: no rows in result set"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user/123", nil)
		req.Header.Add("Authorization", token)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.JSONEq(t, `{"message": "User with id '123' not found"}`, w.Body.String())
	})

	t.Run("GET /user with string id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user/string", nil)
		req.Header.Add("Authorization", token)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, `{"message": "Param 'id' must be int, got 'string'"}`, w.Body.String())
	})

	t.Run("POST /user", func(t *testing.T) {
		mock.ExpectQuery("-- name: CountByEmail").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		rows := sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(1, "test@email.com", "password")
		mock.ExpectQuery("-- name: CreateUser").WithArgs("test@email.com", AnyString{}).WillReturnRows(rows)

		w := httptest.NewRecorder()

		createUserRequest := CreateUserRequest{
			Email:    "test@email.com",
			Password: "password",
		}
		body, _ := json.Marshal(createUserRequest)

		req, _ := http.NewRequest("POST", "/api/user/", bytes.NewReader(body))
		req.Header.Add("Authorization", token)

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"email":"test@email.com", "id":1}`, w.Body.String())
	})

	t.Run("POST /user with missing field", func(t *testing.T) {
		w := httptest.NewRecorder()

		createUserRequest := CreateUserRequest{
			Email: "test@email.com",
		}
		body, _ := json.Marshal(createUserRequest)

		req, _ := http.NewRequest("POST", "/api/user/", bytes.NewReader(body))
		req.Header.Add("Authorization", token)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, `{"error":"Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, w.Body.String())

	})

	t.Run("POST /user with existing email", func(t *testing.T) {
		mock.ExpectQuery("-- name: CountByEmail").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		w := httptest.NewRecorder()
		createUserRequest := CreateUserRequest{
			Email:    "email@email.de",
			Password: "password",
		}
		body, _ := json.Marshal(createUserRequest)

		req, _ := http.NewRequest("POST", "/api/user/", bytes.NewReader(body))
		req.Header.Add("Authorization", token)

		router.ServeHTTP(w, req)

		assert.Equal(t, 409, w.Code)
		assert.JSONEq(t, `{"message": "A user with email 'email@email.de' already exists"}`, w.Body.String())

	})
}
