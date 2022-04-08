-- +goose Up
-- +goose StatementBegin
ALTER TABLE "asset_exchange"
RENAME TO "symbol_exchange";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "symbol_exchange"
RENAME TO "asset_exchange";
-- +goose StatementEnd
