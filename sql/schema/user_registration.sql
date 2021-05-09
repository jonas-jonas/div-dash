-- name: GetUserRegistration :one
SELECT * FROM user_registration
WHERE id = $1 LIMIT 1;

-- name: GetUserRegistrationByUserId :one
SELECT * FROM user_registration
WHERE user_id = $1 LIMIT 1;

-- name: CreateUserRegistration :one
INSERT INTO user_registration (
  id, user_id, timestamp
) VALUES (
  $1, $2, $3
)
RETURNING *;
