-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "symbol"
ADD COLUMN iex_symbol TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "symbol"
DROP COLUMN iex_symbol;
-- +goose StatementEnd
