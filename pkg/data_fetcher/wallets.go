package datafetcher

import (
	"main/pkg/types"
	"main/pkg/utils"
)

func (f *DataFetcher) GetWallets(userID, reporter string) *types.WalletsList {
	response := &types.WalletsList{
		Infos: map[string]*types.ChainWalletsList{},
	}

	wallets, err := f.Database.FindWalletLinksByUserAndReporter(userID, reporter)
	if err != nil {
		response.Error = err
		return response
	}

	chainNames := utils.MapUniq(wallets, func(w *types.WalletLink) string {
		return w.Chain
	})

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

	walletsByChain := utils.GroupBy(wallets, func(w *types.WalletLink) []string {
		return []string{w.Chain}
	})

	chainsMap := utils.GroupSingleBy(chains, func(c *types.Chain) string {
		return c.Name
	})

	for chainName, chainWallets := range walletsByChain {
		chain := chainsMap[chainName]
		response.Infos[chainName] = &types.ChainWalletsList{
			Chain:     chain,
			Explorers: explorers.GetExplorersByChain(chainName),
			Wallets:   chainWallets,
		}
	}

	return response
}
