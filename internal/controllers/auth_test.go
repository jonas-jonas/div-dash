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

		loginRequest := LoginFormRequest{
			Email:    "email@email.de",
			Password: "pass",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/login", string(body))

		assert.Equal(t, 303, w.Code)
		cookies := w.Result().Cookies()
		assert.Equal(t, 1, len(cookies))
		cookie := cookies[0]
		assert.Regexp(t, `v2\.local\..*`, cookie.Value)
	})

	t.Run("POST /login with wrong credentials", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
			AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusActivated)

		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

		loginRequest := LoginFormRequest{
			Email:    "email@email.de",
			Password: "wrong-password",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/login", string(body))

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid password")
	})

	t.Run("POST /login with missing field", func(t *testing.T) {
		loginRequest := LoginFormRequest{
			Email: "test@test.de",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/login", string(body))

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "Password is required")
	})

	t.Run("POST /login for missing user", func(t *testing.T) {
		mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(sqlmock.NewRows([]string{}))

		loginRequest := LoginFormRequest{
			Email:    "non-existent@test.de",
			Password: "wrong-password",
		}
		body, _ := json.Marshal(loginRequest)
		w := PerformRequestWithBody(router, "POST", "/login", string(body))

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid password")
	})
}

func TestPostLoginDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: FindByEmail :one .*$").WithArgs("user@example.de").WillReturnError(errors.New("test error"))

	w := PerformRequestWithBody(router, "POST", "/login", `{"email": "user@example.de", "password": "pass"}`)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error.")
}

func TestLoginWithNonActivatedUser(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
		AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusDeactivated)

	mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

	loginRequest := LoginFormRequest{
		Email:    "email@email.de",
		Password: "pass",
	}
	body, _ := json.Marshal(loginRequest)
	w := PerformRequestWithBody(router, "POST", "/login", string(body))

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "User not activated")
}

func TestLoginWithReferer(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "status"}).
		AddRow(1, "email@email.de", testutil.PasswordHash, db.UserStatusActivated)

	mock.ExpectQuery("^-- name: FindByEmail :one .*$").WillReturnRows(rows)

	loginRequest := LoginFormRequest{
		Email:    "email@email.de",
		Password: "pass",
		Referer:  "http://example.com/referer",
	}
	body, _ := json.Marshal(loginRequest)

	w := PerformRequestWithBody(router, "POST", "/login", string(body))

	assert.Equal(t, 303, w.Code)
	assert.Equal(t, w.Header().Get("Location"), "http://example.com/referer")
}
func TestGetLoginForm(t *testing.T) {

	_, _, router := NewApi()

	w := PerformRequest(router, "GET", "/login")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Login")
	assert.Contains(t, w.Body.String(), "Email")
	assert.Contains(t, w.Body.String(), "Password")
	assert.Contains(t, w.Body.String(), "bulma")
}

func TestGetRegisterForm(t *testing.T) {

	_, _, router := NewApi()

	w := PerformRequest(router, "GET", "/register")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Register")
	assert.Contains(t, w.Body.String(), "Email")
	assert.Contains(t, w.Body.String(), "Password")
	assert.Contains(t, w.Body.String(), "Repeat Password")
	assert.Contains(t, w.Body.String(), "bulma")
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

	registerRequest := RegisterFormRequest{
		Email:          "user@email.de",
		Password:       "password",
		RepeatPassword: "password",
		AcceptTOS:      true,
	}

	body, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	w := PerformRequestWithBody(router, "POST", "/register", string(body))

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Success")
	assert.Contains(t, w.Body.String(), "Please check your emails to verify ownership of")
}

func TestPostRegisterWithMissingField(t *testing.T) {

	_, cleanup, router := NewApi()

	defer cleanup()
	registerRequest := RegisterFormRequest{
		Email: "user@email.de",
	}

	body, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	w := PerformRequestWithBody(router, "POST", "/register", string(body))

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Password is required")
	assert.Contains(t, w.Body.String(), "RepeatPassword is required")
}

func TestPostRegisterWithExistingUserEmail(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@email.de").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("true"))

	registerRequest := RegisterFormRequest{
		Email:          "user@email.de",
		Password:       "password",
		RepeatPassword: "password",
		AcceptTOS:      true,
	}

	body, err := json.Marshal(registerRequest)
	assert.Nil(t, err)

	w := PerformRequestWithBody(router, "POST", "/register", string(body))

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "A user with email &#39;user@email.de&#39; already exists")
}

func TestPostActivate(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now())
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	mock.ExpectQuery("^-- name: IsUserActivated :one .*$").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("false"))

	mock.ExpectExec("^-- name: ActivateUser :exec .*$").WithArgs(1).WillReturnResult(sqlmock.NewResult(-1, 1))

	w := PerformRequest(router, "GET", "/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Success")
	assert.Contains(t, w.Body.String(), "You can now")
}

func TestPostActivateInvalidIdReturnsError(t *testing.T) {
	_, cleanup, router := NewApi()
	defer cleanup()

	w := PerformRequest(router, "GET", "/activate?id=asdsd")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Error")
	assert.Contains(t, w.Body.String(), "Invalid activation id")
}

func TestPostActivateNonExistingIdReturnsError(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(sqlmock.NewRows([]string{}))

	w := PerformRequest(router, "GET", "/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Error")
	assert.Contains(t, w.Body.String(), "Invalid activation id")
}

func TestPostActivateExpiredRegistration(t *testing.T) {
	mock, cleanup, router := NewApi()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now().Add(-25*time.Hour))
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	w := PerformRequest(router, "GET", "/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Error")
	assert.Contains(t, w.Body.String(), "This activation Id expired. Please register again.")
}

func TestPostRegisterDbErrorFindByEmail(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@example.de").WillReturnError(errors.New("test error"))

	w := PerformRequestWithBody(router, "POST", "/register", `{"email": "user@example.de", "password": "pass", "repeatPassword": "pass", "acceptTOS": true}`)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Register")
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

func TestPostRegisterDbErrorCreateUser(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@example.de").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("false"))

	mock.ExpectQuery("^-- name: CreateUser :one .*$").WithArgs("user@example.de", testutil.AnyString{}, db.UserStatusRegistered).WillReturnError(errors.New("test error"))

	w := PerformRequestWithBody(router, "POST", "/register", `{"email": "user@example.de", "password": "pass", "repeatPassword": "pass", "acceptTOS": true}`)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Register")
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

func TestPostRegisterDbErrorCreateUserRegistration(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: ExistsByEmail :one .*$").WithArgs("user@example.de").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("false"))

	rows := sqlmock.NewRows([]string{"id", "email", "password", "status"}).
		AddRow(1, "user@example.de", "password", db.UserStatusRegistered)
	mock.ExpectQuery("^-- name: CreateUser :one .*$").WillReturnRows(rows)

	mock.ExpectQuery("^-- name: CreateUserRegistration :one .*$").WillReturnError(errors.New("test error"))

	w := PerformRequestWithBody(router, "POST", "/register", `{"email": "user@example.de", "password": "pass", "repeatPassword": "pass", "acceptTOS": true}`)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Register")
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

func TestPostActivateDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now())
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)

	mock.ExpectQuery("^-- name: GetUserRegistration :exec .*$").WithArgs(1).WillReturnError(errors.New("test error"))

	w := PerformRequest(router, "GET", "/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Error")
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

func TestPostActivateDbErrorActivateUser(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now())
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)
	mock.ExpectQuery("^-- name: IsUserActivated :one .*$").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("false"))

	mock.ExpectExec("^-- name: ActivateUser :exec .*$").WithArgs(1).WillReturnError(errors.New("test error"))

	w := PerformRequest(router, "GET", "/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Error")
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

func TestPostActivateActivatedUser(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "user_id", "timestamp"}).AddRow("5cf6a941-2517-4e2b-9905-97f507f928c4", 1, time.Now())
	mock.ExpectQuery("^-- name: GetUserRegistration :one .*$").WithArgs("5cf6a941-2517-4e2b-9905-97f507f928c4").WillReturnRows(rows)
	mock.ExpectQuery("^-- name: IsUserActivated :one .*$").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow("true"))

	w := PerformRequest(router, "GET", "/activate?id=5cf6a941-2517-4e2b-9905-97f507f928c4")

	assert.Equal(t, 303, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}
