-- name: GetTransaction :one
SELECT * FROM transaction
WHERE transaction_id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transaction
WHERE portfolio_id = $1
ORDER BY buy_in_date DESC;

-- name: CreateTransaction :one
INSERT INTO transaction (
  symbol, type, transaction_provider, buy_in, buy_in_date, amount, portfolio_id, side
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transaction
WHERE transaction_id = $1;
