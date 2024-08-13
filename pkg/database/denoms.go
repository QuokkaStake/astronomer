package database

import "main/pkg/types"

func (d *Database) InsertDenom(denom *types.Denom) error {
	_, err := d.client.Exec(
		"INSERT INTO denoms (chain, denom, display_denom, denom_exponent, coingecko_currency) VALUES ($1, $2, $3, $4, $5)",
		denom.Chain,
		denom.Denom,
		denom.DisplayDenom,
		denom.DenomExponent,
		denom.CoingeckoCurrency,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert denom")
		return err
	}

	return nil
}

func (d *Database) DeleteDenom(chainName, denom string) (bool, error) {
	result, err := d.client.Exec("DELETE FROM denoms WHERE chain = $1 AND denom = $2", chainName, denom)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete denom")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}
