package datafetcher

import (
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
	"sync"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

func (f *DataFetcher) predicateByQuery(query string) func(v stakingTypes.Validator) bool {
	lowercaseQuery := strings.ToLower(query)

	return func(v stakingTypes.Validator) bool {
		return strings.Contains(strings.ToLower(v.Description.Moniker), lowercaseQuery)
	}
}

func (f *DataFetcher) predicateByValidatorLinks(links []*types.ValidatorLink) func(v stakingTypes.Validator) bool {
	return func(v stakingTypes.Validator) bool {
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
	searchPredicate func(v stakingTypes.Validator) bool,
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

	validatorsResponses := map[string]*stakingTypes.QueryValidatorsResponse{}
	validatorsErrors := map[string]error{}
	signingInfosResponses := map[string]*slashingTypes.QuerySigningInfosResponse{}
	slashingParamsResponses := map[string]*slashingTypes.QueryParamsResponse{}

	for _, chain := range chains {
		wg.Add(3)

		go func(chain *types.Chain) {
			defer wg.Done()

			validators, validatorsErr := f.NodesManager.GetAllValidators(chain)
			mutex.Lock()
			validatorsResponses[chain.Name] = validators
			validatorsErrors[chain.Name] = validatorsErr
			mutex.Unlock()
		}(chain)

		go func(chain *types.Chain) {
			defer wg.Done()

			signingInfos, _ := f.NodesManager.GetAllSigningInfos(chain)
			mutex.Lock()
			if signingInfos != nil {
				signingInfosResponses[chain.Name] = signingInfos
			}
			mutex.Unlock()
		}(chain)

		go func(chain *types.Chain) {
			defer wg.Done()

			slashingParams, _ := f.NodesManager.GetSlashingParams(chain)
			mutex.Lock()
			if slashingParams != nil {
				slashingParamsResponses[chain.Name] = slashingParams
			}
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
		totalVP := utils.GetTotalVP(validatorsResponse.Validators)

		info := types.ChainValidatorsInfo{
			Chain:      chain,
			Explorers:  explorers.GetExplorersByChain(chain.Name),
			Error:      nil,
			Validators: make([]types.ValidatorInfo, len(foundValidators)),
		}

		if chainSlashingParams, ok := slashingParamsResponses[chain.Name]; ok {
			info.SlashingParams = &chainSlashingParams.Params
		}

		for index, validator := range foundValidators {
			validatorTokens := &types.Amount{
				Amount: validator.DelegatorShares,
				Denom:  chain.BaseDenom,
			}

			validatorInfo := types.ValidatorInfo{
				OperatorAddress:         validator.OperatorAddress,
				Jailed:                  validator.Jailed,
				Status:                  validator.Status.String(),
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

			if validator.Status == stakingTypes.Bonded {
				validatorInfo.Rank = utils.FindValidatorRank(validatorsResponse.Validators, validator.OperatorAddress)
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

		for index, validator := range foundValidators {
			signingInfo, found := utils.Find(signingInfos.Info, func(i slashingTypes.ValidatorSigningInfo) bool {
				consAddr := f.Converter.GetValidatorConsAddr(validator)
				equal, _ := f.Converter.CompareTwoBech32(consAddr, i.Address)
				return equal
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
