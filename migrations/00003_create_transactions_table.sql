-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    description VARCHAR(500) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    type INTEGER NOT NULL CHECK (type IN (1, 2)),
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id) WHERE NOT deleted;
CREATE INDEX idx_transactions_category_id ON transactions(category_id) WHERE NOT deleted;
CREATE INDEX idx_transactions_type ON transactions(type) WHERE NOT deleted;
CREATE INDEX idx_transactions_description_search ON transactions USING GIN (to_tsvector('simple', description));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
