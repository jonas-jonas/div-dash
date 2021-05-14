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
	rows := sqlmock.NewRows([]string{"symbol", "total"}).
		AddRow("BTC", float64(23.2312))

	mock.ExpectQuery("^-- name: GetBalance :many .*$").
		WithArgs(testutil.TestUserID).
		WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"cost_basis"}).
		AddRow(int64(4623323))
	mock.ExpectQuery("^-- name: GetCostBasis :one .*$").
		WithArgs("BTC", testutil.TestUserID).
		WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/balance")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `[{"costBasis":46233.23, "symbol":"BTC", "total":23.2312}]`, w.Body.String())
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

func TestGetBalanceDbErrorOnCostBasis(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"symbol", "total"}).
		AddRow("BTC", float64(23.2312))

	mock.ExpectQuery("^-- name: GetBalance :many .*$").
		WithArgs(testutil.TestUserID).
		WillReturnRows(rows)
	mock.ExpectQuery("^-- name: GetCostBasis :one .*$").
		WithArgs("BTC", testutil.TestUserID).
		WillReturnError(errors.New("test-error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/balance")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)

}
