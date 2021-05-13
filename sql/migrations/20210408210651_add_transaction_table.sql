-- +goose Up
-- +goose StatementBegin
CREATE TABLE "transaction" (
    id                    text PRIMARY KEY,
    symbol                text NOT NULL,
    type                  text NOT NULL,
    transaction_provider  text NOT NULL,
    price                 BIGINT NOT NULL,
    date                  TIMESTAMP NOT NULL,
    amount                NUMERIC(20,10) NOT NULL,
    account_id            text NOT NULL,
    side                  text NOT NULL,
    CONSTRAINT uq_id_account_id
      UNIQUE(id, account_id),
    CONSTRAINT fk_transaction_account
      FOREIGN KEY(account_id)
        REFERENCES account(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "transaction";
-- +goose StatementEnd
