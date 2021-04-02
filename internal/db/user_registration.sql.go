// Code generated by sqlc. DO NOT EDIT.
// source: user_registration.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUserRegistration = `-- name: CreateUserRegistration :one
INSERT INTO user_registrations (
  id, user_id, timestamp
) VALUES (
  $1, $2, $3
)
RETURNING id, user_id, timestamp
`

type CreateUserRegistrationParams struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (q *Queries) CreateUserRegistration(ctx context.Context, arg CreateUserRegistrationParams) (UserRegistration, error) {
	row := q.queryRow(ctx, q.createUserRegistrationStmt, createUserRegistration, arg.ID, arg.UserID, arg.Timestamp)
	var i UserRegistration
	err := row.Scan(&i.ID, &i.UserID, &i.Timestamp)
	return i, err
}

const getUserRegistration = `-- name: GetUserRegistration :one
SELECT id, user_id, timestamp FROM user_registrations
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserRegistration(ctx context.Context, id uuid.UUID) (UserRegistration, error) {
	row := q.queryRow(ctx, q.getUserRegistrationStmt, getUserRegistration, id)
	var i UserRegistration
	err := row.Scan(&i.ID, &i.UserID, &i.Timestamp)
	return i, err
}

const getUserRegistrationByUserId = `-- name: GetUserRegistrationByUserId :one
SELECT id, user_id, timestamp FROM user_registrations
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUserRegistrationByUserId(ctx context.Context, userID int64) (UserRegistration, error) {
	row := q.queryRow(ctx, q.getUserRegistrationByUserIdStmt, getUserRegistrationByUserId, userID)
	var i UserRegistration
	err := row.Scan(&i.ID, &i.UserID, &i.Timestamp)
	return i, err
}
