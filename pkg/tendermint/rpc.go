package tendermint

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/pkg/constants"
	converterPkg "main/pkg/converter"
	"main/pkg/http"
	"main/pkg/metrics"
	"main/pkg/types"
	"main/pkg/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govV1beta1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/gogoproto/proto"

	"github.com/rs/zerolog"
)

type RPC struct {
	Chain          *types.Chain
	Host           string
	Client         *http.Client
	Timeout        int
	Logger         zerolog.Logger
	MetricsManager *metrics.Manager
	Converter      *converterPkg.Converter
}

func NewRPC(
	chain *types.Chain,
	timeout int,
	logger *zerolog.Logger,
	converter *converterPkg.Converter,
	metricsManager *metrics.Manager,
) *RPC {
	return &RPC{
		Chain:   chain,
		Client:  http.NewClient(logger, chain.Name),
		Timeout: timeout,
		Logger: logger.With().
			Str("component", "rpc").
			Str("chain", chain.Name).
			Logger(),
		Converter:      converter,
		MetricsManager: metricsManager,
	}
}

func (rpc *RPC) GetAllValidators(hosts []string) (*stakingTypes.QueryValidatorsResponse, types.QueryInfo, error) {
	url := "/cosmos/staking/v1beta1/validators?pagination.count_total=true&pagination.limit=1000"

	var response stakingTypes.QueryValidatorsResponse
	info, err := rpc.Get(hosts, url, "validators", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetAllSigningInfos(hosts []string) (*slashingTypes.QuerySigningInfosResponse, types.QueryInfo, error) {
	url := "/cosmos/slashing/v1beta1/signing_infos?pagination.limit=1000"

	var response slashingTypes.QuerySigningInfosResponse
	info, err := rpc.Get(hosts, url, "signing_infos", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetValidator(address string) (*stakingTypes.QueryValidatorResponse, types.QueryInfo, error) {
	url := "/cosmos/staking/v1beta1/validators/" + address

	var response stakingTypes.QueryValidatorResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "validator", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetStakingParams() (*stakingTypes.QueryParamsResponse, types.QueryInfo, error) {
	url := "/cosmos/staking/v1beta1/params"

	var response stakingTypes.QueryParamsResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "staking_params", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetSlashingParams(hosts []string) (*slashingTypes.QueryParamsResponse, types.QueryInfo, error) {
	url := "/cosmos/slashing/v1beta1/params"

	var response slashingTypes.QueryParamsResponse
	info, err := rpc.Get(hosts, url, "slashing_params", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetGovParams(paramsType string, hosts []string) (*govV1beta1Types.QueryParamsResponse, types.QueryInfo, error) {
	url := "/cosmos/gov/v1beta1/params/" + paramsType

	var response govV1beta1Types.QueryParamsResponse
	info, err := rpc.Get(hosts, url, "gov_params_"+paramsType, &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetMintParams(hosts []string) (*mintTypes.QueryParamsResponse, types.QueryInfo, error) {
	url := "/cosmos/mint/v1beta1/params"

	var response mintTypes.QueryParamsResponse
	info, err := rpc.Get(hosts, url, "mint_params", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetInflation(hosts []string) (*mintTypes.QueryInflationResponse, types.QueryInfo, error) {
	url := "/cosmos/mint/v1beta1/inflation"

	var response mintTypes.QueryInflationResponse
	info, err := rpc.Get(hosts, url, "inflation", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetBalance(address string, hosts []string) (*bankTypes.QueryAllBalancesResponse, types.QueryInfo, error) {
	url := "/cosmos/bank/v1beta1/balances/" + address

	var response bankTypes.QueryAllBalancesResponse
	info, err := rpc.Get(hosts, url, "balance", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetRewards(address string) (*distributionTypes.QueryDelegationTotalRewardsResponse, types.QueryInfo, error) {
	url := "/cosmos/distribution/v1beta1/delegators/" + address + "/rewards"

	var response distributionTypes.QueryDelegationTotalRewardsResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "rewards", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetCommission(address string) (*distributionTypes.QueryValidatorCommissionResponse, types.QueryInfo, error) {
	url := "/cosmos/distribution/v1beta1/validators/" + address + "/commission"

	var response distributionTypes.QueryValidatorCommissionResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "commission", &response)
	if err != nil {
		// not being a validator is acceptable
		if strings.Contains(err.Error(), "validator does not exist") {
			return &distributionTypes.QueryValidatorCommissionResponse{
				Commission: distributionTypes.ValidatorAccumulatedCommission{
					Commission: make(cosmosTypes.DecCoins, 0),
				},
			}, info, nil
		}

		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetDelegations(address string) (*stakingTypes.QueryDelegatorDelegationsResponse, types.QueryInfo, error) {
	url := "/cosmos/staking/v1beta1/delegations/" + address + "?pagination.limit=1000"

	var response stakingTypes.QueryDelegatorDelegationsResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "delegations", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetRedelegations(address string) (*stakingTypes.QueryRedelegationsResponse, types.QueryInfo, error) {
	url := "/cosmos/staking/v1beta1/delegators/" + address + "/redelegations?pagination.limit=1000"

	var response stakingTypes.QueryRedelegationsResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "commission", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetUnbonds(address string) (*stakingTypes.QueryDelegatorUnbondingDelegationsResponse, types.QueryInfo, error) {
	url := "/cosmos/staking/v1beta1/delegators/" + address + "/unbonding_delegations?pagination.limit=1000"

	var response stakingTypes.QueryDelegatorUnbondingDelegationsResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "unbonds", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetPool() (*stakingTypes.QueryPoolResponse, types.QueryInfo, error) {
	url := "/cosmos/staking/v1beta1/pool"

	var response stakingTypes.QueryPoolResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "pool", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetSupply() (*bankTypes.QueryTotalSupplyResponse, types.QueryInfo, error) {
	url := "/cosmos/bank/v1beta1/supply?pagination.limit=10000&pagination.offset=0"

	var response bankTypes.QueryTotalSupplyResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "supply", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetCommunityPool() (*distributionTypes.QueryCommunityPoolResponse, types.QueryInfo, error) {
	url := "/cosmos/distribution/v1beta1/community_pool?pagination.limit=10000&pagination.offset=0"

	var response distributionTypes.QueryCommunityPoolResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "community_pool", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetBlockTime() (time.Duration, error) {
	var newerBlock cmtservice.GetLatestBlockResponse
	_, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, "/cosmos/base/tendermint/v1beta1/blocks/latest", "block", &newerBlock)
	if err != nil {
		return 0, err
	}

	newerHeight := newerBlock.Block.Header.Height - 1000 //nolint:staticcheck

	var olderBlock cmtservice.GetBlockByHeightResponse
	_, err = rpc.Get(
		[]string{rpc.Chain.LCDEndpoint},
		"/cosmos/base/tendermint/v1beta1/blocks/"+strconv.FormatInt(newerHeight, 10),
		"block",
		&olderBlock,
	)
	if err != nil {
		return 0, err
	}

	timeDiff := olderBlock.Block.Header.Time.Sub(newerBlock.Block.Header.Time)    //nolint:staticcheck
	heightDiff := olderBlock.Block.Header.Height - newerBlock.Block.Header.Height //nolint:staticcheck

	return time.Duration(float64(timeDiff.Nanoseconds()) / float64(heightDiff)), nil
}

func (rpc *RPC) GetActiveProposals() ([]*types.Proposal, types.QueryInfo, error) {
	url := "/cosmos/gov/v1/proposals?pagination.limit=1000&proposal_status=PROPOSAL_STATUS_VOTING_PERIOD"

	var response govV1Types.QueryProposalsResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "proposals_v1", &response)
	if err == nil {
		return utils.Map(response.Proposals, types.ProposalFromV1), info, nil
	}

	if !strings.Contains(err.Error(), "Not Implemented") {
		return nil, info, err
	}

	rpc.Logger.Warn().Msg("v1 proposals are not supported, falling back to v1beta1")

	url = "/cosmos/gov/v1beta1/proposals?pagination.limit=1000&proposal_status=2"

	var responsev1beta1 govV1beta1Types.QueryProposalsResponse
	infov1beta1, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "proposals_v1beta1", &responsev1beta1)
	if err != nil {
		return nil, info, err
	}

	for _, proposal := range responsev1beta1.Proposals {
		if err := rpc.Converter.UnpackProposal(proposal); err != nil {
			return nil, info, err
		}
	}

	return utils.Map(responsev1beta1.Proposals, types.ProposalFromV1beta1), infov1beta1, nil
}

func (rpc *RPC) GetSingleProposal(proposalID string) (*types.Proposal, types.QueryInfo, error) {
	url := "/cosmos/gov/v1/proposals/" + proposalID

	var response govV1Types.QueryProposalResponse
	info, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "proposal_v1", &response)
	if err == nil {
		return types.ProposalFromV1(response.Proposal), info, nil
	}

	// failed cases
	if strings.Contains(err.Error(), "doesn't exist") {
		return nil, info, nil
	}

	if !strings.Contains(err.Error(), "Not Implemented") {
		return nil, info, err
	}

	rpc.Logger.Warn().Msg("v1 proposal are not supported, falling back to v1")

	url = "/cosmos/gov/v1beta1/proposals/" + proposalID

	var responsev1beta1 govV1beta1Types.QueryProposalResponse
	infov1beta1, err := rpc.Get([]string{rpc.Chain.LCDEndpoint}, url, "proposal_v1beta1", &responsev1beta1)
	if err != nil {
		if strings.Contains(err.Error(), "doesn't exist") {
			return nil, info, nil
		}

		return nil, info, err
	}

	if err := rpc.Converter.UnpackProposal(responsev1beta1.Proposal); err != nil {
		return nil, info, err
	}

	return types.ProposalFromV1beta1(responsev1beta1.Proposal), infov1beta1, nil
}

func (rpc *RPC) Get(
	hosts []string,
	url string,
	queryName string,
	target proto.Message,
) (types.QueryInfo, error) {
	for attempt := range constants.RetriesCount {
		host := hosts[rand.Int()%len(hosts)]
		queryInfo, err := rpc.GetOne(host, url, queryName, target)
		if err != nil {
			rpc.Logger.Warn().
				Str("host", host).
				Str("url", url).
				Int("attempt", attempt).
				Int("max_attempts", constants.RetriesCount).
				Err(err).
				Msg("LCD request failed, retrying")
		} else {
			return queryInfo, nil
		}
	}

	return types.QueryInfo{}, fmt.Errorf("could not get data after %d attempts", constants.RetriesCount)
}

func (rpc *RPC) GetOne(
	host string,
	url string,
	queryName string,
	target proto.Message,
) (types.QueryInfo, error) {
	bytes, queryInfo, err := rpc.Client.GetPlain(
		host+url,
		queryName,
	)
	if err != nil {
		rpc.Logger.Warn().
			Str("host", host).
			Str("url", url).
			Err(err).Msg("LCD request failed")
		return queryInfo, err
	}

	// check whether the response is error first
	var errorResponse types.LCDError
	if unmarshalErr := json.Unmarshal(bytes, &errorResponse); unmarshalErr == nil {
		// if we successfully unmarshalled it into LCDError, so err == nil,
		// that means the response is indeed an error.
		if errorResponse.Code != 0 {
			rpc.Logger.Warn().
				Str("host", host).
				Str("url", url).
				Err(unmarshalErr).
				Int("code", errorResponse.Code).
				Str("message", errorResponse.Message).
				Msg("LCD request returned an error")
			return queryInfo, errors.New(errorResponse.Message)
		}
	}

	if decodeErr := rpc.Converter.Unmarshal(bytes, target); decodeErr != nil {
		rpc.Logger.Warn().Str("url", url).Err(decodeErr).Msg("JSON unmarshalling failed")
		return queryInfo, decodeErr
	}

	return queryInfo, nil
}
