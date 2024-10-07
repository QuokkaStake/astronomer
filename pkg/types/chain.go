package types

import (
	"fmt"
)

type Chain struct {
	Name                  string `toml:"name"`
	PrettyName            string `toml:"pretty-name"`
	LCDEndpoint           string `toml:"lcd-endpoint"`
	BaseDenom             string `toml:"base-denom"`
	Bech32ValidatorPrefix string
}

func ChainFromArgs(args map[string]string) *Chain {
	chain := &Chain{}

	for key, value := range args {
		switch key {
		case "name":
			chain.Name = value
		case "lcd-endpoint":
			chain.LCDEndpoint = value
		case "pretty-name":
			chain.PrettyName = value
		case "base-denom":
			chain.BaseDenom = value
		case "bech32-validator-prefix":
			chain.Bech32ValidatorPrefix = value
		}
	}

	return chain
}

func (c *Chain) UpdateFromArgs(args map[string]string) {
	for key, value := range args {
		switch key {
		case "lcd-endpoint":
			c.LCDEndpoint = value
		case "pretty-name":
			c.PrettyName = value
		case "base-denom":
			c.BaseDenom = value
		case "bech32-validator-prefix":
			c.Bech32ValidatorPrefix = value
		}
	}
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty chain name")
	}

	if c.LCDEndpoint == "" {
		return fmt.Errorf("empty LCD endpoint")
	}

	if c.BaseDenom == "" {
		return fmt.Errorf("empty base denom")
	}

	if c.Bech32ValidatorPrefix == "" {
		return fmt.Errorf("empty bech32 validator prefix")
	}

	return nil
}

func (c *Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}
