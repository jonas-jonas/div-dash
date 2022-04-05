-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "asset_exchange" (
    symbol_id   TEXT,
    type        TEXT,
    source      TEXT,
    exchange    TEXT REFERENCES exchange,
    symbol      TEXT,
    primary key(symbol_id, type, source, exchange),
    foreign key(symbol_id, type, source) references symbol(symbol_id, type, source)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "asset_exchange"
-- +goose StatementEnd
