-- +goose Up
-- +goose StatementBegin
ALTER TABLE chains ADD COLUMN bech32_validator_prefix TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chains DROP COLUMN bech32_validator_prefix;
-- +goose StatementEnd
