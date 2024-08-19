package datafetcher

import (
	"fmt"
	"main/pkg/types"
	"main/pkg/utils"
	regularMath "math"

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

	prices, err := f.PriceFetchers[0].GetPrices(foundDenoms)
	if err != nil {
		f.Logger.Err(err).Msg("Could not fetch prices")
		return
	}

	denomsMap := foundDenoms.ToMap()

	for _, amount := range amounts {
		if chainDenoms, chainFound := denomsMap[amount.Chain]; chainFound {
			if denom, denomFound := chainDenoms[amount.Amount.Denom]; denomFound {
				power := int64(regularMath.Pow10(denom.DenomExponent))

				amount.Amount.BaseDenom = amount.Amount.Denom
				amount.Amount.Denom = denom.DisplayDenom
				amount.Amount.Amount = amount.Amount.Amount.Quo(math.LegacyNewDec(power))

				if price, found := prices.Get(amount.Chain, amount.Amount.BaseDenom); found {
					singleTokenPrice := math.LegacyMustNewDecFromStr(fmt.Sprintf("%.6f", price))
					amountUSDPrice := amount.Amount.Amount.Mul(singleTokenPrice)
					amount.Amount.PriceUSD = &amountUSDPrice
				}
			}
		}
	}
}
