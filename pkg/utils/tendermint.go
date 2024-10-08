package utils

import (
	"sort"

	"cosmossdk.io/math"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func GetTotalVP(validators []stakingTypes.Validator) math.LegacyDec {
	vp := math.LegacyNewDec(0)

	for _, validator := range validators {
		if validator.Status == stakingTypes.Bonded {
			vp = vp.Add(validator.DelegatorShares)
		}
	}

	return vp
}

func FindValidatorRank(validators []stakingTypes.Validator, valoper string) int {
	valsActive := Filter(validators, func(v stakingTypes.Validator) bool {
		return v.Status == stakingTypes.Bonded
	})

	sort.SliceStable(valsActive, func(i, j int) bool {
		return valsActive[i].DelegatorShares.GT(valsActive[j].DelegatorShares)
	})

	for index, validator := range valsActive {
		if validator.OperatorAddress == valoper {
			return index + 1
		}
	}

	return 0
}
