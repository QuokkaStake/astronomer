package datafetcher

import (
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
	"sync"
)

func (f *DataFetcher) predicateByQuery(query string) func(v *types.Validator) bool {
	lowercaseQuery := strings.ToLower(query)

	return func(v *types.Validator) bool {
		return strings.Contains(strings.ToLower(v.Description.Moniker), lowercaseQuery)
	}
}

func (f *DataFetcher) predicateByValidatorLinks(links []*types.ValidatorLink) func(v *types.Validator) bool {
	return func(v *types.Validator) bool {
		_, found := utils.Find(links, func(l *types.ValidatorLink) bool {
			return l.Address == v.OperatorAddress
		})

		return found
	}
}

func (f *DataFetcher) FindValidator(query string, chainNames []string) types.ValidatorsInfo {
	return f.FindValidatorGeneric(chainNames, f.predicateByQuery(query))
}

func (f *DataFetcher) FindMyValidators(
	chainNames []string,
	userID string,
	reporter string,
) types.ValidatorsInfo {
	validatorLinks, err := f.Database.FindValidatorLinksByUserAndReporter(userID, reporter)
	if err != nil {
		return types.ValidatorsInfo{Error: err}
	}

	return f.FindValidatorGeneric(chainNames, f.predicateByValidatorLinks(validatorLinks))
}

func (f *DataFetcher) FindValidatorGeneric(
	chainNames []string,
	searchPredicate func(v *types.Validator) bool,
) types.ValidatorsInfo {
	response := types.ValidatorsInfo{}

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

	var wg sync.WaitGroup
	var mutex sync.Mutex

	chainValidators := map[string]types.ChainValidatorsInfo{}
	denoms := []*types.AmountWithChain{}

	for _, chain := range chains {
		wg.Add(1)

		go func(chain *types.Chain) {
			defer wg.Done()

			rpc := f.GetRPC(chain)

			validators, _, err := rpc.GetAllValidators()
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				chainValidators[chain.Name] = types.ChainValidatorsInfo{
					Chain: chain,
					Error: err,
				}
				return
			}

			foundValidators := utils.Filter(validators.Validators, searchPredicate)

			totalVP := validators.GetTotalVP()

			info := types.ChainValidatorsInfo{
				Chain:      chain,
				Explorers:  explorers.GetExplorersByChain(chain.Name),
				Error:      nil,
				Validators: make([]types.ValidatorInfo, len(foundValidators)),
			}

			for index, validator := range foundValidators {
				validatorTokens := &types.Amount{
					Amount: validator.DelegatorShares,
					Denom:  chain.BaseDenom,
				}

				validatorInfo := types.ValidatorInfo{
					OperatorAddress:         validator.OperatorAddress,
					Jailed:                  validator.Jailed,
					Status:                  validator.Status,
					Tokens:                  validatorTokens,
					Moniker:                 validator.Description.Moniker,
					Details:                 validator.Description.Details,
					Identity:                validator.Description.Identity,
					Website:                 validator.Description.Website,
					SecurityContact:         validator.Description.SecurityContact,
					Commission:              validator.Commission.CommissionRates.Rate.MustFloat64(),
					CommissionMax:           validator.Commission.CommissionRates.MaxRate.MustFloat64(),
					CommissionMaxChangeRate: validator.Commission.CommissionRates.MaxChangeRate.MustFloat64(),
					VotingPowerPercent:      validator.DelegatorShares.Quo(totalVP).MustFloat64(),
				}

				if validator.Active() {
					validatorInfo.Rank = validators.FindValidatorRank(validator.OperatorAddress)
				}

				info.Validators[index] = validatorInfo
				denoms = append(denoms, &types.AmountWithChain{
					Chain:  chain.Name,
					Amount: validatorTokens,
				})
			}

			chainValidators[chain.Name] = info
		}(chain)
	}

	wg.Wait()

	f.PopulateDenoms(denoms)

	response.Chains = chainValidators
	return response
}
