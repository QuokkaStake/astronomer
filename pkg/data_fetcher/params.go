package datafetcher

import (
	"main/pkg/tendermint"
	"main/pkg/types"
	"sync"
)

func (f *DataFetcher) GetChainsParams(chainNames []string) types.ChainsParams {
	response := types.ChainsParams{}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	chains, err := f.Database.GetChainsByNames(chainNames)
	if err != nil {
		response.Error = err
		return response
	}

	chainsParams := map[string]*types.ChainParams{}

	for _, chain := range chains {
		chainsParams[chain.Name] = &types.ChainParams{
			Chain: chain,
		}

		rpc := tendermint.NewRPC(chain, 10, f.Logger)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetStakingParams()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsParams[chain.Name].StakingParamsError = err
			} else {
				chainsParams[chain.Name].StakingParams = params.Params
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetSlashingParams()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsParams[chain.Name].SlashingParamsError = err
			} else {
				chainsParams[chain.Name].SlashingParams = params.Params
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetGovParams("voting")
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsParams[chain.Name].VotingParamsError = err
			} else {
				chainsParams[chain.Name].VotingParams = params.VotingParams
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetGovParams("deposit")
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsParams[chain.Name].DepositParamsError = err
			} else {
				chainsParams[chain.Name].DepositParams = params.DepositParams
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetGovParams("tallying")
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsParams[chain.Name].TallyParamsError = err
			} else {
				chainsParams[chain.Name].TallyParams = params.TallyParams
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			blockTime, err := rpc.GetBlockTime()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsParams[chain.Name].BlockTimeError = err
			} else {
				chainsParams[chain.Name].BlockTime = blockTime
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			params, _, err := rpc.GetMintParams()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsParams[chain.Name].MintParamsError = err
			} else {
				chainsParams[chain.Name].MintParams = params.Params
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			inflation, _, err := rpc.GetInflation()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsParams[chain.Name].InflationError = err
			} else {
				chainsParams[chain.Name].Inflation = inflation.Inflation
			}
		}(chain, rpc)
	}

	wg.Wait()

	response.Params = chainsParams
	return response
}
