CREATE TABLE IF NOT EXISTS explorers (
    chain TEXT NOT NULL REFERENCES chains(name) ON DELETE CASCADE,
    name TEXT NOT NULL,
    proposal_link_pattern TEXT NOT NULL,
    wallet_link_pattern TEXT NOT NULL,
    validator_link_pattern TEXT NOT NULL,
    main_link TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain, name)
);