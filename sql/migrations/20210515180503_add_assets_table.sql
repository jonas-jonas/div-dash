-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "asset" (
    asset_name  TEXT NOT NULL,
    type        TEXT NOT NULL,
    source      TEXT NOT NULL,
    precision   INT DEFAULT 2,
    PRIMARY KEY (asset_name, type, source)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "asset";
-- +goose StatementEnd
