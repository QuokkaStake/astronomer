package datafetcher

import (
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"main/pkg/types"
)

func (f *DataFetcher) DoesValidatorExist(chain *types.Chain, address string) (*stakingTypes.Validator, error) {
	rpc := f.GetRPC(chain)

	validator, _, err := rpc.GetValidator(address)
	if err != nil {
		return nil, err
	}

	return &validator.Validator, nil
}
