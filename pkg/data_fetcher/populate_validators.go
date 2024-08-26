package datafetcher

import (
	"main/pkg/types"
	"sync"
)

func (f *DataFetcher) PopulateValidators(validators []*types.ValidatorAddressWithMoniker) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(len(validators))

	for _, validator := range validators {
		go func(validator *types.ValidatorAddressWithMoniker) {
			defer wg.Done()

			rpc := f.GetRPC(validator.Chain)
			validatorFromChain, _, err := rpc.GetValidator(validator.Address)
			if err != nil {
				f.Logger.Error().Err(err).Msg("Could not get validator from chain")
				return
			}

			mutex.Lock()
			validator.Moniker = validatorFromChain.Validator.Description.Moniker
			mutex.Unlock()
		}(validator)
	}

	wg.Wait()
}
