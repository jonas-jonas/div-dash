-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id            text      PRIMARY KEY,
  email         text      NOT NULL,
  password_hash text      NOT NULL,
  status        text      NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd