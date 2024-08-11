package database

import (
	"main/pkg/types"

	"github.com/lib/pq"
)

func (d *Database) InsertExplorer(explorer *types.Explorer) error {
	_, err := d.client.Exec(
		"INSERT INTO explorers (chain, name, proposal_link_pattern, wallet_link_pattern, validator_link_pattern, main_link) VALUES ($1, $2, $3, $4, $5, $6)",
		explorer.Chain,
		explorer.Name,
		explorer.ProposalLinkPattern,
		explorer.WalletLinkPattern,
		explorer.ValidatorLinkPattern,
		explorer.MainLink,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert explorer")
		return err
	}

	return nil
}

func (d *Database) GetExplorersByChains(chains []string) (types.Explorers, error) {
	explorers := make(types.Explorers, 0)

	rows, err := d.client.Query(
		"SELECT chain, name, proposal_link_pattern, wallet_link_pattern, validator_link_pattern, main_link FROM explorers WHERE chain = any($1)",
		pq.Array(chains),
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Error getting explorers by names")
		return explorers, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	for rows.Next() {
		explorer := &types.Explorer{}

		err = rows.Scan(
			&explorer.Chain,
			&explorer.Name,
			&explorer.ProposalLinkPattern,
			&explorer.WalletLinkPattern,
			&explorer.ValidatorLinkPattern,
			&explorer.MainLink,
		)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting chain bind")
			return explorers, err
		}

		explorers = append(explorers, explorer)
	}

	return explorers, nil
}

func (d *Database) UpdateExplorer(explorer *types.Explorer) (bool, error) {
	result, err := d.client.Exec(
		"UPDATE explorers SET proposal_link_pattern = $1, wallet_link_pattern = $2, validator_link_pattern = $3, main_link = $4 WHERE name = $5 AND chain = $6",
		explorer.ProposalLinkPattern,
		explorer.WalletLinkPattern,
		explorer.ValidatorLinkPattern,
		explorer.MainLink,
		explorer.Name,
		explorer.Chain,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not update explorer")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

func (d *Database) DeleteExplorer(chainName, explorerName string) (bool, error) {
	result, err := d.client.Exec("DELETE FROM explorer WHERE chain = $1 AND name = $2", chainName, explorerName)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete explorer")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}
