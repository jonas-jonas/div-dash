// Code generated by sqlc. DO NOT EDIT.
// source: asset.sql

package db

import (
	"context"
)

const addAsset = `-- name: AddAsset :exec
INSERT INTO "asset" (asset_name, type, source, precision)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING
`

type AddAssetParams struct {
	AssetName string `json:"assetName"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Precision int32  `json:"precision"`
}

func (q *Queries) AddAsset(ctx context.Context, arg AddAssetParams) error {
	_, err := q.exec(ctx, q.addAssetStmt, addAsset,
		arg.AssetName,
		arg.Type,
		arg.Source,
		arg.Precision,
	)
	return err
}

const assetExists = `-- name: AssetExists :one
SELECT EXISTS(
  SELECT 1 FROM "asset"
  WHERE asset_name = $1
)
`

func (q *Queries) AssetExists(ctx context.Context, assetName string) (bool, error) {
	row := q.queryRow(ctx, q.assetExistsStmt, assetExists, assetName)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getAsset = `-- name: GetAsset :one
SELECT asset_name, type, source, precision
FROM "asset"
WHERE asset_name = $1
`

func (q *Queries) GetAsset(ctx context.Context, assetName string) (Asset, error) {
	row := q.queryRow(ctx, q.getAssetStmt, getAsset, assetName)
	var i Asset
	err := row.Scan(
		&i.AssetName,
		&i.Type,
		&i.Source,
		&i.Precision,
	)
	return i, err
}
