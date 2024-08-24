-- +goose Up
ALTER TABLE wallet_links DROP CONSTRAINT wallet_links_pkey;
ALTER TABLE wallet_links ADD CONSTRAINT wallet_links_pkey PRIMARY KEY (chain, reporter, user_id, address);

-- +goose Down
ALTER TABLE wallet_links DROP CONSTRAINT wallet_links_pkey;
ALTER TABLE wallet_links ADD CONSTRAINT wallet_links_pkey PRIMARY KEY (chain, reporter, user_id);
