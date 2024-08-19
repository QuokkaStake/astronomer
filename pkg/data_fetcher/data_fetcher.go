package datafetcher

import (
	"fmt"
	"main/pkg/database"
	priceFetcher "main/pkg/price_fetcher"
	"main/pkg/tendermint"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Logger        zerolog.Logger
	Database      *database.Database
	PriceFetchers []priceFetcher.PriceFetcher
}

func NewDataFetcher(logger zerolog.Logger, database *database.Database) *DataFetcher {
	priceFetchers := []priceFetcher.PriceFetcher{
		priceFetcher.NewCoingeckoPriceFetcher(logger),
	}

	return &DataFetcher{
		Logger:        logger.With().Str("component", "data_fetcher").Logger(),
		Database:      database,
		PriceFetchers: priceFetchers,
	}
}

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

func (f *DataFetcher) GetActiveProposals(chainNames []string) types.ActiveProposals {
	response := types.ActiveProposals{}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	chains, err := f.Database.GetChainsByNames(chainNames)
	if err != nil {
		response.Error = err
		return response
	}

	explorers, err := f.Database.GetExplorersByChains(chainNames)
	if err != nil {
		response.Error = err
		return response
	}

	chainProposals := map[string]*types.ChainActiveProposals{}

	for _, chain := range chains {
		chainProposals[chain.Name] = &types.ChainActiveProposals{
			Chain:     chain,
			Explorers: explorers.GetExplorersByChain(chain.Name),
		}

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			rpc := tendermint.NewRPC(chain, 10, f.Logger)

			proposals, _, err := rpc.GetActiveProposals()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainProposals[chain.Name].ProposalsError = err
			} else {
				chainProposals[chain.Name].Proposals = proposals
			}
		}(chain)
	}

	wg.Wait()

	response.Proposals = chainProposals

	return response
}

func (f *DataFetcher) GetSingleProposal(chainName string, proposalID string) types.SingleProposal {
	response := types.SingleProposal{}

	chains, err := f.Database.GetChainsByNames([]string{chainName})
	if err != nil {
		response.Error = err
		return response
	}

	explorers, err := f.Database.GetExplorersByChains([]string{chainName})
	if err != nil {
		response.Error = err
		return response
	}

	if len(chains) != 1 {
		response.Error = fmt.Errorf("chain '%s' is not found", chainName)
		return response
	}

	response.Chain = chains[0]
	response.Explorers = explorers.GetExplorersByChain(chainName)

	rpc := tendermint.NewRPC(chains[0], 10, f.Logger)
	proposal, _, err := rpc.GetSingleProposal(proposalID)

	if err != nil {
		response.Error = err
	} else {
		response.Proposal = proposal
	}

	return response
}
