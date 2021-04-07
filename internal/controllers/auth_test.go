package controllers

import (
	"bytes"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/util/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	testutil.SetupConfig()
	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)
	defer sdb.Close()

	router := gin.Default()
	RegisterRoutes(router)

	t.Run("POST /login", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
			AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusActivated)

		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

		w := httptest.NewRecorder()
		loginRequest := LoginRequest{
			Email:    "email@email.de",
			Password: "pass",
		}
		body, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/api/login", bytes.NewReader(body))

		router.ServeHTTP(w, req)

		// assert.Equal(t, 200, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.NotEmpty(t, response["token"])
		assert.Regexp(t, `v2\.local\..*`, response["token"])
	})

	t.Run("POST /login with wrong credentials", func(t *testing.T) {
		w := httptest.NewRecorder()
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
			AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusActivated)

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

func setUpDb() (sqlmock.Sqlmock, func()) {
	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)

	return mock, func() {
		sdb.Close()
	}
}

func TestLoginWithNonActivatedUser(t *testing.T) {

	mock, cleanup := setUpDb()
	defer cleanup()
	router := gin.Default()
	RegisterRoutes(router)
	w := httptest.NewRecorder()
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
		AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusDeactivated)

	mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

	loginRequest := LoginRequest{
		Email:    "email@email.de",
		Password: "pass",
	}
	body, _ := json.Marshal(loginRequest)
	req, err := http.NewRequest("POST", "/api/login", bytes.NewReader(body))

	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.JSONEq(t, `{"message": "User not activated"}`, w.Body.String())
}

func TestPostRegister(t *testing.T) {
	testutil.SetupConfig()

	mock, cleanup := setUpDb()
	defer cleanup()
	router := gin.Default()
	RegisterRoutes(router)
	w := httptest.NewRecorder()
	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@email.de").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("false"))

	rows := sqlmock.NewRows([]string{"id", "email", "password", "status"}).
		AddRow(1, "email@email.de", "password", db.UserStatusRegistered)
	mock.ExpectQuery("^-- name: CreateUser :one .*$").WillReturnRows(rows)

	registerUuid, _ := uuid.NewRandom()
	rows = sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).
		AddRow(registerUuid, 1, time.Now())
	mock.ExpectQuery("^-- name: CreateUserRegistration :one .*$").WillReturnRows(rows)

	registerRequest := RegisterRequest{
		Email:    "user@email.de",
		Password: "password",
	}

	body, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/api/register", bytes.NewReader(body))
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPostRegisterWithMissingField(t *testing.T) {

	router := gin.Default()
	RegisterRoutes(router)
	w := httptest.NewRecorder()
	registerRequest := RegisterRequest{
		Email: "user@email.de",
	}

	body, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/api/register", bytes.NewReader(body))
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, w.Body.String())
}

func TestPostRegisterWithExistingUserEmail(t *testing.T) {
	mock, cleanup := setUpDb()
	defer cleanup()
	router := gin.Default()
	RegisterRoutes(router)
	w := httptest.NewRecorder()
	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@email.de").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("true"))

	registerRequest := RegisterRequest{
		Email:    "user@email.de",
		Password: "password",
	}

	body, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/api/register", bytes.NewReader(body))
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
	assert.JSONEq(t, `{"message": "A user with email 'user@email.de' already exists"}`, w.Body.String())
}

func TestPostActivate(t *testing.T) {
	mock, cleanup := setUpDb()
	defer cleanup()
	router := gin.Default()
	RegisterRoutes(router)
	w := httptest.NewRecorder()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now())
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	mock.ExpectExec("^-- name: ActivateUser :exec .*$").WithArgs(1).WillReturnResult(sqlmock.NewResult(-1, 1))

	req, err := http.NewRequest("GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Empty(t, w.Body.String())
	assert.Equal(t, 200, w.Code)
}

func TestPostActivateInvalidIdReturnsError(t *testing.T) {
	router := gin.Default()
	RegisterRoutes(router)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/api/activate?id=asdsd", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"message": "Invalid id format"}`, w.Body.String())
}

func TestPostActivateNonExistingIdReturnsError(t *testing.T) {
	mock, cleanup := setUpDb()
	defer cleanup()
	router := gin.Default()
	RegisterRoutes(router)
	w := httptest.NewRecorder()
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(sqlmock.NewRows([]string{}))

	req, err := http.NewRequest("GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"message": "Invalid id"}`, w.Body.String())
}

func TestPostActivateExpiredRegistration(t *testing.T) {
	mock, cleanup := setUpDb()
	defer cleanup()
	router := gin.Default()
	RegisterRoutes(router)
	w := httptest.NewRecorder()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now().Add(-25*time.Hour))
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"message": "Registration expired"}`, w.Body.String())
}
