// Code generated by sqlc. DO NOT EDIT.
// source: account.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO account (
  id, name, user_id
) VALUES (
  $1, $2, $3
)
RETURNING id, name, user_id
`

type CreateAccountParams struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"userID"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.queryRow(ctx, q.createAccountStmt, createAccount, arg.ID, arg.Name, arg.UserID)
	var i Account
	err := row.Scan(&i.ID, &i.Name, &i.UserID)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id string) error {
	_, err := q.exec(ctx, q.deleteAccountStmt, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, name, user_id FROM account
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id string) (Account, error) {
	row := q.queryRow(ctx, q.getAccountStmt, getAccount, id)
	var i Account
	err := row.Scan(&i.ID, &i.Name, &i.UserID)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, name, user_id FROM account
WHERE user_id = $1
ORDER BY id
`

func (q *Queries) ListAccounts(ctx context.Context, userID string) ([]Account, error) {
	rows, err := q.query(ctx, q.listAccountsStmt, listAccounts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(&i.ID, &i.Name, &i.UserID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE account
SET name = $2
WHERE id = $1
RETURNING id, name, user_id
`

type UpdateAccountParams struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.queryRow(ctx, q.updateAccountStmt, updateAccount, arg.ID, arg.Name)
	var i Account
	err := row.Scan(&i.ID, &i.Name, &i.UserID)
	return i, err
}
