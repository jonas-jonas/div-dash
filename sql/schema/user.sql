-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM "user"
ORDER BY id;

-- name: CreateUser :one
INSERT INTO "user" (
  id, email, password_hash, status
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;

-- name: FindByEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: CountByEmail :one
SELECT count(*) FROM "user"
WHERE email = $1;

-- name: ExistsByEmail :one
SELECT EXISTS(
  SELECT 1 FROM "user"
  WHERE email = $1
);

-- name: ActivateUser :exec
UPDATE "user"
SET status = 'activated'
WHERE id = $1;

-- name: IsUserActivated :one
SELECT EXISTS (
  SELECT status
  FROM "user"
  WHERE id = $1 AND status = 'activated'
);
