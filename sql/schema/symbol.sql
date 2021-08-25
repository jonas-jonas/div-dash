-- name: GetSymbol :one
SELECT *
FROM "symbol"
WHERE symbol_id = $1;

-- name: GetSymbolsByType :many
SELECT *
FROM "symbol" s
WHERE s.type = $1
LIMIT $2;

-- name: GetSymbols :many
SELECT *
FROM "symbol"
LIMIT $1;

-- name: GetSymbolCount :one
SELECT COUNT(*)
FROM "symbol";

-- name: GetSymbolCountByType :one
SELECT COUNT(*)
FROM "symbol" s
WHERE s.type = @symbolType;

-- name: AddSymbol :exec
INSERT INTO "symbol" (symbol_id, type, source, precision, symbol_name)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING;

-- name: ConnectSymbolWithExchange :exec
INSERT INTO "asset_exchange" (symbol_id, exchange, symbol)
VALUES ($1, $2, $3)
ON CONFLICT DO UPDATE SET symbol = $3;

-- name: SymbolExists :one
SELECT EXISTS(
  SELECT 1 FROM "symbol"
  WHERE symbol_id = $1
);

-- name: UpdateSymbol :exec
UPDATE "symbol"
SET type = $2, source = $3, precision = $4, symbol_name = $5
WHERE symbol_id = $1;

-- name: SearchSymbol :many
SELECT *
FROM "symbol"
WHERE symbol_id LIKE @search OR symbol_name LIKE @search
LIMIT @count;

-- name: AddISINAndWKN :exec
UPDATE "symbol"
SET isin = $1, wkn = $2
WHERE symbol_id = $3;
