package database

import "main/pkg/types"

func (d *Database) InsertWalletLink(link *types.WalletLink) error {
	_, err := d.client.Exec(
		"INSERT INTO wallet_links (chain, reporter, user_id, address, alias) VALUES ($1, $2, $3, $4, $5)",
		link.Chain,
		link.Reporter,
		link.UserID,
		link.Address,
		link.Alias,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert wallet link")
		return err
	}

	return nil
}
