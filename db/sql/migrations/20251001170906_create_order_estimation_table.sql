-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_estimations (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    total_price INTEGER NOT NULL,
    estimated_delivery_time_minutes INTEGER NOT NULL,
    total_distance_km DECIMAL(10, 2) NOT NULL
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_estimations;
-- +goose StatementEnd