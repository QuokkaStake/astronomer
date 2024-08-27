package datafetcher

import (
	"fmt"
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

			fmt.Printf("set pool: %+v\n", pool)

			chainsSupplies[chain.Name].BondedTokens = &types.Amount{
				Amount: pool.Pool.BondedTokens,
				Denom:  chain.BaseDenom,
			}
			chainsSupplies[chain.Name].NotBondedTokens = &types.Amount{
				Amount: pool.Pool.NotBondedTokens,
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
	}

	wg.Wait()

	f.PopulateDenoms(amounts)

	response.Supplies = chainsSupplies
	return response
}
