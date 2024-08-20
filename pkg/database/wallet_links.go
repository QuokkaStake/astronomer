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

func (d *Database) DeleteWalletLink(
	chain string,
	reporter string,
	address string,
	userID string,
) (bool, error) {
	result, err := d.client.Exec(
		"DELETE FROM wallet_links WHERE chain = $1 AND reporter = $2 AND address = $3 AND user_id = $4",
		chain,
		reporter,
		address,
		userID,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete wallet link")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}
