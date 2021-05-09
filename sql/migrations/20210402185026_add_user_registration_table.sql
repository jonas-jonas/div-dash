-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_registration (
    id          UUID PRIMARY KEY NOT NULL,
    user_id     TEXT NOT NULL,
    timestamp   TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_registration_user
      FOREIGN KEY(user_id)
        REFERENCES "user"(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_registration;
-- +goose StatementEnd
