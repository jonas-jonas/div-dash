package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {

	router := gin.Default()
	RegisterRoutes(router)

	t.Run("POST /login", func(t *testing.T) {

		w := httptest.NewRecorder()
		loginRequest := LoginRequest{
			Email:    "test@test.de",
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
		loginRequest := LoginRequest{
			Email:    "test@test.de",
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

}
