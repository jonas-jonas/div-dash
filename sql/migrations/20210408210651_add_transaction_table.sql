-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE transaction_provider AS ENUM ('binance');
CREATE TYPE transaction_type AS ENUM ('crypto');
CREATE TYPE transaction_side AS ENUM ('sell', 'buy');
CREATE TABLE transaction(
    transaction_id        BIGSERIAL PRIMARY KEY,
    symbol                text NOT NULL,
    type                  transaction_type NOT NULL,
    transaction_provider  transaction_provider NOT NULL,
    buy_in                BIGINT NOT NULL,
    buy_in_date           TIMESTAMP NOT NULL,
    amount                NUMERIC(20,10) NOT NULL,
    portfolio_id          BIGSERIAL NOT NULL,
    side                  transaction_side NOT NULL,
    CONSTRAINT fk_portfolio_portfolio_id_portfolio_id
      FOREIGN KEY(portfolio_id)
        REFERENCES portfolio(portfolio_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transaction;
DROP TYPE transaction_type;
DROP TYPE transaction_provider;
DROP TYPE transaction_side;
-- +goose StatementEnd
