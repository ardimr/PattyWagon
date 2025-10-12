-- +goose Up
-- +goose StatementBegin
-- ALTER TABLE items ADD COLUMN price INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- ALTER TABLE items DROP COLUMN IF EXISTS price;

-- +goose StatementEnd
