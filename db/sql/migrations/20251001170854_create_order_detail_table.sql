-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_details (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    merchant_name VARCHAR(255),
    merchant_category VARCHAR(255),
    merchant_image_url VARCHAR(255),
    merchant_latitude DECIMAL(12, 6),
    merchant_longitude DECIMAL(12, 6),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(order_id, merchant_id)
);

CREATE INDEX idx_order_details_order_id ON order_details(order_id);
CREATE INDEX idx_order_details_merchant_id ON order_details(merchant_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
