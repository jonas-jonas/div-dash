-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_registration (
    id          UUID PRIMARY KEY NOT NULL,
    user_id     TEXT NOT NULL REFERENCES "user",
    timestamp   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_registration;
-- +goose StatementEnd
