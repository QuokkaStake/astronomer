package types

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Config struct {
	DatabaseConfig DatabaseConfig `toml:"database"`
	LogConfig      LogConfig      `toml:"log"`
	TelegramConfig TelegramConfig `toml:"telegram"`
}

type TelegramConfig struct {
	Chat   int64   `toml:"chat"`
	Token  string  `toml:"token"`
	Admins []int64 `default:"[]" toml:"admins"`
}

type DiscordConfig struct {
	Guild   string `toml:"guild"`
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

func (c *Config) Validate() error {
	if err := c.DatabaseConfig.Validate(); err != nil {
		return fmt.Errorf("database config is invalid: %s", err)
	}
	return nil
}

func (c *Config) DisplayWarnings() []Warning {
	warnings := make([]Warning, 0)

	return warnings
}

func (c *Config) LogWarnings(logger *zerolog.Logger, warnings []Warning) {
	for _, warning := range warnings {
		entry := logger.Warn()

		for key, label := range warning.Labels {
			entry = entry.Str(key, label)
		}

		entry.Msg(warning.Message)
	}
}
