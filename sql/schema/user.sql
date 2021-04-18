-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  email, password_hash, status
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: FindByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CountByEmail :one
SELECT count(*) FROM users
WHERE email = $1;

-- name: ExistsByEmail :one
SELECT EXISTS(
  SELECT 1 FROM users
  WHERE email = $1
);

-- name: ActivateUser :exec
UPDATE users
SET status = 'activated'
WHERE id = $1;

-- name: IsUserActivated :one
SELECT EXISTS (
  SELECT status
  FROM users
  WHERE id = $1 AND status = 'activated'
);
