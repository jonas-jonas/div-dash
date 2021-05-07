-- +goose Up
-- +goose StatementBegin
CREATE TABLE "transaction"(
    transaction_id        BIGSERIAL PRIMARY KEY,
    symbol                text NOT NULL,
    type                  text NOT NULL,
    transaction_provider  text NOT NULL,
    buy_in                BIGINT NOT NULL,
    buy_in_date           TIMESTAMP NOT NULL,
    amount                NUMERIC(20,10) NOT NULL,
    portfolio_id          BIGSERIAL NOT NULL,
    side                  text NOT NULL,
    CONSTRAINT fk_portfolio_portfolio_id_portfolio_id
      FOREIGN KEY(portfolio_id)
        REFERENCES portfolio(portfolio_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "transaction";
-- +goose StatementEnd
