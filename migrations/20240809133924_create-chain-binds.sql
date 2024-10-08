-- +goose Up
CREATE TABLE chain_binds (
    id SERIAL PRIMARY KEY,
    reporter TEXT NOT NULL,
    chat_id TEXT NOT NULL,
    chat_name TEXT NOT NULL,
    chain TEXT NOT NULL REFERENCES chains(name),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (reporter, chat_id, chain)
);

-- +goose Down
DROP TABLE chain_binds;
