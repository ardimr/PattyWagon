-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_detail_id BIGINT NOT NULL REFERENCES order_details(id) ON DELETE CASCADE,
    item_id BIGINT NOT NULL REFERENCES items(id),
    item_name VARCHAR(255),
    product_category VARCHAR(255),
    item_image_url VARCHAR(255),
    price_per_item NUMERIC(12,2) NOT NULL,
    quantity INTEGER NOT NULL,
    total_price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT order_items_quantity_positive CHECK (quantity > 0)
);

CREATE INDEX idx_order_items_order_detail_id ON order_items(order_detail_id);
CREATE INDEX idx_order_items_item_id ON order_items(item_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_order_items_item_id;
DROP INDEX IF EXISTS idx_order_items_order_detail_id;
DROP TABLE IF EXISTS order_items;
SELECT 'down SQL query';
-- +goose StatementEnd
