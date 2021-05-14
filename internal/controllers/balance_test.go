package controllers

import (
	"div-dash/util/testutil"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"symbol", "cost_basis", "total"}).
		AddRow("BTC", float64(10000), float64(20))

	mock.ExpectQuery("^-- name: GetBalance :many .*$").
		WithArgs(testutil.TestUserID).
		WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/balance")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `[{"costBasis":5, "symbol":"BTC", "total":20}]`, w.Body.String())
}

func TestGetBalanceDbErrorOnBalance(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: GetBalance :many .*$").
		WithArgs(testutil.TestUserID).
		WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/balance")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)

}
