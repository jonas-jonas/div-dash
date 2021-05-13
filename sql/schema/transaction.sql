-- name: GetTransaction :one
SELECT * FROM "transaction"
WHERE id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM "transaction"
WHERE account_id = $1
ORDER BY date DESC;

-- name: CreateTransaction :one
INSERT INTO "transaction" (
  id, symbol, type, "transaction_provider", price, "date", amount, account_id, side
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id;

-- name: DeleteTransaction :exec
DELETE FROM "transaction"
WHERE id = $1;
