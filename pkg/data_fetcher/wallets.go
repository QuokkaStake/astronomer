package datafetcher

import (
	"fmt"
	"main/pkg/tendermint"
	"main/pkg/types"
	"main/pkg/utils"
	"sync"
)

func (f *DataFetcher) GetBalances(userID, reporter string) types.WalletsBalancesInfo { //nolint:gocyclo,maintidx // need to refactor this
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
	validators := []*types.ValidatorAddressWithMoniker{}

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
			// balances
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := tendermint.NewRPC(chain, 10, f.Logger)

				balances, _, err := rpc.GetBalance(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				balanceInfo, ok := chainInfos[chain.Name].BalancesInfo[chainWallet.Address]
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

			// rewards
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := tendermint.NewRPC(chain, 10, f.Logger)

				rewards, _, err := rpc.GetRewards(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				balanceInfo, ok := chainInfos[chain.Name].BalancesInfo[chainWallet.Address]
				if !ok {
					balanceInfo = &types.WalletBalancesInfo{
						Address: chainWallet,
					}
				}

				if err != nil {
					balanceInfo.RewardsError = err
				} else {
					balanceInfo.Rewards = utils.Map(rewards.Total, func(b types.SdkAmount) *types.Amount {
						return b.ToAmount()
					})

					for _, amount := range balanceInfo.Rewards {
						amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
							Chain:  chain.Name,
							Amount: amount,
						})
					}
				}

				chainInfos[chain.Name].BalancesInfo[chainWallet.Address] = balanceInfo
			}(chain, chainWallet)

			// commission
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				valoper, err := utils.ConvertBech32Prefix(chainWallet.Address, chain.Bech32ValidatorPrefix)
				if err != nil {
					mutex.Lock()
					balanceInfo, ok := chainInfos[chain.Name].BalancesInfo[chainWallet.Address]
					if !ok {
						balanceInfo = &types.WalletBalancesInfo{
							Address: chainWallet,
						}
					}

					balanceInfo.CommissionsError = err

					chainInfos[chain.Name].BalancesInfo[chainWallet.Address] = balanceInfo
					mutex.Unlock()
					return
				}

				rpc := tendermint.NewRPC(chain, 10, f.Logger)

				rewards, _, err := rpc.GetCommission(valoper)
				mutex.Lock()
				defer mutex.Unlock()

				balanceInfo, ok := chainInfos[chain.Name].BalancesInfo[chainWallet.Address]
				if !ok {
					balanceInfo = &types.WalletBalancesInfo{
						Address: chainWallet,
					}
				}

				if err != nil {
					balanceInfo.CommissionsError = err
				} else {
					balanceInfo.Commissions = utils.Map(rewards.Commission.Commission, func(b types.SdkAmount) *types.Amount {
						return b.ToAmount()
					})

					for _, amount := range balanceInfo.Commissions {
						amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
							Chain:  chain.Name,
							Amount: amount,
						})
					}
				}

				chainInfos[chain.Name].BalancesInfo[chainWallet.Address] = balanceInfo
			}(chain, chainWallet)

			// delegations
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := tendermint.NewRPC(chain, 10, f.Logger)

				delegations, _, err := rpc.GetDelegations(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				balanceInfo, ok := chainInfos[chain.Name].BalancesInfo[chainWallet.Address]
				if !ok {
					balanceInfo = &types.WalletBalancesInfo{
						Address: chainWallet,
					}
				}

				if err != nil {
					balanceInfo.DelegationsError = err
				} else {
					balanceInfo.Delegations = utils.Map(delegations.Delegations, func(b types.SdkDelegation) *types.Delegation {
						return &types.Delegation{
							Amount: b.Balance.ToAmount(),
							Validator: &types.ValidatorAddressWithMoniker{
								Chain:   chain,
								Address: b.Delegation.ValidatorAddress,
							},
						}
					})

					for _, delegation := range balanceInfo.Delegations {
						amountsWithChains = append(amountsWithChains, &types.AmountWithChain{
							Chain:  chain.Name,
							Amount: delegation.Amount,
						})

						validators = append(validators, delegation.Validator)
					}
				}

				chainInfos[chain.Name].BalancesInfo[chainWallet.Address] = balanceInfo
			}(chain, chainWallet)

			// redelegations
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := tendermint.NewRPC(chain, 10, f.Logger)

				redelegations, _, err := rpc.GetRedelegations(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				balanceInfo, ok := chainInfos[chain.Name].BalancesInfo[chainWallet.Address]
				if !ok {
					balanceInfo = &types.WalletBalancesInfo{
						Address: chainWallet,
					}
				}

				if err != nil {
					balanceInfo.RedelegationsError = err
				} else {
					balanceInfo.Redelegations = []*types.Redelegation{}

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

							balanceInfo.Redelegations = append(balanceInfo.Redelegations, &types.Redelegation{
								Amount:         amount,
								SrcValidator:   srcValidator,
								DstValidator:   dstValidator,
								CompletionTime: entry.Entry.CompletionTime,
							})
						}
					}
				}

				chainInfos[chain.Name].BalancesInfo[chainWallet.Address] = balanceInfo
			}(chain, chainWallet)

			// unbonds
			wg.Add(1)
			go func(chain *types.Chain, chainWallet *types.WalletLink) {
				defer wg.Done()

				rpc := tendermint.NewRPC(chain, 10, f.Logger)

				unbonds, _, err := rpc.GetUnbonds(chainWallet.Address)
				mutex.Lock()
				defer mutex.Unlock()

				balanceInfo, ok := chainInfos[chain.Name].BalancesInfo[chainWallet.Address]
				if !ok {
					balanceInfo = &types.WalletBalancesInfo{
						Address: chainWallet,
					}
				}

				if err != nil {
					balanceInfo.UnbondsError = err
				} else {
					balanceInfo.Unbonds = []*types.Unbond{}

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

							balanceInfo.Unbonds = append(balanceInfo.Unbonds, &types.Unbond{
								Amount:         amount,
								Validator:      validator,
								CompletionTime: entry.CompletionTime,
							})
						}
					}
				}

				chainInfos[chain.Name].BalancesInfo[chainWallet.Address] = balanceInfo
			}(chain, chainWallet)
		}
	}

	wg.Wait()

	f.PopulateDenoms(amountsWithChains)
	f.PopulateValidators(validators)
	response.Infos = chainInfos

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
