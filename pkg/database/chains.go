package database

import (
	"main/pkg/types"

	"github.com/lib/pq"
)

func (d *Database) GetChainsByNames(names []string) ([]*types.Chain, error) {
	chains := make([]*types.Chain, 0)

	rows, err := d.client.Query(
		"SELECT name, pretty_name, lcd_endpoint, base_denom, bech32_validator_prefix FROM chains WHERE name = any($1)",
		pq.Array(names),
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Error getting chains by names")
		return chains, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	for rows.Next() {
		chain := &types.Chain{}

		err = rows.Scan(&chain.Name, &chain.PrettyName, &chain.LCDEndpoint, &chain.BaseDenom, &chain.Bech32ValidatorPrefix)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting chain bind")
			return chains, err
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

func (d *Database) GetAllChains() ([]*types.Chain, error) {
	chains := make([]*types.Chain, 0)

	rows, err := d.client.Query("SELECT name, pretty_name, lcd_endpoint, base_denom, bech32_validator_prefix FROM chains")
	if err != nil {
		d.logger.Error().Err(err).Msg("Error getting all chains")
		return chains, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	for rows.Next() {
		chain := &types.Chain{}

		err = rows.Scan(&chain.Name, &chain.PrettyName, &chain.LCDEndpoint, &chain.BaseDenom, &chain.Bech32ValidatorPrefix)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting chain")
			return chains, err
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

func (d *Database) InsertChain(chain *types.Chain) error {
	_, err := d.client.Exec(
		"INSERT INTO chains (name, pretty_name, lcd_endpoint, base_denom, bech32_validator_prefix) VALUES ($1, $2, $3, $4, $5)",
		chain.Name,
		chain.PrettyName,
		chain.LCDEndpoint,
		chain.BaseDenom,
		chain.Bech32ValidatorPrefix,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert chain")
		return err
	}

	return nil
}

func (d *Database) UpdateChain(chain *types.Chain) (bool, error) {
	result, err := d.client.Exec(
		"UPDATE chains SET pretty_name = $1, lcd_endpoint = $2 , base_denom = $3, bech32_validator_prefix = $4 WHERE name = $5",
		chain.PrettyName,
		chain.LCDEndpoint,
		chain.BaseDenom,
		chain.Bech32ValidatorPrefix,
		chain.Name,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not update chain")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

func (d *Database) DeleteChain(chainName string) (bool, error) {
	result, err := d.client.Exec("DELETE FROM chains WHERE name = $1", chainName)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete chain")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}
