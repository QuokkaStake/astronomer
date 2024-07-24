CREATE TABLE IF NOT EXISTS chain_binds (
    id SERIAL PRIMARY KEY,
    reporter TEXT NOT NULL,
    chat_id TEXT NOT NULL,
    chain TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (reporter, chat_id, chain)
);