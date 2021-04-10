package controllers

import (
	"database/sql/driver"
	"div-dash/internal/db"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
func TestUser(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	t.Run("GET /user with valid id", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "status"}).
			AddRow(1, "email@email.de", "password", db.UserStatusActivated)

		mock.ExpectQuery("^-- name: GetUser :one .*$").WillReturnRows(rows)

		w := PerformAuthenticatedRequest(router, "GET", "/api/user/"+strconv.Itoa(1))

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"id":1,"email":"email@email.de"}`, w.Body.String())
	})

	t.Run("GET /user with invalid id", func(t *testing.T) {
		mock.ExpectQuery("^-- name: GetUser :one .*$").WillReturnError(fmt.Errorf("sql: no rows in result set"))

		w := PerformAuthenticatedRequest(router, "GET", "/api/user/123")

		assert.Equal(t, 404, w.Code)
		AssertErrorObject(t, "The requested resource could not be found", 404, w.Body)
	})

	t.Run("GET /user with string id", func(t *testing.T) {
		w := PerformAuthenticatedRequest(router, "GET", "/api/user/string")

		assert.Equal(t, 400, w.Code)
		AssertErrorObject(t, "User id is invalid", 400, w.Body)
	})

}

func TestGetUserDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: GetUser :one .*$").WithArgs(1).WillReturnError(errors.New("test-error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/user/1")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}
