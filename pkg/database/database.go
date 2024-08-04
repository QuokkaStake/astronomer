package database

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	migrationsPkg "main/migrations"
	"main/pkg/types"
	"strings"

	"github.com/rs/zerolog"

	_ "github.com/lib/pq"
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
	fmt.Printf("insert: %s %s %s %s\n", reporter, chatID, chatName, chain)

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
