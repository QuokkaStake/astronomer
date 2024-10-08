package datafetcher

import (
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
	"sync"

	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
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

	validatorsResponses := map[string]*types.ValidatorsResponse{}
	validatorsErrors := map[string]error{}
	signingInfosResponses := map[string]*slashingTypes.QuerySigningInfosResponse{}

	for _, chain := range chains {
		wg.Add(2)

		go func(chain *types.Chain) {
			defer wg.Done()

			rpc := f.GetRPC(chain)

			validators, _, err := rpc.GetAllValidators()
			mutex.Lock()
			validatorsResponses[chain.Name] = validators
			validatorsErrors[chain.Name] = err
			mutex.Unlock()
		}(chain)

		go func(chain *types.Chain) {
			defer wg.Done()

			rpc := f.GetRPC(chain)

			signingInfos, _, _ := rpc.GetAllSigningInfos()
			mutex.Lock()
			signingInfosResponses[chain.Name] = signingInfos
			mutex.Unlock()
		}(chain)
	}

	wg.Wait()

	validatorsInfos := map[string]types.ChainValidatorsInfo{}
	denoms := []*types.AmountWithChain{}

	for _, chain := range chains {
		if chainErr, ok := validatorsErrors[chain.Name]; ok && chainErr != nil {
			validatorsInfos[chain.Name] = types.ChainValidatorsInfo{
				Chain: chain,
				Error: chainErr,
			}
			continue
		}

		validatorsResponse := validatorsResponses[chain.Name]
		foundValidators := utils.Filter(validatorsResponse.Validators, searchPredicate)
		totalVP := validatorsResponse.GetTotalVP()

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
				validatorInfo.Rank = validatorsResponse.FindValidatorRank(validator.OperatorAddress)
			}

			info.Validators[index] = validatorInfo
			denoms = append(denoms, &types.AmountWithChain{
				Chain:  chain.Name,
				Amount: validatorTokens,
			})
		}

		signingInfos, ok := signingInfosResponses[chain.Name]
		if !ok {
			validatorsInfos[chain.Name] = info
			continue
		}

		for index := range foundValidators {
			signingInfo, found := utils.Find(signingInfos.Info, func(i slashingTypes.ValidatorSigningInfo) bool {
				return false
			})

			if !found {
				continue
			}

			info.Validators[index].SigningInfo = &signingInfo
		}

		validatorsInfos[chain.Name] = info
	}

	f.PopulateDenoms(denoms)

	response.Chains = validatorsInfos
	return response
}
