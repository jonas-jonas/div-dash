-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO "account_type" (
    "account_type", "label"
)
VALUES (
    'comdirect', 'comdirect'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM "account_type"
WHERE account_type = 'comdirect';
-- +goose StatementEnd
