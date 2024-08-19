package database

import (
	"fmt"
	"main/pkg/types"
	"strings"
)

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

func (d *Database) FindDenoms(denoms []types.ChainWithDenom) (types.Denoms, error) {
	if len(denoms) == 0 {
		return []*types.Denom{}, nil
	}

	subqueries := make([]string, len(denoms))
	args := make([]interface{}, len(denoms)*2)

	for index, denom := range denoms {
		subqueries[index] = fmt.Sprintf("(chain = $%d AND denom = $%d)", index*2+1, index*2+2)
		args[index*2] = denom.Chain
		args[index*2+1] = denom.Denom
	}

	query := "SELECT chain, denom, display_denom, denom_exponent, coingecko_currency FROM denoms WHERE " + strings.Join(subqueries, " OR ")

	rows, err := d.client.Query(query, args...)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not find denoms")
		return []*types.Denom{}, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	returnDenoms := []*types.Denom{}

	for rows.Next() {
		denom := &types.Denom{}

		err = rows.Scan(&denom.Chain, &denom.Denom, &denom.DisplayDenom, &denom.DenomExponent, &denom.CoingeckoCurrency)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting denom")
			return returnDenoms, err
		}

		returnDenoms = append(returnDenoms, denom)
	}

	return returnDenoms, nil
}
