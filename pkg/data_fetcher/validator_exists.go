package datafetcher

import (
	"main/pkg/types"
)

func (f *DataFetcher) DoesValidatorExist(chain *types.Chain, address string) (*types.Validator, error) {
	rpc := f.GetRPC(chain)

	validator, _, err := rpc.GetValidator(address)
	if err != nil {
		return nil, err
	}

	return &validator.Validator, nil
}
