package controllers

import (
	"div-dash/internal/db"
	"div-dash/util/testutil"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	t.Run("GET /user with valid id", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "status"}).
			AddRow(testutil.TestUserID, "email@email.de", "password", db.UserStatusActivated)

		mock.ExpectQuery("^-- name: GetUser :one .*$").WillReturnRows(rows)

		w := PerformAuthenticatedRequest(router, "GET", "/api/user/"+testutil.TestUserID)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"id":"`+testutil.TestUserID+`","email":"email@email.de"}`, w.Body.String())
	})

	t.Run("GET /user with invalid id", func(t *testing.T) {
		mock.ExpectQuery("^-- name: GetUser :one .*$").WillReturnError(fmt.Errorf("sql: no rows in result set"))

		w := PerformAuthenticatedRequest(router, "GET", "/api/user/INVALID-ID")

		assert.Equal(t, 404, w.Code)
		AssertErrorObject(t, "The requested resource could not be found", 404, w.Body)
	})

}

func TestGetUserDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: GetUser :one .*$").WithArgs(testutil.TestUserID).WillReturnError(errors.New("test-error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/user/"+testutil.TestUserID)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}
