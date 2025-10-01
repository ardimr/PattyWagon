-- +goose Up
-- +goose StatementBegin
CREATE TABLE merchants (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  category VARCHAR(255),
  image_url VARCHAR(255) NOT NULL,
  latitude NUMERIC(12,2) NOT NULL,
  longitude NUMERIC(12,2) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS merchants;
-- +goose StatementEnd
