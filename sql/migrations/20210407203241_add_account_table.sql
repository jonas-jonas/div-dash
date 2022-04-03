-- +goose Up
-- +goose StatementBegin
CREATE TABLE account (
    id              text PRIMARY KEY,
    name            text NOT NULL,
    user_id         text NOT NULL REFERENCES "user"
    
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE account;
-- +goose StatementEnd
