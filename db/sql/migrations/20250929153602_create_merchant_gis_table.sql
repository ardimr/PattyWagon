-- +goose Up
-- +goose StatementBegin
CREATE TABLE merchant_locations (
  id BIGSERIAL PRIMARY KEY,
  merchant_id BIGINT NOT NULL,
  h3_index BIGINT NOT NULL,
  resolution SMALLINT NOT NULL,
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
DROP TABLE IF EXISTS merchant_locations;
-- +goose StatementEnd
