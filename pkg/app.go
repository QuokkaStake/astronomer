package pkg

import (
	"main/pkg/fs"
	"main/pkg/logger"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type App struct {
	Logger *zerolog.Logger
	Config *types.Config

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

	return &App{
		Logger:      log,
		Config:      config,
		StopChannel: make(chan bool),
	}
}

func (a *App) Start() {
	a.Logger.Info().Str("interval", a.Config.Interval).Msg("Scheduled proposals reporting")

	<-a.StopChannel
	a.Logger.Info().Msg("Shutting down...")
}

func (a *App) Stop() {
	a.StopChannel <- true
}
