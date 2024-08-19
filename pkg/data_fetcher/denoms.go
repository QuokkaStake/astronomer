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

	denomsByPriceFetcher := utils.GroupBy(foundDenoms, func(d *types.Denom) constants.PriceFetcherName {
		if d.CoingeckoCurrency.IsZero() {
			return ""
		}

		return constants.PriceFetcherNameCoingecko
	})

	var wg sync.WaitGroup
	var mutex sync.Mutex
	prices := map[constants.PriceFetcherName]priceFetcher.Prices{}

	for priceFetcherName, denoms := range denomsByPriceFetcher {
		wg.Add(1)

		go func(priceFetcherName constants.PriceFetcherName, denoms []*types.Denom) {
			defer wg.Done()

			foundPriceFetcher, ok := f.PriceFetchers[priceFetcherName]
			if !ok {
				return
			}

			fetcherPrices, denomFetchError := foundPriceFetcher.GetPrices(denoms)
			if denomFetchError != nil {
				f.Logger.Err(denomFetchError).Msg("Could not fetch prices")
			} else {
				mutex.Lock()
				prices[priceFetcherName] = fetcherPrices
			}
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
