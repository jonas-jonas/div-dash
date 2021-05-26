-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "asset_exchange" (
    symbol      TEXT,
    exchange    TEXT,
    PRIMARY KEY (symbol, exchange),
    CONSTRAINT fk_asset_exchange_asset
      FOREIGN KEY(symbol)
        REFERENCES "asset"(assetName),
    CONSTRAINT fk_asset_exchange_exchange
      FOREIGN KEY(exchange)
        REFERENCES "exchange"(exchange)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "asset_exchange"
-- +goose StatementEnd
