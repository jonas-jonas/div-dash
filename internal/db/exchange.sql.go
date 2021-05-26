// Code generated by sqlc. DO NOT EDIT.
// source: exchange.sql

package db

import (
	"context"
)

const createExchange = `-- name: CreateExchange :exec
INSERT INTO "exchange" (
    exchange, region, description, mic, exchange_suffix
) VALUES (
    $1, $2, $3, $4, $5
) ON CONFLICT DO NOTHING
`

type CreateExchangeParams struct {
	Exchange       string `json:"exchange"`
	Region         string `json:"region"`
	Description    string `json:"description"`
	Mic            string `json:"mic"`
	ExchangeSuffix string `json:"exchangeSuffix"`
}

func (q *Queries) CreateExchange(ctx context.Context, arg CreateExchangeParams) error {
	_, err := q.exec(ctx, q.createExchangeStmt, createExchange,
		arg.Exchange,
		arg.Region,
		arg.Description,
		arg.Mic,
		arg.ExchangeSuffix,
	)
	return err
}

const getExchangesOfAsset = `-- name: GetExchangesOfAsset :many
SELECT e.exchange, e.exchange_suffix, e.region, e.description, e.mic
FROM "asset_exchange" ae
JOIN "exchange" e
    ON ae.exchange = e.exchange
WHERE ae.symbol = $1
`

func (q *Queries) GetExchangesOfAsset(ctx context.Context, symbol string) ([]Exchange, error) {
	rows, err := q.query(ctx, q.getExchangesOfAssetStmt, getExchangesOfAsset, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Exchange
	for rows.Next() {
		var i Exchange
		if err := rows.Scan(
			&i.Exchange,
			&i.ExchangeSuffix,
			&i.Region,
			&i.Description,
			&i.Mic,
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