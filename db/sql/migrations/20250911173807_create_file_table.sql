-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
  id BIGSERIAL PRIMARY KEY,
  uri VARCHAR(255),
  thumbnail_uri VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files;
-- +goose StatementEnd
