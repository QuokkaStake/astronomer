package types

import (
	"fmt"
)

type Chain struct {
	Name        string `toml:"name"`
	PrettyName  string `toml:"pretty-name"`
	LCDEndpoint string `toml:"lcd-endpoint"`
	BaseDenom   string `toml:"base-denom"`
}

func ChainFromArgs(args map[string]string) *Chain {
	chain := &Chain{}

	for key, value := range args {
		switch key {
		case "name":
			chain.Name = value
		case "lcd_endpoint":
			chain.LCDEndpoint = value
		case "lcd-endpoint":
			chain.LCDEndpoint = value
		case "pretty_name":
			chain.PrettyName = value
		case "pretty-name":
			chain.PrettyName = value
		case "base_denom":
			chain.BaseDenom = value
		case "base-denom":
			chain.BaseDenom = value
		}
	}

	return chain
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

	return nil
}

func (c *Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}
