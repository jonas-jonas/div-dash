-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS user_registrations(
    id          UUID PRIMARY KEY NOT NULL,
    user_id     BIGSERIAL NOT NULL,
    timestamp   TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_userId
      FOREIGN KEY(user_id)
        REFERENCES users(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE user_registrations;
-- +goose StatementEnd
