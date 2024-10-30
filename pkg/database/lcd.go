package database

import (
	"errors"
	"main/pkg/types"
)

func (d *Database) GetLCDHosts(chain *types.Chain) ([]string, error) {
	hosts := []string{}

	rows, err := d.client.Query(
		"SELECT host FROM lcd WHERE chain = $1",
		chain.Name,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Error getting LCDs for chain")
		return hosts, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		host := ""

		if scanErr := rows.Scan(&host); scanErr != nil {
			d.logger.Error().Err(scanErr).Msg("Error getting chain LCD")
			return hosts, err
		}

		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return hosts, errors.New("no LCD hosts found")
	}

	return hosts, nil
}

func (d *Database) InsertLCDHost(chain *types.Chain, host string) error {
	_, err := d.client.Exec(
		"INSERT INTO lcd (chain, host) VALUES ($1, $2)",
		chain.Name,
		host,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert LCD host")
		return err
	}

	return nil
}
