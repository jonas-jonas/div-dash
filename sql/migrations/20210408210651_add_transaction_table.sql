-- +goose Up
-- +goose StatementBegin
CREATE TABLE "transaction" (
    id                    text PRIMARY KEY,
    symbol                text NOT NULL,
    type                  text NOT NULL,
    transaction_provider  text NOT NULL,
    buy_in                BIGINT NOT NULL,
    buy_in_date           TIMESTAMP NOT NULL,
    amount                NUMERIC(20,10) NOT NULL,
    portfolio_id          text NOT NULL,
    side                  text NOT NULL,
    CONSTRAINT uq_id_portfolio_id
      UNIQUE(id, portfolio_id),
    CONSTRAINT fk_transaction_portfolio
      FOREIGN KEY(portfolio_id)
        REFERENCES portfolio(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "transaction";
-- +goose StatementEnd
