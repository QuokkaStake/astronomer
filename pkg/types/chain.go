package types

import (
	"fmt"
)

type Chain struct {
	Name           string    `toml:"name"`
	PrettyName     string    `toml:"pretty-name"`
	KeplrName      string    `toml:"keplr-name"`
	LCDEndpoints   []string  `toml:"lcd-endpoints"`
	ProposalsType  string    `default:"v1beta1"          toml:"proposals-type"`
	MintscanPrefix string    `toml:"mintscan-prefix"`
	PingPrefix     string    `toml:"ping-prefix"`
	PingHost       string    `default:"https://ping.pub" toml:"ping-host"`
	Explorer       *Explorer `toml:"explorer"`

	Type string `default:"cosmos" toml:"type"`
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty chain name")
	}

	if len(c.LCDEndpoints) == 0 {
		return fmt.Errorf("no LCD endpoints provided")
	}

	return nil
}

func (c *Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}

func (c *Chain) GetExplorer() *Explorer {
	if c.MintscanPrefix != "" {
		return &Explorer{
			ProposalLinkPattern: fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", c.MintscanPrefix),
			WalletLinkPattern:   fmt.Sprintf("https://mintscan.io/%s/account/%%s", c.MintscanPrefix),
		}
	}

	if c.PingPrefix != "" {
		return &Explorer{
			ProposalLinkPattern: fmt.Sprintf("%s/%s/gov/%%s", c.PingHost, c.PingPrefix),
			WalletLinkPattern:   fmt.Sprintf("%s/%s/account/%%s", c.PingHost, c.PingPrefix),
		}
	}

	return c.Explorer
}

func (c *Chain) DisplayWarnings() []Warning {
	warnings := make([]Warning, 0)

	if c.Explorer == nil {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{"chain": c.Name},
			Message: "explorer is not set, cannot generate links",
		})
	} else {
		warnings = append(warnings, c.Explorer.DisplayWarnings(c.Name)...)
	}

	if c.KeplrName == "" {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{"chain": c.Name},
			Message: "keplr-name is not set, cannot generate Keplr link to proposal",
		})
	}

	return warnings
}
