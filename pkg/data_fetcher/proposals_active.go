package datafetcher

import (
	"main/pkg/types"
	"sync"
)

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

			rpc := f.GetRPC(chain)

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
