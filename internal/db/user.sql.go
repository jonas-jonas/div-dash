// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
)

const activateUser = `-- name: ActivateUser :exec
UPDATE users
SET status = 'activated'
WHERE id = $1
`

func (q *Queries) ActivateUser(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.activateUserStmt, activateUser, id)
	return err
}

const countByEmail = `-- name: CountByEmail :one
SELECT count(*) FROM users
WHERE email = $1
`

func (q *Queries) CountByEmail(ctx context.Context, email string) (int64, error) {
	row := q.queryRow(ctx, q.countByEmailStmt, countByEmail, email)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  email, password_hash, status
) VALUES (
  $1, $2, $3
)
RETURNING id, email, password_hash, status
`

type CreateUserParams struct {
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	Status       UserStatus `json:"status"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser, arg.Email, arg.PasswordHash, arg.Status)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.Status,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteUserStmt, deleteUser, id)
	return err
}

const existsByEmail = `-- name: ExistsByEmail :one
SELECT EXISTS(
  SELECT 1 FROM users
  WHERE email = $1
)
`

func (q *Queries) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	row := q.queryRow(ctx, q.existsByEmailStmt, existsByEmail, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const findByEmail = `-- name: FindByEmail :one
SELECT id, email, password_hash, status FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) FindByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.findByEmailStmt, findByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.Status,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, email, password_hash, status FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.queryRow(ctx, q.getUserStmt, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.Status,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, password_hash, status FROM users
ORDER BY id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.query(ctx, q.listUsersStmt, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.PasswordHash,
			&i.Status,
		); err != nil {
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
