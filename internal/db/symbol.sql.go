// Code generated by sqlc. DO NOT EDIT.
// source: symbol.sql

package db

import (
	"context"
	"database/sql"
)

const addISINAndWKN = `-- name: AddISINAndWKN :exec
UPDATE "symbol"
SET isin = $1, wkn = $2
WHERE symbol_id = $3
`

type AddISINAndWKNParams struct {
	Isin     sql.NullString `json:"isin"`
	Wkn      sql.NullString `json:"wkn"`
	SymbolID string         `json:"symbolID"`
}

func (q *Queries) AddISINAndWKN(ctx context.Context, arg AddISINAndWKNParams) error {
	_, err := q.exec(ctx, q.addISINAndWKNStmt, addISINAndWKN, arg.Isin, arg.Wkn, arg.SymbolID)
	return err
}

const addSymbol = `-- name: AddSymbol :exec
INSERT INTO "symbol" (symbol_id, type, source, precision, symbol_name)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING
`

type AddSymbolParams struct {
	SymbolID   string         `json:"symbolID"`
	Type       string         `json:"type"`
	Source     string         `json:"source"`
	Precision  int32          `json:"precision"`
	SymbolName sql.NullString `json:"symbolName"`
}

func (q *Queries) AddSymbol(ctx context.Context, arg AddSymbolParams) error {
	_, err := q.exec(ctx, q.addSymbolStmt, addSymbol,
		arg.SymbolID,
		arg.Type,
		arg.Source,
		arg.Precision,
		arg.SymbolName,
	)
	return err
}

const connectSymbolWithExchange = `-- name: ConnectSymbolWithExchange :exec
INSERT INTO "asset_exchange" (symbol, exchange)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`

type ConnectSymbolWithExchangeParams struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

func (q *Queries) ConnectSymbolWithExchange(ctx context.Context, arg ConnectSymbolWithExchangeParams) error {
	_, err := q.exec(ctx, q.connectSymbolWithExchangeStmt, connectSymbolWithExchange, arg.Symbol, arg.Exchange)
	return err
}

const getSymbol = `-- name: GetSymbol :one
SELECT symbol_id, type, source, precision, symbol_name, isin, wkn
FROM "symbol"
WHERE symbol_id = $1
`

func (q *Queries) GetSymbol(ctx context.Context, symbolID string) (Symbol, error) {
	row := q.queryRow(ctx, q.getSymbolStmt, getSymbol, symbolID)
	var i Symbol
	err := row.Scan(
		&i.SymbolID,
		&i.Type,
		&i.Source,
		&i.Precision,
		&i.SymbolName,
		&i.Isin,
		&i.Wkn,
	)
	return i, err
}

const searchSymbol = `-- name: SearchSymbol :many
SELECT symbol_id, type, source, precision, symbol_name, isin, wkn
FROM "symbol"
WHERE symbol_id LIKE $1 OR symbol_name LIKE $1
LIMIT $2
`

type SearchSymbolParams struct {
	Search string `json:"search"`
	Count  int32  `json:"count"`
}

func (q *Queries) SearchSymbol(ctx context.Context, arg SearchSymbolParams) ([]Symbol, error) {
	rows, err := q.query(ctx, q.searchSymbolStmt, searchSymbol, arg.Search, arg.Count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Symbol
	for rows.Next() {
		var i Symbol
		if err := rows.Scan(
			&i.SymbolID,
			&i.Type,
			&i.Source,
			&i.Precision,
			&i.SymbolName,
			&i.Isin,
			&i.Wkn,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const symbolExists = `-- name: SymbolExists :one
SELECT EXISTS(
  SELECT 1 FROM "symbol"
  WHERE symbol_id = $1
)
`

func (q *Queries) SymbolExists(ctx context.Context, symbolID string) (bool, error) {
	row := q.queryRow(ctx, q.symbolExistsStmt, symbolExists, symbolID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const updateSymbol = `-- name: UpdateSymbol :exec
UPDATE "symbol"
SET type = $2, source = $3, precision = $4, symbol_name = $5
WHERE symbol_id = $1
`

type UpdateSymbolParams struct {
	SymbolID   string         `json:"symbolID"`
	Type       string         `json:"type"`
	Source     string         `json:"source"`
	Precision  int32          `json:"precision"`
	SymbolName sql.NullString `json:"symbolName"`
}

func (q *Queries) UpdateSymbol(ctx context.Context, arg UpdateSymbolParams) error {
	_, err := q.exec(ctx, q.updateSymbolStmt, updateSymbol,
		arg.SymbolID,
		arg.Type,
		arg.Source,
		arg.Precision,
		arg.SymbolName,
	)
	return err
}
