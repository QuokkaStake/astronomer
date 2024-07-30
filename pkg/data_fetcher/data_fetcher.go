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

func (f *DataFetcher) GetChainsParams(chains []string) map[string]*types.ChainParams {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	response := map[string]*types.ChainParams{}

	for index, chain := range f.Chains {
		if len(chains) > 0 && !slices.Contains(chains, chain.Name) {
			continue
		}

		response[chain.Name] = &types.ChainParams{
			Chain: chain,
		}

		rpc := f.RPCs[index]

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetStakingParams()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name].StakingParamsError = err
			} else {
				response[chain.Name].StakingParams = params.Params
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetSlashingParams()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name].SlashingParamsError = err
			} else {
				response[chain.Name].SlashingParams = params.Params
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetGovParams("voting")
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name].VotingParamsError = err
			} else {
				response[chain.Name].VotingParams = params.VotingParams
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetGovParams("deposit")
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name].DepositParamsError = err
			} else {
				response[chain.Name].DepositParams = params.DepositParams
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetGovParams("tallying")
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name].TallyParamsError = err
			} else {
				response[chain.Name].TallyParams = params.TallyParams
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			blockTime, err := rpc.GetBlockTime()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name].BlockTimeError = err
			} else {
				response[chain.Name].BlockTime = blockTime
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetMintParams()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name].MintParamsError = err
			} else {
				response[chain.Name].MintParams = params.Params
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			inflation, _, err := rpc.GetInflation()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				response[chain.Name].InflationError = err
			} else {
				response[chain.Name].Inflation = inflation.Inflation
			}
		}(chain, rpc)
	}

	wg.Wait()

	return response
}
