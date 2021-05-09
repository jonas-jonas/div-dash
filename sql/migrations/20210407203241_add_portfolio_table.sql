-- +goose Up
-- +goose StatementBegin
CREATE TABLE portfolio (
    id              text PRIMARY KEY,
    name            text NOT NULL,
    user_id         text NOT NULL,
    CONSTRAINT fk_portfolio_user
      FOREIGN KEY(user_id)
        REFERENCES "user"(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE portfolio;
-- +goose StatementEnd
