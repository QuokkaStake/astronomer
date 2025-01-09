package pricefetcher

import "main/pkg/types"

type Prices map[string]map[string]float64

func (p *Prices) Set(chain string, denom string, value float64) {
	if _, ok := (*p)[chain]; !ok {
		(*p)[chain] = make(map[string]float64)
	}

	(*p)[chain][denom] = value
}

func (p *Prices) Get(chain string, denom string) (float64, bool) {
	chainPrices, chainPricesOk := (*p)[chain]
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
