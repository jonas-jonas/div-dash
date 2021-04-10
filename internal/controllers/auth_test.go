package controllers

import (
	"div-dash/internal/db"
	"div-dash/util/testutil"
	"encoding/json"
	"errors"
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
		AssertErrorObject(t, "wrong credentials", 401, w.Body)
	})

	t.Run("POST /login with missing field", func(t *testing.T) {
		loginRequest := LoginRequest{
			Email: "test@test.de",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/api/login", string(body))

		assert.Equal(t, 400, w.Code)
		AssertErrorObject(t, "Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag", 400, w.Body)
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
		AssertErrorObject(t, "wrong credentials", 401, w.Body)
	})
}

func TestPostLoginDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: FindByEmail :one .*$").WithArgs("user@example.de").WillReturnError(errors.New("test error"))

	w := PerformRequestWithBody(router, "POST", "/api/login", `{"email": "user@example.de", "password": "pass"}`)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
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
	AssertErrorObject(t, "User not activated", 401, w.Body)
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
	AssertErrorObject(t, "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag", 400, w.Body)
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
	AssertErrorObject(t, "A user with email 'user@email.de' already exists", 409, w.Body)
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
	AssertErrorObject(t, "Activation id is in wrong format", 400, w.Body)
}

func TestPostActivateNonExistingIdReturnsError(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(sqlmock.NewRows([]string{}))

	w := PerformRequest(router, "GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 400, w.Code)
	AssertErrorObject(t, "Invalid id", 400, w.Body)
}

func TestPostActivateExpiredRegistration(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now().Add(-25*time.Hour))
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	w := PerformRequest(router, "GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 400, w.Code)
	AssertErrorObject(t, "Registration expired", 400, w.Body)
}

func TestPostRegisterDbErrorFindByEmail(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: FindByEmail :one .*$").WithArgs("user@example.de").WillReturnError(errors.New("test error"))

	w := PerformRequestWithBody(router, "POST", "/api/register", `{"email": "user@example.de", "password": "pass"}`)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestPostRegisterDbErrorCreateUser(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@example.de").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("false"))

	mock.ExpectQuery("^-- name: CreateUser :one .*$").WithArgs("user@example.de", testutil.AnyString{}, db.UserStatusRegistered).WillReturnError(errors.New("test error"))

	w := PerformRequestWithBody(router, "POST", "/api/register", `{"email": "user@example.de", "password": "pass"}`)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestPostRegisterDbErrorCreateUserRegistration(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@example.de").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("false"))

	rows := sqlmock.NewRows([]string{"id", "email", "password", "status"}).
		AddRow(1, "user@example.de", "password", db.UserStatusRegistered)
	mock.ExpectQuery("^-- name: CreateUser :one .*$").WillReturnRows(rows)

	mock.ExpectQuery("^-- name: CreateUserRegistration :one .*$").WillReturnError(errors.New("test error"))

	w := PerformRequestWithBody(router, "POST", "/api/register", `{"email": "user@example.de", "password": "pass"}`)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestPostActivateDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now())
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	mock.ExpectQuery("^-- name: GetUserRegistration :exec .*$").WithArgs(1).WillReturnError(errors.New("test error"))

	w := PerformRequest(router, "GET", "/api/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}
