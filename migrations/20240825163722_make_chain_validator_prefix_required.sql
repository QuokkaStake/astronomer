-- +goose Up
-- +goose StatementBegin
ALTER TABLE chains ALTER COLUMN bech32_validator_prefix SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chains ALTER COLUMN bech32_validator_prefix DROP NOT NULL;
-- +goose StatementEnd
