package database

import (
	"main/pkg/types"
)

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

func (d *Database) FindWalletLinksByUserAndReporter(userID, reporter string) ([]*types.WalletLink, error) {
	walletLinks := make([]*types.WalletLink, 0)

	rows, err := d.client.Query(
		"SELECT chain, reporter, user_id, address, alias FROM wallet_links WHERE user_id = $1 AND reporter = $2",
		userID,
		reporter,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Error getting wallet links")
		return walletLinks, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	for rows.Next() {
		walletLink := &types.WalletLink{}

		err = rows.Scan(&walletLink.Chain, &walletLink.Reporter, &walletLink.UserID, &walletLink.Address, &walletLink.Alias)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting wallet link")
			return walletLinks, err
		}

		walletLinks = append(walletLinks, walletLink)
	}

	return walletLinks, nil
}
