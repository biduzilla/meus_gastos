-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name VARCHAR(500) NOT NULL,
    type INTEGER NOT NULL CHECK (type IN (1, 2)),
    color VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_categories_user_id ON categories(user_id) WHERE NOT deleted;
CREATE INDEX idx_categories_name ON categories(name) WHERE NOT deleted;
CREATE INDEX idx_categories_type ON categories(type) WHERE NOT deleted;
CREATE INDEX idx_categories_user_type ON categories(user_id, type) WHERE NOT deleted;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd