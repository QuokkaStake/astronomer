-- +goose Up
CREATE TABLE chains (
    name TEXT NOT NULL,
    pretty_name TEXT NOT NULL,
    lcd_endpoint TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (name),
    PRIMARY KEY (name)
);

-- +goose Down
DROP TABLE chains;
