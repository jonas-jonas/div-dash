-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_status AS ENUM ('registered', 'activated', 'deactivated');
ALTER TABLE users
ADD status user_status;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN status;

DROP TYPE user_status;
-- +goose StatementEnd
