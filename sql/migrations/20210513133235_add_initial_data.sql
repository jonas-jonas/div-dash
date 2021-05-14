-- +goose Up
-- +goose StatementBegin
INSERT INTO "user" (id, email, password_hash, status)
VALUES ('YH8UFLWMGXQO4KPD', 'admin@example.com', '$2a$10$EZVctwNMgfZjNkKNjGgJqOhH0hgpMHlbMYznz.rjEYy6ZAxMMimQa', 'activated') ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
