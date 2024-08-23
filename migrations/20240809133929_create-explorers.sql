-- +goose Up
CREATE TABLE explorers (
    chain TEXT NOT NULL REFERENCES chains(name),
    name TEXT NOT NULL,
    proposal_link_pattern TEXT NOT NULL,
    wallet_link_pattern TEXT NOT NULL,
    validator_link_pattern TEXT NOT NULL,
    main_link TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain, name)
);

-- +goose Down
DROP TABLE explorers;
