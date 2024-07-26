package database

import (
	"database/sql"
	"github.com/rs/zerolog"
	migrationsPkg "main/migrations"
	"main/pkg/types"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

type Database struct {
	logger zerolog.Logger
	config types.DatabaseConfig
	client *sql.DB
	mutex  sync.Mutex
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
	var db *sql.DB = d.InitPostgresDatabase()
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

		_ = statement.Close()
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
