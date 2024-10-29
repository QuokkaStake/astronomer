package database

import (
	"database/sql"
	"errors"
	"main/pkg/constants"
	"main/pkg/types"
)

func (d *Database) GetRandomLCDHost(chain *types.Chain) (string, error) {
	host := ""

	row := d.client.QueryRow(
		"SELECT host FROM lcd WHERE chain = $1 LIMIT 1 ORDER BY RANDOM ()",
		chain.Name,
	)

	if err := row.Scan(&host); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", constants.ErrChainNotFound
		}

		d.logger.Error().Err(err).Msg("Error getting chain LCD")
		return "", err
	}

	return host, nil
}
