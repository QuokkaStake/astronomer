package datafetcher

import (
	"main/pkg/types"
)

func (f *DataFetcher) DoesValidatorExist(chain *types.Chain, address string) error {
	rpc := f.GetRPC(chain)

	_, _, err := rpc.GetValidator(address)
	if err != nil {
		return err
	}

	return nil
}
