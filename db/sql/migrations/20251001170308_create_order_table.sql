-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_estimation_id BIGINT NOT NULL,
    is_purchased BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_user_purchased ON orders(user_id, is_purchased);
CREATE INDEX idx_orders_created_at ON orders(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_user_purchased;
DROP INDEX IF EXISTS idx_orders_created_at;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd