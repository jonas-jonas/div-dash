-- +goose Up
-- +goose StatementBegin
CREATE TABLE portfolio(
    portfolio_id    BIGSERIAL PRIMARY KEY,
    name            text NOT NULL,
    user_id         text NOT NULL,
    CONSTRAINT fk_portfolio_user_id_user_id
      FOREIGN KEY(user_id)
        REFERENCES users(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE portfolio;
-- +goose StatementEnd
