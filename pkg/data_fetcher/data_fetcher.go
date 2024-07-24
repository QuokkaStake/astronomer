package datafetcher

import (
	"main/pkg/tendermint"
	"main/pkg/types"
	"main/pkg/utils"
	"slices"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Chains types.Chains
	RPCs   []*tendermint.RPC
	Logger zerolog.Logger
}

func NewDataFetcher(config *types.Config, logger *zerolog.Logger) *DataFetcher {
	rpcs := make([]*tendermint.RPC, len(config.Chains))

	for index, chain := range config.Chains {
		rpcs[index] = tendermint.NewRPC(chain, 10, logger)
	}

	return &DataFetcher{
		Logger: logger.With().Str("component", "data_fetcher").Logger(),
		Chains: config.Chains, RPCs: rpcs,
	}
}

func (f *DataFetcher) FindValidator(query string, chains []string) map[string]types.ValidatorInfo {
	lowercaseQuery := strings.ToLower(query)

	var wg sync.WaitGroup
	var mutex sync.Mutex

	response := map[string]types.ValidatorInfo{}

	for index, chain := range f.Chains {
		if len(chains) > 0 && !slices.Contains(chains, chain.Name) {
			continue
		}

		rpc := f.RPCs[index]

		wg.Add(1)

		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			validators, _, err := rpc.GetAllValidators()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name] = types.ValidatorInfo{
					Chain:     chain,
					Validator: nil,
					Error:     err,
				}
				return
			}

			validator, _ := utils.Find(validators.Validators, func(v *types.Validator) bool {
				return strings.Contains(strings.ToLower(v.Description.Moniker), lowercaseQuery)
			})

			totalVP := validators.GetTotalVP()

			info := types.ValidatorInfo{
				Chain:     chain,
				Validator: validator,
				Error:     nil,
			}

			if validator != nil {
				info.VotingPowerPercent = validator.DelegatorShares.Quo(totalVP).MustFloat64()
				if validator.Active() {
					info.Rank = validators.FindValidatorRank(validator.OperatorAddress)
				}
			}

			response[chain.Name] = info
		}(chain, rpc)
	}

	wg.Wait()

	return response
}
