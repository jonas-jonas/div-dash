-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO "account_type" (
    "account_type", "label"
)
VALUES (
    'scalable_capital', 'scalable.capital'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM "account"
WHERE account_type = 'scalable_capital';

DELETE FROM "account_type"
WHERE "account_type" = 'scalable_capital';
-- +goose StatementEnd
