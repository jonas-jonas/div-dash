-- name: GetAsset :one
SELECT *
FROM "asset"
WHERE asset_name = $1;

-- name: AddAsset :exec
INSERT INTO "asset" (asset_name, type, source, precision)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;

-- name: ConnectAssetWithExchange :exec
INSERT INTO "asset_exchange" (symbol, exchange)
VALUES ($1, $2);

-- name: AssetExists :one
SELECT EXISTS(
  SELECT 1 FROM "asset"
  WHERE asset_name = $1
);
