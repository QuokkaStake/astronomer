package database

import (
	"context"
	"database/sql"
	"errors"
	"main/pkg/constants"
	"main/pkg/types"

	"github.com/lib/pq"
)

func (d *Database) GetChainsByNames(names []string) ([]*types.Chain, error) {
	chains := make([]*types.Chain, 0)

	rows, err := d.client.Query(
		"SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains WHERE name = any($1)",
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

		err = rows.Scan(&chain.Name, &chain.PrettyName, &chain.BaseDenom, &chain.Bech32ValidatorPrefix)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting chains by names")
			return chains, err
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

func (d *Database) GetChainByName(name string) (*types.Chain, error) {
	chain := &types.Chain{}
	row := d.client.QueryRow(
		"SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains WHERE name = $1 LIMIT 1",
		name,
	)

	err := row.Scan(
		&chain.Name,
		&chain.PrettyName,
		&chain.BaseDenom,
		&chain.Bech32ValidatorPrefix,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrChainNotFound
		}

		d.logger.Error().Err(err).Msg("Error getting chain by name")
		return nil, err
	}

	return chain, nil
}

func (d *Database) GetAllChains() ([]*types.Chain, error) {
	chains := make([]*types.Chain, 0)

	rows, err := d.client.Query("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains")
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

		err = rows.Scan(&chain.Name, &chain.PrettyName, &chain.BaseDenom, &chain.Bech32ValidatorPrefix)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting chain")
			return chains, err
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

func (d *Database) InsertChain(chain *types.ChainWithLCD) error {
	tx, err := d.client.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	_, err = tx.Exec(
		"INSERT INTO chains (name, pretty_name, base_denom, bech32_validator_prefix) VALUES ($1, $2, $3, $4)",
		chain.Chain.Name,
		chain.Chain.PrettyName,
		chain.Chain.BaseDenom,
		chain.Chain.Bech32ValidatorPrefix,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert chain")
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO lcd (chain, host) VALUES ($1, $2)",
		chain.Chain.Name,
		chain.LCDEndpoint,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert LCD")
		return err
	}

	if err = tx.Commit(); err != nil {
		d.logger.Error().Err(err).Msg("Error committing transaction when inserting chain")
		return err
	}

	return nil
}

func (d *Database) UpdateChain(chain *types.Chain) (bool, error) {
	result, err := d.client.Exec(
		"UPDATE chains SET pretty_name = $1, base_denom = $2, bech32_validator_prefix = $3 WHERE name = $4",
		chain.PrettyName,
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
	tx, err := d.client.BeginTx(context.Background(), nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback() //nolint:errcheck

	_, err = tx.Exec("DELETE FROM lcd WHERE chain = $1", chainName)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete LCD when deleting chains")
		return false, err
	}

	result, err := tx.Exec("DELETE FROM chains WHERE name = $1", chainName)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete chain")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()

	if err = tx.Commit(); err != nil {
		d.logger.Error().Err(err).Msg("Error committing transaction when inserting chain")
		return false, err
	}

	return rowsAffected > 0, nil
}
