-- +goose Up
ALTER TABLE chains ADD COLUMN base_denom TEXT NOT NULL;

-- +goose Down
ALTER TABLE chains DROP COLUMN base_denom;
