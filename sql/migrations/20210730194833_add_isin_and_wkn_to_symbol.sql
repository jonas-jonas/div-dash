-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "symbol"
ADD COLUMN "isin" VARCHAR(12);

ALTER TABLE "symbol"
ADD COLUMN "wkn" VARCHAR(6);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "symbol"
DROP COLUMN "isin";

ALTER TABLE "symbol"
DROP COLUMN "wkn";
-- +goose StatementEnd
