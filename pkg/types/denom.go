package types

import (
	"fmt"
	"strconv"

	"github.com/creasty/defaults"
	"github.com/guregu/null/v5"
)

type Denom struct {
	Chain             string
	Denom             string
	DisplayDenom      string
	DenomExponent     int `default:"6"`
	CoingeckoCurrency null.String
}

func (d *Denom) Validate() error {
	if d.Chain == "" {
		return fmt.Errorf("chain is required")
	}

	if d.Denom == "" {
		return fmt.Errorf("denom is required")
	}

	if d.DisplayDenom == "" {
		return fmt.Errorf("display-denom is required")
	}

	if d.DenomExponent < 1 {
		return fmt.Errorf("denom-exponent must be positive")
	}

	return nil
}

func DenomFromArgs(args map[string]string) *Denom {
	denom := &Denom{}

	for key, value := range args {
		switch key {
		case "denom":
			denom.Denom = value
		case "chain":
			denom.Chain = value
		case "display-denom":
			denom.DisplayDenom = value
		case "display_denom":
			denom.DisplayDenom = value
		case "denom-exponent":
			fallthrough
		case "denom_exponent":
			if exponent, err := strconv.Atoi(value); err == nil {
				denom.DenomExponent = exponent
			}
		case "coingecko-currency":
			denom.CoingeckoCurrency = null.StringFrom(value)
		case "coingecko_currency":
			denom.CoingeckoCurrency = null.StringFrom(value)
		}
	}

	defaults.MustSet(denom)
	return denom
}

type Denoms []*Denom

func (denoms Denoms) ToMap() map[string]map[string]*Denom {
	m := make(map[string]map[string]*Denom)

	for _, denom := range denoms {
		if _, ok := m[denom.Chain]; !ok {
			m[denom.Chain] = make(map[string]*Denom)
		}

		m[denom.Chain][denom.Denom] = denom
	}

	return m
}

type ChainWithDenom struct {
	Chain string
	Denom string
}

type AmountWithChain struct {
	Chain  string
	Amount *Amount
}
