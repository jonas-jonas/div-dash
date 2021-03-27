-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id       BIGSERIAL PRIMARY KEY,
  email    text      NOT NULL,
  password text      NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd