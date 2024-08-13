CREATE TABLE denoms (
    chain TEXT NOT NULL REFERENCES chains(name),
    denom TEXT NOT NULL,
    display_denom TEXT NOT NULL,
    denom_exponent INT NOT NULL,
    coingecko_currency TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain, denom)
);
