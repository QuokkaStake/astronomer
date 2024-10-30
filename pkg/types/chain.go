package types

import (
	"fmt"
)

type Chain struct {
	Name                  string `toml:"name"`
	PrettyName            string `toml:"pretty-name"`
	BaseDenom             string `toml:"base-denom"`
	Bech32ValidatorPrefix string
}

type ChainWithLCD struct {
	Chain       Chain
	LCDEndpoint string
}

func ChainFromArgs(args map[string]string) *ChainWithLCD {
	chain := &ChainWithLCD{
		Chain: Chain{},
	}

	for key, value := range args {
		switch key {
		case "name":
			chain.Chain.Name = value
		case "lcd-endpoint":
			chain.LCDEndpoint = value
		case "pretty-name":
			chain.Chain.PrettyName = value
		case "base-denom":
			chain.Chain.BaseDenom = value
		case "bech32-validator-prefix":
			chain.Chain.Bech32ValidatorPrefix = value
		}
	}

	return chain
}

func (c *Chain) UpdateFromArgs(args map[string]string) {
	for key, value := range args {
		switch key {
		case "pretty-name":
			c.PrettyName = value
		case "base-denom":
			c.BaseDenom = value
		case "bech32-validator-prefix":
			c.Bech32ValidatorPrefix = value
		}
	}
}

func (c *ChainWithLCD) Validate() error {
	if c.LCDEndpoint == "" {
		return fmt.Errorf("empty LCD endpoint")
	}

	if err := c.Chain.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty chain name")
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
