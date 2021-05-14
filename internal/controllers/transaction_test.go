package controllers

import (
	"div-dash/internal/db"
	"div-dash/util/testutil"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"transaction_id", "symbol", "type", "transaction_provider", "price", "date", "amount", "account_id", "user_id", "side"}).
		AddRow("T1", "BTC", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, "A1", testutil.TestUserID, db.TransactionSideBuy)

	mock.ExpectQuery("^-- name: GetTransaction :one .*$").WithArgs("T1", "A1", testutil.TestUserID).WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/account/A1/transaction/T1")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{
		"transactionId":"T1",
		"symbol":"BTC",
		"type":"crypto",
		"transactionProvider":"binance",
		"price":34972.23,
		"date":
		"2021-04-09T10:00:00+02:00",
		"amount":"0.00034",
		"accountId":"A1",
		"side": "buy"
	}`, w.Body.String())
}

func TestGetTransactionDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: GetTransaction :one .*$").WithArgs("A1", "T1", testutil.TestUserID).WillReturnError(errors.New("test-error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/account/A1/transaction/T1")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestPostTransaction(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"transaction_id"}).AddRow("T1")
	mock.ExpectQuery("^-- name: CreateTransaction :one .*$").
		WithArgs(testutil.AnyTransactionId{}, "BTC", "crypto", "binance", 3497223, testutil.AnyTime{}, "0.00032", "A1", testutil.TestUserID, "buy").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"transaction_id", "symbol", "type", "transaction_provider", "price", "date", "amount", "account_id", "user_id", "side"}).
		AddRow("T1", "BTC", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, "A1", testutil.TestUserID, db.TransactionSideBuy)

	mock.ExpectQuery("^-- name: GetTransaction :one .*$").WithArgs("T1", "A1", testutil.TestUserID).WillReturnRows(rows)

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/account/A1/transaction", `{
		"symbol": "BTC",
		"type": "crypto",
		"transactionProvider": "binance",
		"price": 34972.23,
		"date": "2021-04-09T18:24:12+00:00",
		"amount": 0.00032,
		"side": "buy"
	}`)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{
		"amount":"0.00034",
		"price":34972.23,
		"date":"2021-04-09T10:00:00+02:00",
		"accountId":"A1",
		"transactionId":"T1",
		"transactionProvider":"binance",
		"symbol":"BTC",
		"type":"crypto",
		"side": "buy"
	}`, w.Body.String())
}

func TestPostTransactionMissingField(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/account/1/transaction", `{
		"symbol": "BTC",
		"type": "crypto",
		"transactionProvider": "binance",
		"price": 34972.23,
		"date": "2021-04-09T18:24:12+00:00"
	}`)

	assert.Equal(t, 400, w.Code)
	AssertErrorObject(t, "Key: 'createTransactionRequest.Amount' Error:Field validation for 'Amount' failed on the 'required' tag\nKey: 'createTransactionRequest.Side' Error:Field validation for 'Side' failed on the 'required' tag", 400, w.Body)
}

func TestPostTransactionDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: CreateTransaction :one .*$").WillReturnError(errors.New("test-error"))

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/account/A1/transaction", `{
		"symbol": "BTC",
		"type": "crypto",
		"transactionProvider": "binance",
		"price": 34972.23,
		"date": "2021-04-09T18:24:12+00:00",
		"amount": 0.00032,
		"side": "buy"
	}`)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestGetTransactions(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"transaction_id", "symbol", "type", "transaction_provider", "price", "date", "amount", "account_id", "user_id", "side"}).
		AddRow("T1", "BTC", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, "A1", testutil.TestUserID, db.TransactionSideBuy).
		AddRow("T2", "ETH", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, "A1", testutil.TestUserID, db.TransactionSideBuy).
		AddRow("T3", "DOT", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, "A1", testutil.TestUserID, db.TransactionSideBuy)
	mock.ExpectQuery("^-- name: ListTransactions :many .*$").WithArgs("A1", testutil.TestUserID).WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/account/A1/transaction")

	assert.Equal(t, 200, w.Code)

	var r []transactionResponse

	err := json.Unmarshal(w.Body.Bytes(), &r)

	assert.Nil(t, err)

	assert.Equal(t, 3, len(r))
}

func TestGetTransactionsDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: ListTransactions :many .*$").WithArgs("A1", testutil.TestUserID).WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/account/A1/transaction")

	assert.Equal(t, 500, w.Code)

	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}
