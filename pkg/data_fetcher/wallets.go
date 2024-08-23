package datafetcher

import (
	"fmt"
	"main/pkg/tendermint"
	"main/pkg/types"
	"main/pkg/utils"
	"sync"
)

func (f *DataFetcher) GetBalances(userID, reporter string) types.WalletsBalancesInfo {
	response := types.WalletsBalancesInfo{}

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

	var wg sync.WaitGroup
	var mutex sync.Mutex

	chainInfos := map[string]types.ChainWalletsBalancesInfo{}
	amountsWithChains := []*types.AmountWithChain{}

	for chainName, chainWallets := range walletsByChain {
		chain, ok := chainsMap[chainName]
		if !ok {
			panic(fmt.Errorf("chain %s not found", chainName))
		}

		chainInfos[chainName] = types.ChainWalletsBalancesInfo{
			Chain:        chain,
			Explorers:    explorers.GetExplorersByChain(chain.Name),
			BalancesInfo: map[string]*types.WalletBalancesInfo{},
		}

		for _, chainWallet := range chainWallets {
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := tendermint.NewRPC(chain, 10, f.Logger)

				balances, _, err := rpc.GetBalance(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				balanceInfo, ok := chainInfos[chainName].BalancesInfo[chainWallet.Address]
				if !ok {
					balanceInfo = &types.WalletBalancesInfo{
						Address: chainWallet,
					}
				}

				if err != nil {
					balanceInfo.BalancesError = err
				} else {
					balanceInfo.Balances = utils.Map(balances.Balances, func(b types.SdkAmount) *types.Amount {
						return b.ToAmount()
					})

					for _, amount := range balanceInfo.Balances {
						amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
							Chain:  chain.Name,
							Amount: amount,
						})
					}
				}

				chainInfos[chain.Name].BalancesInfo[chainWallet.Address] = balanceInfo
			}(chain, chainWallet)
		}
	}

	wg.Wait()

	f.PopulateDenoms(amountsWithChains)
	response.Infos = chainInfos

	for _, chainBalances := range response.Infos {
		for _, walletBalances := range chainBalances.BalancesInfo {
			walletBalances.Balances = utils.Filter(walletBalances.Balances, func(a *types.Amount) bool {
				return a.PriceUSD != nil
			})
		}
	}

	return response
}
