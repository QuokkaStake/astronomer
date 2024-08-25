-- +goose Up
CREATE TABLE queries (
    id SERIAL PRIMARY KEY,
    reporter TEXT NOT NULL,
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    chat_id TEXT NOT NULL,
    command TEXT NOT NULL,
    query TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE queries;
