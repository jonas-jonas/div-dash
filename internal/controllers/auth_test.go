package controllers

import (
	"bytes"
	"div-dash/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)
	defer sdb.Close()

	router := gin.Default()
	RegisterRoutes(router)

	t.Run("POST /login", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(1, "email@email.de", "password")

		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

		w := httptest.NewRecorder()
		loginRequest := LoginRequest{
			Email:    "email@email.de",
			Password: "password",
		}
		body, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/api/login", bytes.NewReader(body))

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fail()
		}
		assert.NotNil(t, response["token"])
		assert.Regexp(t, `v2\.local\..*`, response["token"])
	})

	t.Run("POST /login with wrong credentials", func(t *testing.T) {
		w := httptest.NewRecorder()
		rows := sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(1, "email@email.de", "password")

		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

		loginRequest := LoginRequest{
			Email:    "email@email.de",
			Password: "wrong-password",
		}
		body, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/api/login", bytes.NewReader(body))

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.JSONEq(t, `{"message": "wrong credentials"}`, w.Body.String())
	})

	t.Run("POST /login with missing field", func(t *testing.T) {
		w := httptest.NewRecorder()
		loginRequest := LoginRequest{
			Email: "test@test.de",
		}
		body, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/api/login", bytes.NewReader(body))

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, `{"error": "Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, w.Body.String())
	})

	t.Run("POST /login for missing user", func(t *testing.T) {
		w := httptest.NewRecorder()
		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(sqlmock.NewRows([]string{}))

		loginRequest := LoginRequest{
			Email:    "non-existent@test.de",
			Password: "wrong-password",
		}
		body, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/api/login", bytes.NewReader(body))

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.JSONEq(t, `{"message":"wrong credentials"}`, w.Body.String())
	})
}
