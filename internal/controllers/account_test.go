package controllers

import (
	"div-dash/util/testutil"
	"errors"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestGetAccount(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"account_id", "name", "user_id"}).
		AddRow("1", "Test Account", testutil.TestUserID)

	mock.ExpectQuery("^-- name: GetAccount :one .*$").WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/account/1")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"id": "1", "name": "Test Account", "userID": "`+testutil.TestUserID+`"}`, w.Body.String())
}

func TestGetAccountNoAccount(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{})
	mock.ExpectQuery("^-- name: GetAccount :one .*$").WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/account/1")

	assert.Equal(t, 404, w.Code)
	AssertErrorObject(t, "The requested resource could not be found", 404, w.Body)
}

func TestGetAccountDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: GetAccount :one .*$").WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/account/1")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestPostAccount(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"account_id", "name", "user_id"}).
		AddRow(1, "Test Account", testutil.TestUserID)

	mock.ExpectQuery("^-- name: CreateAccount :one .*$").WithArgs(testutil.AnyAccountId{}, "Test Account", testutil.TestUserID).WillReturnRows(rows)

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/account", `{"name": "Test Account"}`)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"id": "1", "name": "Test Account", "userID": "`+testutil.TestUserID+`"}`, w.Body.String())
}

func TestPostAccountDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: CreateAccount :one .*$").WithArgs("Test Account", 1).WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/account", `{"name": "Test Account"}`)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestPostAccountMissingField(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/account", `{}`)

	assert.Equal(t, 400, w.Code)
	AssertErrorObject(t, "Key: 'createAccountRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag", 400, w.Body)
}

func TestPutAccount(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"account_id", "name", "user_id"}).
		AddRow(1, "New Test Account", testutil.TestUserID)

	mock.ExpectQuery("^-- name: UpdateAccount :one .*$").WithArgs("1", "New Test Account").WillReturnRows(rows)

	w := PerformAuthenticatedRequestWithBody(router, "PUT", "/api/account/1", `{"name": "New Test Account"}`)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"id": "1", "name": "New Test Account", "userID": "`+testutil.TestUserID+`"}`, w.Body.String())
}

func TestPutAccountDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: UpdateAccount :one .*$").WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequestWithBody(router, "PUT", "/api/account/1", `{"name": "New Test Account"}`)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestPutAccountMissingField(t *testing.T) {

	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequestWithBody(router, "PUT", "/api/account/1", `{}`)

	assert.Equal(t, 400, w.Code)
	AssertErrorObject(t, "Key: 'updateAccountRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag", 400, w.Body)
}

func TestDeleteAccount(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectExec("^-- name: DeleteAccount :exec .*$").WithArgs("1").WillReturnResult(sqlmock.NewResult(0, 1))

	w := PerformAuthenticatedRequest(router, "DELETE", "/api/account/1")

	assert.Equal(t, 200, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestDeleteAccountDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectExec("^-- name: DeleteAccount :exec .*$").WithArgs(1).WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "DELETE", "/api/account/1")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestGetAccounts(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"account_id", "name", "user_id"}).
		AddRow("1", "Test Account 1", testutil.TestUserID).
		AddRow("2", "Test Account 2", testutil.TestUserID).
		AddRow("3", "Test Account 3", testutil.TestUserID)

	mock.ExpectQuery("^-- name: ListAccounts :many .*$").WithArgs(testutil.TestUserID).WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/account")

	accountCount := gjson.Get(w.Body.String(), "#")

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int64(3), accountCount.Int())
	for i := 0; i < 3; i++ {
		account := gjson.Get(w.Body.String(), strconv.Itoa(i)+".name")
		assert.Equal(t, "Test Account "+strconv.Itoa(i+1), account.String())

	}
}

func TestGetAccountsDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: ListAccounts :many .*$").WithArgs(testutil.TestUserID).WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/account")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestGetAccountsNoResults(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: ListAccounts :many .*$").WithArgs(testutil.TestUserID).WillReturnRows(sqlmock.NewRows([]string{}))

	w := PerformAuthenticatedRequest(router, "GET", "/api/account")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `[]`, w.Body.String())
}
