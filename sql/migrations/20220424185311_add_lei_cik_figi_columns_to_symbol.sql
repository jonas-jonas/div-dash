-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "symbol"
ADD COLUMN figi VARCHAR(12),
ADD COLUMN cik VARCHAR(10),
ADD COLUMN lei VARCHAR(20);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "symbol"
DROP COLUMN lei,
DROP COLUMN cik,
DROP COLUMN figi;
-- +goose StatementEnd
