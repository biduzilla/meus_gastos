-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions DROP COLUMN IF EXISTS type;
DROP INDEX IF EXISTS idx_transactions_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions ADD COLUMN type INTEGER NOT NULL CHECK (type IN (1, 2));
CREATE INDEX idx_transactions_type ON transactions(type) WHERE NOT deleted;
-- +goose StatementEnd
