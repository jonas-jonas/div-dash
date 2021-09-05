-- name: GetTransaction :one
SELECT * FROM "transaction"
WHERE id = $1 AND account_id = $2 AND user_id = $3
LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM "transaction"
WHERE account_id = $1 AND user_id = $2
ORDER BY date DESC;

-- name: CreateTransaction :one
INSERT INTO "transaction" (
  id, symbol, type, "transaction_provider", price, "date", amount, account_id, user_id, side, external_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id;

-- name: DeleteTransaction :exec
DELETE FROM "transaction"
WHERE id = $1 AND account_id = $2 AND user_id = $3;

-- name: TransactionExists :one
SELECT EXISTS (
  SELECT id
  FROM "transaction"
  WHERE external_id = $1
);

