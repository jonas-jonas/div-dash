-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "symbol" (
    symbol_id    TEXT NOT NULL,
    type        TEXT NOT NULL,
    source      TEXT NOT NULL,
    precision   INT NOT NULL DEFAULT 2,
    symbol_name  TEXT,
    PRIMARY KEY (symbol_id, type, source)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "symbol";
-- +goose StatementEnd
