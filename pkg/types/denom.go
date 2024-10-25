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
	Ignored           bool
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

func (d *Denom) PrintCoingeckoCurrency() string {
	if d.CoingeckoCurrency.IsZero() {
		return "not set"
	}

	return d.CoingeckoCurrency.String
}

func DenomFromArgs(args map[string]string) *Denom {
	denom := &Denom{}

	for key, value := range args {
		switch key {
		case "denom":
			denom.Denom = value
		case "chain":
			denom.Chain = value
		case "display-denom", "display_denom":
			denom.DisplayDenom = value
		case "ignore", "ignored":
			if ignored, err := strconv.ParseBool(value); err == nil {
				denom.Ignored = ignored
			}
		case "denom-exponent", "denom_exponent":
			if exponent, err := strconv.Atoi(value); err == nil {
				denom.DenomExponent = exponent
			}
		case "coingecko-currency", "coingecko_currency":
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

func (denoms Denoms) Find(amount *AmountWithChain) (*Denom, bool) {
	for _, denom := range denoms {
		if denom.Chain == amount.Chain && denom.Denom == amount.Amount.Denom {
			return denom, true
		}
	}

	return nil, false
}

type ChainWithDenom struct {
	Chain string
	Denom string
}

type AmountWithChain struct {
	Chain  string
	Amount *Amount
}
