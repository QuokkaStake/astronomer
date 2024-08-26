package pkg

import (
	datafetcher "main/pkg/data_fetcher"
	databasePkg "main/pkg/database"
	"main/pkg/fs"
	interacterPkg "main/pkg/interacter"
	"main/pkg/interacter/telegram"
	"main/pkg/logger"
	"main/pkg/metrics"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type App struct {
	Logger  *zerolog.Logger
	Config  *types.Config
	Version string

	Interacters    []interacterPkg.Interacter
	MetricsManager *metrics.Manager
	Database       *databasePkg.Database

	StopChannel chan bool
}

func NewApp(configPath string, filesystem fs.FS, version string) *App {
	config, err := GetConfig(filesystem, configPath)
	if err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Could not load config")
	}

	if err = config.Validate(); err != nil {
		logger.GetDefaultLogger().Panic().Err(err).Msg("Provided config is invalid!")
	}

	if warnings := config.DisplayWarnings(); len(warnings) > 0 {
		config.LogWarnings(logger.GetDefaultLogger(), warnings)
	} else {
		logger.GetDefaultLogger().Info().Msg("Provided config is valid.")
	}

	log := logger.GetLogger(config.LogConfig)
	database := databasePkg.NewDatabase(log, config.DatabaseConfig)
	metricsManager := metrics.NewManager(log, config.MetricsConfig)
	dataFetcher := datafetcher.NewDataFetcher(*log, database)
	interacters := []interacterPkg.Interacter{
		telegram.NewInteracter(config.TelegramConfig, version, log, dataFetcher, database, metricsManager),
	}

	return &App{
		Logger:         log,
		Config:         config,
		Version:        version,
		Interacters:    interacters,
		Database:       database,
		MetricsManager: metricsManager,
		StopChannel:    make(chan bool),
	}
}

func (a *App) Start() {
	a.Logger.Info().Msg("Listening")

	a.MetricsManager.LogAppVersion(a.Version)
	go a.MetricsManager.Start()

	a.Database.Init()

	for _, interacter := range a.Interacters {
		interacter.Init()

		a.MetricsManager.LogReporterEnabled(interacter.Name(), interacter.Enabled())

		if interacter.Enabled() {
			a.Logger.Info().Str("name", interacter.Name()).Msg("Interacter is enabled")
			go interacter.Start()
		} else {
			a.Logger.Info().Str("name", interacter.Name()).Msg("Interacter is disabled")
		}
	}

	<-a.StopChannel
	a.Logger.Info().Msg("Shutting down...")
}

func (a *App) Stop() {
	a.StopChannel <- true
}
