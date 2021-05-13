-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM account
WHERE user_id = $1
ORDER BY id;

-- name: CreateAccount :one
INSERT INTO account (
  id, name, user_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;

-- name: UpdateAccount :one
UPDATE account
SET name = $2
WHERE id = $1
RETURNING *;
