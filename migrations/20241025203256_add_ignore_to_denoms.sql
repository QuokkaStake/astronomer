-- +goose Up
ALTER TABLE denoms ADD COLUMN ignored BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE denoms DROP COLUMN ignored;
