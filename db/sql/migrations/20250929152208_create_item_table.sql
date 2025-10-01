-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
  id BIGSERIAL PRIMARY KEY,
  merchant_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  category VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_merchant
    FOREIGN KEY (merchant_id)
    REFERENCES merchants(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
