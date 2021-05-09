-- name: GetPortfolio :one
SELECT * FROM portfolio
WHERE id = $1 LIMIT 1;

-- name: ListPortfolios :many
SELECT * FROM portfolio
WHERE user_id = $1
ORDER BY id;

-- name: CreatePortfolio :one
INSERT INTO portfolio (
  id, name, user_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeletePortfolio :exec
DELETE FROM portfolio
WHERE id = $1;

-- name: UpdatePortfolio :one
UPDATE portfolio
SET name = $2
WHERE id = $1
RETURNING *;
