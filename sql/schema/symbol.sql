-- name: GetSymbol :one
SELECT *
FROM "symbol"
WHERE symbol_id = $1;

-- name: AddSymbol :exec
INSERT INTO "symbol" (symbol_id, type, source, precision, symbol_name)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING;

-- name: ConnectSymbolWithExchange :exec
INSERT INTO "asset_exchange" (symbol, exchange)
VALUES ($1, $2);

-- name: SymbolExists :one
SELECT EXISTS(
  SELECT 1 FROM "symbol"
  WHERE symbol_id = $1
);

-- name: SearchSymbol :many
SELECT *
FROM "symbol"
WHERE symbol_id LIKE @search OR symbol_name LIKE @search
LIMIT @count;

-- name: AddISINAndWKN :exec
UPDATE "symbol"
SET isin = $1, wkn = $2
WHERE symbol_id = $3;
