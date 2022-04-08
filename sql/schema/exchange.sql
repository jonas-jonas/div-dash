-- name: CreateExchange :exec
INSERT INTO "exchange" (
    exchange, region, description, mic, exchange_suffix
) VALUES (
    $1, $2, $3, $4, $5
) ON CONFLICT DO NOTHING;

-- name: GetExchangesOfSymbol :many
SELECT e.*
FROM "symbol_exchange" ae
JOIN "exchange" e
    ON ae.exchange = e.exchange
WHERE ae.symbol_id = $1;

-- name: GetSymbolOfSymbolAndExchange :one
SELECT symbol
FROM "symbol_exchange"
WHERE symbol_id = $1 AND exchange = $2;

-- name: DoesExchangeExist :one
SELECT EXISTS(
    SELECT *
    FROM "exchange"
    WHERE exchange = $1
);
