package types

import "github.com/guregu/null/v5"

type MetricsConfig struct {
	Enabled    null.Bool `default:"true"  toml:"enabled"`
	ListenAddr string    `default:":9590" toml:"listen-addr"`
}
