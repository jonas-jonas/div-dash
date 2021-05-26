-- name: CreateExchange :exec
INSERT INTO "exchange" (
    exchange, region, description, mic, exchange_suffix
) VALUES (
    $1, $2, $3, $4, $5
) ON CONFLICT DO NOTHING;

-- name: GetExchangesOfAsset :many
SELECT e.*
FROM "asset_exchange" ae
JOIN "exchange" e
    ON ae.exchange = e.exchange
WHERE ae.symbol = $1;
