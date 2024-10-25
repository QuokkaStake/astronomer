package datafetcher

import (
	"main/pkg/types"
	"sync"
)

func (f *DataFetcher) GetSupply(chainNames []string) types.SupplyInfo {
	response := types.SupplyInfo{}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	chains, err := f.Database.GetChainsByNames(chainNames)
	if err != nil {
		response.Error = err
		return response
	}

	chainsSupplies := map[string]*types.ChainSupply{}
	amounts := []*types.AmountWithChain{}

	for _, chain := range chains {
		chainsSupplies[chain.Name] = &types.ChainSupply{Chain: chain}

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			rpc := f.GetRPC(chain)
			pool, _, err := rpc.GetPool()

			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsSupplies[chain.Name].PoolError = err
				return
			}

			chainsSupplies[chain.Name].BondedTokens = &types.Amount{
				Amount: pool.Pool.BondedTokens.ToLegacyDec(),
				Denom:  chain.BaseDenom,
			}
			chainsSupplies[chain.Name].NotBondedTokens = &types.Amount{
				Amount: pool.Pool.NotBondedTokens.ToLegacyDec(),
				Denom:  chain.BaseDenom,
			}

			amounts = append(amounts, &types.AmountWithChain{
				Chain:  chain.Name,
				Amount: chainsSupplies[chain.Name].BondedTokens,
			}, &types.AmountWithChain{
				Chain:  chain.Name,
				Amount: chainsSupplies[chain.Name].NotBondedTokens,
			})
		}(chain)

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			rpc := f.GetRPC(chain)
			supply, _, err := rpc.GetSupply()

			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsSupplies[chain.Name].SupplyError = err
				return
			}

			chainsSupplies[chain.Name].AllSupplies = make(map[string]*types.Amount, len(supply.Supply))

			for _, supply := range supply.Supply {
				supplyAmount := types.AmountFrom(supply)
				amounts = append(amounts, &types.AmountWithChain{Chain: chain.Name, Amount: supplyAmount})
				chainsSupplies[chain.Name].AllSupplies[supplyAmount.Denom] = supplyAmount
			}
		}(chain)

		wg.Add(1)
		go func(chain *types.Chain) {
			defer wg.Done()

			rpc := f.GetRPC(chain)
			communityPool, _, err := rpc.GetCommunityPool()

			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainsSupplies[chain.Name].CommunityPoolError = err
				return
			}

			chainsSupplies[chain.Name].AllCommunityPool = make(map[string]*types.Amount, len(communityPool.Pool))

			for _, communityPoolEntry := range communityPool.Pool {
				communityPoolAmount := types.AmountFromDec(communityPoolEntry)
				amounts = append(amounts, &types.AmountWithChain{Chain: chain.Name, Amount: communityPoolAmount})
				chainsSupplies[chain.Name].AllCommunityPool[communityPoolAmount.Denom] = communityPoolAmount
			}
		}(chain)
	}

	wg.Wait()

	f.PopulateDenoms(amounts)

	response.Supplies = chainsSupplies

	// TODO: refactor
	for _, chainSupplies := range response.Supplies {
		for supplyKey, supply := range chainSupplies.AllSupplies {
			if supply.IsIgnored() {
				delete(chainSupplies.AllSupplies, supplyKey)
			}
		}

		for communityPoolKey, communityPool := range chainSupplies.AllCommunityPool {
			if communityPool.IsIgnored() {
				delete(chainSupplies.AllCommunityPool, communityPoolKey)
			}
		}
	}

	return response
}
