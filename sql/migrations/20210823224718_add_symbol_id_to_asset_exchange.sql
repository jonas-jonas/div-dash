-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "asset_exchange" 
RENAME COLUMN symbol TO symbol_id;

ALTER TABLE "asset_exchange"
ADD COLUMN symbol TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE "asset_exchange"
REMOVE COLUMN symbol TEXT;

ALTER TABLE "asset_exchange"
RENAME COLUMN symbol_id TO symbol;
-- +goose StatementEnd
