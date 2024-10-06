-- +goose Up
CREATE TABLE validator_links (
    chain TEXT NOT NULL REFERENCES chains(name),
    reporter TEXT NOT NULL,
    user_id TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain, reporter, user_id, address)
);

-- +goose Down
DROP TABLE validator_links;
