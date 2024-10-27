-- +goose Up
-- +goose StatementBegin
CREATE TABLE lcd (
     chain TEXT NOT NULL REFERENCES chains(name),
     host TEXT NOT NULL,
     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
     PRIMARY KEY (chain, host)
);
INSERT INTO lcd (chain, host) (SELECT name chain, lcd_endpoint host FROM chains);
-- +goose StatementEnd

-- +goose Down
DROP TABLE lcd;
