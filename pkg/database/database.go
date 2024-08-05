package database

import (
	"database/sql"
	migrationsPkg "main/migrations"
	"main/pkg/types"
	"strings"

	"github.com/lib/pq"

	"github.com/rs/zerolog"
)

type Database struct {
	logger zerolog.Logger
	config types.DatabaseConfig
	client *sql.DB
}

func NewDatabase(
	logger *zerolog.Logger,
	config types.DatabaseConfig,
) *Database {
	return &Database{
		logger: logger.With().Str("component", "state_manager").Logger(),
		config: config,
	}
}

func (d *Database) Init() {
	var db = d.InitPostgresDatabase()
	migrations, err := migrationsPkg.EmbedFS.ReadDir(".")
	if err != nil {
		d.logger.Panic().
			Err(err).
			Msg("Error reading migrations dir")
	}

	for _, entry := range migrations {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		d.logger.Info().
			Str("name", entry.Name()).
			Msg("Applying sqlite migration")

		content, err := migrationsPkg.EmbedFS.ReadFile(entry.Name())
		if err != nil {
			d.logger.Fatal().
				Str("name", entry.Name()).
				Err(err).
				Msg("Could not read migration content")
		}

		statement, err := db.Prepare(string(content))
		if err != nil {
			d.logger.Fatal().
				Str("name", entry.Name()).
				Err(err).
				Msg("Could not prepare migration")
		}
		if _, err := statement.Exec(); err != nil {
			d.logger.Fatal().
				Str("name", entry.Name()).
				Err(err).
				Msg("Could not execute migration")
		}

		_ = statement.Close() //nolint:sqlclosecheck // false positive
	}

	d.client = db
}

func (d *Database) InitPostgresDatabase() *sql.DB {
	db, err := sql.Open("postgres", d.config.Path)

	if err != nil {
		d.logger.Panic().Err(err).Msg("Could not open PostgreSQL database")
	}

	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)

	if err != nil {
		d.logger.Panic().Err(err).Msg("Could not query PostgreSQL database")
	}

	d.logger.Info().
		Str("version", version).
		Str("path", d.config.Path).
		Msg("PostgreSQL database connected")

	return db
}

func (d *Database) GetAllChainBinds(chatID string) ([]string, error) {
	chains := make([]string, 0)

	rows, err := d.client.Query(
		"SELECT chain FROM chain_binds WHERE chat_id = $1",
		chatID,
	)
	if err != nil {
		d.logger.Error().Str("chat", chatID).Err(err).Msg("Error getting chain binds")
		return chains, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	for rows.Next() {
		var chain string

		err = rows.Scan(&chain)
		if err != nil {
			d.logger.Error().Str("chat", chatID).Err(err).Msg("Error getting chain bind")
			return chains, err
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

func (d *Database) InsertChainBind(
	reporter string,
	chatID string,
	chatName string,
	chain string,
) error {
	_, err := d.client.Exec(
		"INSERT INTO chain_binds (reporter, chat_id, chat_name, chain) VALUES ($1, $2, $3, $4)",
		reporter,
		chatID,
		chatName,
		chain,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert chain bind")
		return err
	}

	return nil
}

func (d *Database) InsertChain(chain *types.Chain) error {
	_, err := d.client.Exec(
		"INSERT INTO chains (name, pretty_name, lcd_endpoint) VALUES ($1, $2, $3)",
		chain.Name,
		chain.PrettyName,
		chain.LCDEndpoint,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert chain")
		return err
	}

	return nil
}

func (d *Database) UpdateChain(chain *types.Chain) (bool, error) {
	result, err := d.client.Exec(
		"UPDATE chains SET pretty_name = $1, lcd_endpoint = $2 WHERE name = $3",
		chain.PrettyName,
		chain.LCDEndpoint,
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

func (d *Database) DeleteChainBind(
	reporter string,
	chatID string,
	chain string,
) (bool, error) {
	result, err := d.client.Exec(
		"DELETE FROM chain_binds WHERE reporter = $1 AND chat_id = $2 AND chain = $3",
		reporter,
		chatID,
		chain,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete chain bind")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

func (d *Database) GetChainsByNames(names []string) ([]*types.Chain, error) {
	chains := make([]*types.Chain, 0)

	rows, err := d.client.Query(
		"SELECT name, pretty_name, lcd_endpoint FROM chains WHERE name = any($1)",
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

		err = rows.Scan(&chain.Name, &chain.PrettyName, &chain.LCDEndpoint)
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

	rows, err := d.client.Query("SELECT name, pretty_name, lcd_endpoint FROM chains")
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

		err = rows.Scan(&chain.Name, &chain.PrettyName, &chain.LCDEndpoint)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting chain")
			return chains, err
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

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
