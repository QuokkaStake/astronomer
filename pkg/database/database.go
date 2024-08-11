package database

import (
	"database/sql"
	"errors"
	migrationsPkg "main/migrations"
	"main/pkg/types"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/rs/zerolog"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Logger struct {
	Logger zerolog.Logger
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Debug().Msgf(strings.TrimSpace(format), v...)
}

func (l *Logger) Verbose() bool {
	return true
}

type Database struct {
	logger         zerolog.Logger
	config         types.DatabaseConfig
	client         *sql.DB
	databaseLogger *Logger
	migrator       *migrate.Migrate
}

func NewDatabase(
	logger *zerolog.Logger,
	config types.DatabaseConfig,
) *Database {
	return &Database{
		logger: logger.With().Str("component", "database").Logger(),
		config: config,
		databaseLogger: &Logger{
			Logger: logger.With().Str("component", "migrations").Logger(),
		},
	}
}

func (d *Database) Init() {
	d.client = d.InitPostgresDatabase()

	filesystem, err := iofs.New(migrationsPkg.EmbedFS, "migrations")
	if err != nil {
		d.logger.Panic().Err(err).Msg("Could not init filesystem")
	}

	driver, err := postgres.WithInstance(d.client, &postgres.Config{})
	if err != nil {
		d.logger.Panic().Err(err).Msg("Could not init migrations driver")
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		filesystem,
		"postgres",
		driver,
	)
	if err != nil {
		d.logger.Panic().Err(err).Msg("Could not init migrate instance")
	}

	m.Log = d.databaseLogger
	d.migrator = m
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

func (d *Database) Migrate() {
	if err := d.migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			d.logger.Debug().Msg("No pending migrations.")
			return
		}
		d.logger.Panic().Err(err).Msg("Could not run migrations")
	}
}

func (d *Database) Rollback() {
	if err := d.migrator.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			d.logger.Debug().Msg("No migrations to rollback.")
			return
		}
		d.logger.Panic().Err(err).Msg("Could not rollback migrations")
	}
}
