package main

import (
	"main/pkg"
	databasePkg "main/pkg/database"
	"main/pkg/fs"
	"main/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	version = "unknown"
)

func ExecuteMain(configPath string) {
	filesystem := &fs.OsFS{}
	app := pkg.NewApp(configPath, filesystem, version)
	app.Start()
}

func ExecuteValidateConfig(configPath string) {
	filesystem := &fs.OsFS{}

	config, err := pkg.GetConfig(filesystem, configPath)
	if err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Could not load config!")
	}

	if err := config.Validate(); err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Config is invalid!")
	}

	if warnings := config.DisplayWarnings(); len(warnings) > 0 {
		config.LogWarnings(logger.GetDefaultLogger(), warnings)
	} else {
		logger.GetDefaultLogger().Info().Msg("Provided config is valid.")
	}
}

func ExecuteMigrate(configPath string) {
	filesystem := &fs.OsFS{}

	config, err := pkg.GetConfig(filesystem, configPath)
	if err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Could not load config!")
	}

	if err := config.Validate(); err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Config is invalid!")
	}

	database := databasePkg.NewDatabase(logger.GetDefaultLogger(), config.DatabaseConfig)
	database.Init()
	database.Migrate()
}

func ExecuteRollback(configPath string) {
	filesystem := &fs.OsFS{}

	config, err := pkg.GetConfig(filesystem, configPath)
	if err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Could not load config!")
	}

	if err := config.Validate(); err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Config is invalid!")
	}

	database := databasePkg.NewDatabase(logger.GetDefaultLogger(), config.DatabaseConfig)
	database.Init()
	database.Rollback()
}

func main() {
	var ConfigPath string

	rootCmd := &cobra.Command{
		Use:     "astronomer --config [config path]",
		Long:    "A multi-chain explorer/wallet as a Telegram/Discord bot.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteMain(ConfigPath)
		},
	}

	validateConfigCmd := &cobra.Command{
		Use:     "validate-config --config [config path]",
		Long:    "Validate application config.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteValidateConfig(ConfigPath)
		},
	}

	migrateCmd := &cobra.Command{
		Use:     "migrate --config [config path]",
		Long:    "Perform a database migration.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteMigrate(ConfigPath)
		},
	}

	rollbackCmd := &cobra.Command{
		Use:     "rollback --config [config path]",
		Long:    "Rollback all database migrations.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteRollback(ConfigPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	_ = rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.AddCommand(validateConfigCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Could not start application")
	}
}
