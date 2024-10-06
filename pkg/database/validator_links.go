package database

import (
	"main/pkg/types"
)

func (d *Database) InsertValidatorLink(link *types.ValidatorLink) error {
	_, err := d.client.Exec(
		"INSERT INTO validator_links (chain, reporter, user_id, address) VALUES ($1, $2, $3, $4)",
		link.Chain,
		link.Reporter,
		link.UserID,
		link.Address,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert validator link")
		return err
	}

	return nil
}

func (d *Database) DeleteValidatorLink(
	chain string,
	reporter string,
	address string,
	userID string,
) (bool, error) {
	result, err := d.client.Exec(
		"DELETE FROM validator_links WHERE chain = $1 AND reporter = $2 AND address = $3 AND user_id = $4",
		chain,
		reporter,
		address,
		userID,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete validator link")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

func (d *Database) FindValidatorLinksByUserAndReporter(userID, reporter string) ([]*types.ValidatorLink, error) {
	validatorLinks := make([]*types.ValidatorLink, 0)

	rows, err := d.client.Query(
		"SELECT chain, reporter, user_id, address FROM validator_links WHERE user_id = $1 AND reporter = $2",
		userID,
		reporter,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Error getting validator links")
		return validatorLinks, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	for rows.Next() {
		validatorLink := &types.ValidatorLink{}

		err = rows.Scan(&validatorLink.Chain, &validatorLink.Reporter, &validatorLink.UserID, &validatorLink.Address)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting validator link")
			return validatorLinks, err
		}

		validatorLinks = append(validatorLinks, validatorLink)
	}

	return validatorLinks, nil
}
