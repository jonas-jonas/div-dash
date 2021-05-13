-- +goose Up
-- +goose StatementBegin
CREATE TABLE account (
    id              text PRIMARY KEY,
    name            text NOT NULL,
    user_id         text NOT NULL,
    CONSTRAINT uq_id_user_id
      UNIQUE(id, user_id),
    CONSTRAINT fk_account_user
      FOREIGN KEY(user_id)
        REFERENCES "user"(id)
    
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE account;
-- +goose StatementEnd
