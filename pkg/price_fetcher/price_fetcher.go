package pricefetcher

import "main/pkg/types"

type Prices map[string]map[string]float64

func (p Prices) Get(chain string, denom string) (float64, bool) {
	chainPrices, chainPricesOk := p[chain]
	if !chainPricesOk {
		return 0, false
	}

	denomPrice, denomPriceOk := chainPrices[denom]
	return denomPrice, denomPriceOk
}

type PriceFetcher interface {
	GetPrices(denomInfos []*types.Denom) (Prices, error)
	Name() string
}
