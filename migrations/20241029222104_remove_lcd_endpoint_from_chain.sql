-- +goose Up
-- +goose StatementBegin
ALTER TABLE chains DROP COLUMN lcd_endpoint;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chains ADD COLUMN lcd_endpoint TEXT;
-- +goose StatementEnd
