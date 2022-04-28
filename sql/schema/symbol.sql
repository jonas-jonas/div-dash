-- name: GetSymbol :one
SELECT *
FROM "symbol"
WHERE symbol_id = $1;

-- name: GetSymbolByWKN :one
SELECT *
FROM "symbol"
WHERE wkn = $1;

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

-- name: BulkImportSymbol :exec
INSERT INTO "symbol" (symbol_id, type, source, precision, symbol_name, figi, cik, lei, iex_symbol)
SELECT unnest(@symbol_ids::text[]) AS symbol_id,
  unnest(@types::text[]) AS type,
  unnest(@sources::text[]) AS source,
  unnest(@precisions::int[]) AS precision,
  unnest(@symbol_names::text[]) AS symbol_name,
  unnest(@figis::text[]) AS figi,
  unnest(@ciks::text[]) AS cik,
  unnest(@leis::text[]) AS lei,
  unnest(@iex_symbols::text[]) AS iex_symbol
ON CONFLICT DO NOTHING;
  
-- name: BulkImportSymbolExchange :exec
INSERT INTO "symbol_exchange"
SELECT unnest(@symbol_ids::text[]) AS symbol_id,
  unnest(@types::text[]) AS type,
  unnest(@sources::text[]) AS source,
  unnest(@exchanges::text[]) AS exchange,
  unnest(@symbols::text[]) AS symbol
ON CONFLICT DO NOTHING;

-- name: ConnectSymbolWithExchange :exec
INSERT INTO "symbol_exchange" (symbol_id, type, source, exchange, symbol)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING;

-- name: SymbolExists :one
SELECT EXISTS(
  SELECT 1 FROM "symbol"
  WHERE symbol_id = $1 AND type = $2 AND source = $3
);

-- name: UpdateSymbol :exec
UPDATE "symbol"
SET precision = $4, symbol_name = $5
WHERE symbol_id = $1 AND type = $2 AND source = $3;

-- name: SearchSymbol :many
SELECT *
FROM "symbol"
WHERE symbol_id LIKE @search OR symbol_name LIKE @search
LIMIT @count;

-- name: AddISINAndWKN :exec
UPDATE "symbol"
SET isin = $1, wkn = $2
WHERE symbol_id = $3;
