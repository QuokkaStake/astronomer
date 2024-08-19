package datafetcher

import (
	"fmt"
	"main/pkg/constants"
	priceFetcher "main/pkg/price_fetcher"
	"main/pkg/types"
	"main/pkg/utils"
	regularMath "math"
	"sync"

	"cosmossdk.io/math"
)

func (f *DataFetcher) GetDenomCacheKey(chain, denom string) string {
	return fmt.Sprintf("denom_%s_%s", chain, denom)
}

func (f *DataFetcher) PopulateDenoms(amounts []*types.AmountWithChain) {
	chainWithDenoms := utils.Map(amounts, func(a *types.AmountWithChain) types.ChainWithDenom {
		return types.ChainWithDenom{
			Chain: a.Chain,
			Denom: a.Amount.Denom,
		}
	})

	foundDenoms, err := f.Database.FindDenoms(chainWithDenoms)
	if err != nil {
		f.Logger.Err(err).Msg("Could not fetch denoms")
		return
	}

	denomsByPriceFetcher := utils.GroupBy(foundDenoms, func(d *types.Denom) []constants.PriceFetcherName {
		if d.CoingeckoCurrency.IsZero() {
			return []constants.PriceFetcherName{}
		}

		return []constants.PriceFetcherName{constants.PriceFetcherNameCoingecko}
	})

	var wg sync.WaitGroup
	var mutex sync.Mutex
	prices := map[constants.PriceFetcherName]priceFetcher.Prices{}

	for priceFetcherName, denoms := range denomsByPriceFetcher {
		wg.Add(1)

		go func(priceFetcherName constants.PriceFetcherName, denoms []*types.Denom) {
			defer wg.Done()

			notCachedDenoms := []*types.Denom{}
			allPrices := priceFetcher.Prices{}

			for _, denom := range denoms {
				value, cached := f.Cache.Get(f.GetDenomCacheKey(denom.Chain, denom.Denom))
				if !cached {
					notCachedDenoms = append(notCachedDenoms, denom)
					continue
				}

				valueFloat, _ := value.(float64)
				allPrices.Set(denom.Chain, denom.Denom, valueFloat)
			}

			if len(notCachedDenoms) == 0 {
				f.Logger.Debug().
					Str("price_fetcher", string(priceFetcherName)).
					Msg("All denoms prices are cached, not fetching")
				return
			}

			foundPriceFetcher, ok := f.PriceFetchers[priceFetcherName]
			if !ok {
				mutex.Lock()
				prices[priceFetcherName] = allPrices
				mutex.Unlock()
				return
			}

			if fetcherPrices, denomFetchError := foundPriceFetcher.GetPrices(notCachedDenoms); denomFetchError != nil {
				f.Logger.Err(denomFetchError).
					Str("price_fetcher", string(priceFetcherName)).
					Msg("Could not fetch prices")
			} else {
				for chain, chainPrices := range fetcherPrices {
					for denom, value := range chainPrices {
						f.Cache.Set(f.GetDenomCacheKey(chain, denom), value)
						allPrices.Set(chain, denom, value)
					}
				}
			}

			mutex.Lock()
			prices[priceFetcherName] = allPrices
			mutex.Unlock()
		}(priceFetcherName, denoms)
	}

	wg.Wait()

	denomsMap := foundDenoms.ToMap()

	for _, amount := range amounts {
		chainDenoms, chainFound := denomsMap[amount.Chain]
		if !chainFound {
			continue
		}

		denom, denomFound := chainDenoms[amount.Amount.Denom]
		if !denomFound {
			continue
		}

		power := int64(regularMath.Pow10(denom.DenomExponent))

		amount.Amount.BaseDenom = amount.Amount.Denom
		amount.Amount.Denom = denom.DisplayDenom
		amount.Amount.Amount = amount.Amount.Amount.Quo(math.LegacyNewDec(power))

		for _, priceFetcherPrices := range prices {
			if price, found := priceFetcherPrices.Get(amount.Chain, amount.Amount.BaseDenom); found {
				singleTokenPrice := math.LegacyMustNewDecFromStr(fmt.Sprintf("%.6f", price))
				amountUSDPrice := amount.Amount.Amount.Mul(singleTokenPrice)
				amount.Amount.PriceUSD = &amountUSDPrice
			}
		}
	}
}
