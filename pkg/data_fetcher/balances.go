package datafetcher

import (
	"fmt"
	"main/pkg/types"
	"main/pkg/utils"
	"sync"
)

func (f *DataFetcher) GetBalances(userID, reporter string) *types.WalletsBalancesInfo { //nolint:maintidx
	response := &types.WalletsBalancesInfo{
		Infos: map[string]*types.ChainWalletsBalancesInfo{},
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

	var wg sync.WaitGroup
	var mutex sync.Mutex

	amountsWithChains := []*types.AmountWithChain{}
	validators := []*types.ValidatorAddressWithMoniker{}

	for chainName, chainWallets := range walletsByChain {
		chain, ok := chainsMap[chainName]
		if !ok {
			panic(fmt.Errorf("chain %s not found", chainName))
		}

		response.SetChain(chain, explorers.GetExplorersByChain(chain.Name))

		for _, chainWallet := range chainWallets {
			response.SetAddressInfo(chain.Name, chainWallet)

			// balances
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := f.GetRPC(chain)

				balances, _, err := rpc.GetBalance(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				if err != nil {
					response.SetBalancesError(chain.Name, chainWallet, err)
					return
				}

				walletBalances := utils.Map(balances.Balances, types.AmountFrom)

				response.SetBalances(chain.Name, chainWallet, walletBalances)

				for _, amount := range walletBalances {
					amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
						Chain:  chain.Name,
						Amount: amount,
					})
				}
			}(chain, chainWallet)

			// rewards
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := f.GetRPC(chain)

				rewards, _, err := rpc.GetRewards(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				if err != nil {
					response.SetRewardsError(chain.Name, chainWallet, err)
					return
				}

				walletRewards := utils.Map(rewards.Total, types.AmountFromDec)

				response.SetRewards(chain.Name, chainWallet, walletRewards)

				for _, amount := range walletRewards {
					amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
						Chain:  chain.Name,
						Amount: amount,
					})
				}
			}(chain, chainWallet)

			// commission
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				valoper, err := utils.ConvertBech32Prefix(chainWallet.Address, chain.Bech32ValidatorPrefix)
				if err != nil {
					mutex.Lock()
					response.SetCommissionsError(chain.Name, chainWallet, err)
					mutex.Unlock()
					return
				}

				rpc := f.GetRPC(chain)

				rewards, _, err := rpc.GetCommission(valoper)
				mutex.Lock()
				defer mutex.Unlock()

				if err != nil {
					response.SetCommissionsError(chain.Name, chainWallet, err)
					return
				}

				walletCommissions := utils.Map(rewards.Commission.Commission, func(b types.SdkAmount) *types.Amount {
					return b.ToAmount()
				})

				response.SetCommissions(chain.Name, chainWallet, walletCommissions)

				for _, amount := range walletCommissions {
					amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
						Chain:  chain.Name,
						Amount: amount,
					})
				}
			}(chain, chainWallet)

			// delegations
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := f.GetRPC(chain)

				delegations, _, err := rpc.GetDelegations(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				if err != nil {
					response.SetDelegationsError(chain.Name, chainWallet, err)
					return
				}
				walletDelegations := utils.Map(delegations.Delegations, func(b types.SdkDelegation) *types.Delegation {
					return &types.Delegation{
						Amount: b.Balance.ToAmount(),
						Validator: &types.ValidatorAddressWithMoniker{
							Chain:   chain,
							Address: b.Delegation.ValidatorAddress,
						},
					}
				})

				response.SetDelegations(chain.Name, chainWallet, walletDelegations)

				for _, delegation := range walletDelegations {
					amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
						Chain:  chain.Name,
						Amount: delegation.Amount,
					})

					validators = append(validators, delegation.Validator)
				}
			}(chain, chainWallet)

			// redelegations
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := f.GetRPC(chain)

				redelegations, _, err := rpc.GetRedelegations(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				if err != nil {
					response.SetRedelegationsError(chain.Name, chainWallet, err)
					return
				}
				walletRedelegations := []*types.Redelegation{}

				for _, redelegation := range redelegations.Redelegations {
					for _, entry := range redelegation.Entries {
						amount := &types.Amount{
							Amount: entry.Balance,
							Denom:  chain.BaseDenom,
						}

						srcValidator := &types.ValidatorAddressWithMoniker{
							Chain:   chain,
							Address: redelegation.Redelegation.ValidatorSrcAddress,
						}

						dstValidator := &types.ValidatorAddressWithMoniker{
							Chain:   chain,
							Address: redelegation.Redelegation.ValidatorDstAddress,
						}

						amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
							Chain:  chain.Name,
							Amount: amount,
						})

						validators = append(validators, srcValidator, dstValidator)

						walletRedelegations = append(walletRedelegations, &types.Redelegation{
							Amount:         amount,
							SrcValidator:   srcValidator,
							DstValidator:   dstValidator,
							CompletionTime: entry.Entry.CompletionTime,
						})
					}
				}

				response.SetRedelegations(chain.Name, chainWallet, walletRedelegations)
			}(chain, chainWallet)

			// unbonds
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := f.GetRPC(chain)

				unbonds, _, err := rpc.GetUnbonds(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				if err != nil {
					response.SetUnbondsError(chain.Name, chainWallet, err)
					return
				}

				walletUnbonds := []*types.Unbond{}

				for _, unbond := range unbonds.Unbonds {
					for _, entry := range unbond.Entries {
						amount := &types.Amount{
							Amount: entry.Balance,
							Denom:  chain.BaseDenom,
						}

						validator := &types.ValidatorAddressWithMoniker{
							Chain:   chain,
							Address: unbond.ValidatorAddress,
						}

						amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
							Chain:  chain.Name,
							Amount: amount,
						})

						validators = append(validators, validator)

						walletUnbonds = append(walletUnbonds, &types.Unbond{
							Amount:         amount,
							Validator:      validator,
							CompletionTime: entry.CompletionTime,
						})
					}
				}

				response.SetUnbonds(chain.Name, chainWallet, walletUnbonds)
			}(chain, chainWallet)
		}
	}

	wg.Wait()

	f.PopulateDenoms(amountsWithChains)
	f.PopulateValidators(validators)

	// TODO: refactor
	for _, chainBalances := range response.Infos {
		for _, walletBalances := range chainBalances.BalancesInfo {
			walletBalances.Balances = utils.Filter(walletBalances.Balances, func(a *types.Amount) bool {
				return a.PriceUSD != nil
			})

			walletBalances.Rewards = utils.Filter(walletBalances.Rewards, func(a *types.Amount) bool {
				return a.PriceUSD != nil
			})

			walletBalances.Commissions = utils.Filter(walletBalances.Commissions, func(a *types.Amount) bool {
				return a.PriceUSD != nil
			})

			walletBalances.Delegations = utils.Filter(walletBalances.Delegations, func(d *types.Delegation) bool {
				return d.Amount.PriceUSD != nil
			})
		}
	}

	return response
}
