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
    account_id            text NOT NULL REFERENCES account,
    user_id               text NOT NULL REFERENCES "user",
    side                  text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "transaction";
-- +goose StatementEnd
