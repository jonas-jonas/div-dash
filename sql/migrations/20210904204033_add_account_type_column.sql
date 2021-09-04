-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "account_type" (
    account_type    TEXT PRIMARY KEY,
    label           TEXT NOT NULL
);

ALTER TABLE "account"
ADD COLUMN "account_type" TEXT REFERENCES account_type(account_type);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE "account"
DROP COLUMN "account_type";

DROP TABLE "account_type";
-- +goose StatementEnd
