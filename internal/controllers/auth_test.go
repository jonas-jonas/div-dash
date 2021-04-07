package controllers

import (
	"div-dash/internal/db"
	"div-dash/util/testutil"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	t.Run("POST /login", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
			AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusActivated)

		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

		loginRequest := LoginRequest{
			Email:    "email@email.de",
			Password: "pass",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/api/login", string(body))

		// assert.Equal(t, 200, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.NotEmpty(t, response["token"])
		assert.Regexp(t, `v2\.local\..*`, response["token"])
	})

	t.Run("POST /login with wrong credentials", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
			AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusActivated)

		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

		loginRequest := LoginRequest{
			Email:    "email@email.de",
			Password: "wrong-password",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/api/login", string(body))

		assert.Equal(t, 401, w.Code)
		assert.JSONEq(t, `{"message": "wrong credentials"}`, w.Body.String())
	})

	t.Run("POST /login with missing field", func(t *testing.T) {
		loginRequest := LoginRequest{
			Email: "test@test.de",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/api/login", string(body))

		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, `{"error": "Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, w.Body.String())
	})

	t.Run("POST /login for missing user", func(t *testing.T) {
		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(sqlmock.NewRows([]string{}))

		loginRequest := LoginRequest{
			Email:    "non-existent@test.de",
			Password: "wrong-password",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/api/login", string(body))

		assert.Equal(t, 401, w.Code)
		assert.JSONEq(t, `{"message":"wrong credentials"}`, w.Body.String())
	})
}

func TestLoginWithNonActivatedUser(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
		AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusDeactivated)

	mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

	loginRequest := LoginRequest{
		Email:    "email@email.de",
		Password: "pass",
	}
	body, _ := json.Marshal(loginRequest)
	w := PerformRequestWithBody(router, "POST", "/api/login", string(body))

	assert.Equal(t, 401, w.Code)
	assert.JSONEq(t, `{"message": "User not activated"}`, w.Body.String())
}

func TestPostRegister(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
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

	w := PerformRequestWithBody(router, "POST", "/api/register", string(body))

	assert.Equal(t, 200, w.Code)
}

func TestPostRegisterWithMissingField(t *testing.T) {

	_, cleanup, router := NewApi()

	defer cleanup()
	registerRequest := RegisterRequest{
		Email: "user@email.de",
	}

	body, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	w := PerformRequestWithBody(router, "POST", "/api/register", string(body))

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, w.Body.String())
}

func TestPostRegisterWithExistingUserEmail(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@email.de").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("true"))

	registerRequest := RegisterRequest{
		Email:    "user@email.de",
		Password: "password",
	}

	body, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	w := PerformRequestWithBody(router, "POST", "/api/register", string(body))

	assert.Equal(t, 409, w.Code)
	assert.JSONEq(t, `{"message": "A user with email 'user@email.de' already exists"}`, w.Body.String())
}

func TestPostActivate(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now())
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	mock.ExpectExec("^-- name: ActivateUser :exec .*$").WithArgs(1).WillReturnResult(sqlmock.NewResult(-1, 1))

	w := PerformRequest(router, "GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Empty(t, w.Body.String())
	assert.Equal(t, 200, w.Code)
}

func TestPostActivateInvalidIdReturnsError(t *testing.T) {
	_, cleanup, router := NewApi()
	defer cleanup()

	w := PerformRequest(router, "GET", "/api/activate?id=asdsd")

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"message": "Invalid id format"}`, w.Body.String())
}

func TestPostActivateNonExistingIdReturnsError(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(sqlmock.NewRows([]string{}))

	w := PerformRequest(router, "GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"message": "Invalid id"}`, w.Body.String())
}

func TestPostActivateExpiredRegistration(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now().Add(-25*time.Hour))
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	w := PerformRequest(router, "GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"message": "Registration expired"}`, w.Body.String())
}
