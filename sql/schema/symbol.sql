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
