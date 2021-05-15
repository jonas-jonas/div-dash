-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "job_history" (
    id              INTEGER PRIMARY KEY,
    name            TEXT NOT NULL,
    started         BIGINT NOT NULL,
    finished        BIGINT,
    error_message   TEXT
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "job_history";
-- +goose StatementEnd
