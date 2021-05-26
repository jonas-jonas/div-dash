-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "exchange" (
    exchange TEXT PRIMARY KEY,
    exchange_suffix TEXT NOT NULL,
    region TEXT NOT NULL,
    description TEXT NOT NULL,
    mic TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "exchange";
-- +goose StatementEnd
