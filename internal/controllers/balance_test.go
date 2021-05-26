package controllers

import (
	"div-dash/util/testutil"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestGetBalance(t *testing.T) {
	mock, cleanup, router := NewApi()

	// TODO: Mock Price Services here

	defer cleanup()
	rows := sqlmock.NewRows([]string{"symbol", "cost_basis", "amount"}).
		AddRow("BTC", float64(10000), float64(20))

	mock.ExpectQuery("^-- name: GetBalance :many .*$").
		WithArgs(testutil.TestUserID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"asset_name", "type", "source", "precision"}).
		AddRow("BTC", "crypto", "binance", 8)

	mock.ExpectQuery("^-- name: GetAsset :one .*$").
		WithArgs("BTC").
		WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/balance")

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()

	result := gjson.Parse(body)

	assert.Equal(t, result.Get("#").Int(), int64(1))
	assert.Equal(t, result.Get("0.asset.assetName").String(), "BTC")
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
