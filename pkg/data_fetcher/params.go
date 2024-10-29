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

		rpc := f.GetRPC(chain)

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
		go func(chain *types.Chain) {
			defer wg.Done()

			params, slashingParamsErr := f.NodesManager.GetSlashingParams(chain)
			mutex.Lock()
			defer mutex.Unlock()

			if slashingParamsErr != nil {
				chainsParams[chain.Name].SlashingParamsError = slashingParamsErr
			} else {
				chainsParams[chain.Name].SlashingParams = params.Params
			}
		}(chain)

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			params, paramsErr := f.NodesManager.GetGovParams(chain, "voting")
			mutex.Lock()
			defer mutex.Unlock()

			if paramsErr != nil {
				chainsParams[chain.Name].VotingParamsError = paramsErr
			} else {
				chainsParams[chain.Name].VotingParams = params.VotingParams
			}
		}(chain)

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			params, paramsErr := f.NodesManager.GetGovParams(chain, "deposit")
			mutex.Lock()
			defer mutex.Unlock()

			if paramsErr != nil {
				chainsParams[chain.Name].DepositParamsError = paramsErr
			} else {
				chainsParams[chain.Name].DepositParams = params.DepositParams
			}
		}(chain)

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			params, paramsErr := f.NodesManager.GetGovParams(chain, "tallying")
			mutex.Lock()
			defer mutex.Unlock()

			if paramsErr != nil {
				chainsParams[chain.Name].TallyParamsError = paramsErr
			} else {
				chainsParams[chain.Name].TallyParams = params.TallyParams
			}
		}(chain)

		wg.Add(1)
		go func(chain *types.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			blockTime, blockTimeErr := f.NodesManager.GetBlockTime(chain)
			mutex.Lock()
			defer mutex.Unlock()

			if blockTimeErr != nil {
				chainsParams[chain.Name].BlockTimeError = blockTimeErr
			} else {
				chainsParams[chain.Name].BlockTime = blockTime
			}
		}(chain, rpc)

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			params, paramsErr := f.NodesManager.GetMintParams(chain)
			mutex.Lock()
			defer mutex.Unlock()

			if paramsErr != nil {
				chainsParams[chain.Name].MintParamsError = paramsErr
			} else {
				chainsParams[chain.Name].MintParams = params.Params
			}
		}(chain)

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			inflation, inflationErr := f.NodesManager.GetInflation(chain)
			mutex.Lock()
			defer mutex.Unlock()

			if inflationErr != nil {
				chainsParams[chain.Name].InflationError = inflationErr
			} else {
				chainsParams[chain.Name].Inflation = inflation.Inflation
			}
		}(chain)
	}

	wg.Wait()

	response.Params = chainsParams
	return response
}
