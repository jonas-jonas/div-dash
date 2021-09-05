-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "transaction"
ADD COLUMN external_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE "transaction"
DROP COLUMN external_id;
-- +goose StatementEnd
