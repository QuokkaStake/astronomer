package tendermint

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/pkg/http"
	"main/pkg/metrics"
	"main/pkg/types"
	"main/pkg/utils"
	"strconv"
	"strings"
	"time"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govV1beta1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/gogoproto/proto"

	"github.com/rs/zerolog"
)

type RPC struct {
	Chain          *types.Chain
	Client         *http.Client
	Timeout        int
	Logger         zerolog.Logger
	MetricsManager *metrics.Manager

	registry   codecTypes.InterfaceRegistry
	parseCodec *codec.ProtoCodec
}

func NewRPC(
	chain *types.Chain,
	timeout int,
	logger zerolog.Logger,
	metricsManager *metrics.Manager,
) *RPC {
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	parseCodec := codec.NewProtoCodec(interfaceRegistry)

	return &RPC{
		Chain:   chain,
		Client:  http.NewClient(logger, chain.Name),
		Timeout: timeout,
		Logger: logger.With().
			Str("component", "rpc").
			Str("chain", chain.Name).
			Logger(),
		MetricsManager: metricsManager,
		registry:       interfaceRegistry,
		parseCodec:     parseCodec,
	}
}

func (rpc *RPC) GetAllValidators() (*stakingTypes.QueryValidatorsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/validators?pagination.count_total=true&pagination.limit=1000"

	var response stakingTypes.QueryValidatorsResponse
	info, err := rpc.Get(url, "validators", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetAllSigningInfos() (*slashingTypes.QuerySigningInfosResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/slashing/v1beta1/signing_infos?pagination.limit=1000"

	var response slashingTypes.QuerySigningInfosResponse
	info, err := rpc.Get(url, "signing_infos", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetValidator(address string) (*stakingTypes.QueryValidatorResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/validators/" + address

	var response stakingTypes.QueryValidatorResponse
	info, err := rpc.Get(url, "validator", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetStakingParams() (*stakingTypes.QueryParamsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/params"

	var response stakingTypes.QueryParamsResponse
	info, err := rpc.Get(url, "staking_params", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetSlashingParams() (*slashingTypes.QueryParamsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/slashing/v1beta1/params"

	var response slashingTypes.QueryParamsResponse
	info, err := rpc.Get(url, "slashing_params", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetGovParams(paramsType string) (*govV1beta1Types.QueryParamsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/gov/v1beta1/params/" + paramsType

	var response govV1beta1Types.QueryParamsResponse
	info, err := rpc.Get(url, "gov_params_"+paramsType, &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetMintParams() (*mintTypes.QueryParamsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/mint/v1beta1/params"

	var response mintTypes.QueryParamsResponse
	info, err := rpc.Get(url, "mint_params", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetInflation() (*mintTypes.QueryInflationResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/mint/v1beta1/inflation"

	var response mintTypes.QueryInflationResponse
	info, err := rpc.Get(url, "inflation", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetBalance(address string) (*bankTypes.QueryAllBalancesResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/bank/v1beta1/balances/" + address

	var response bankTypes.QueryAllBalancesResponse
	info, err := rpc.Get(url, "balance", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetRewards(address string) (*distributionTypes.QueryDelegationTotalRewardsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/distribution/v1beta1/delegators/" + address + "/rewards"

	var response distributionTypes.QueryDelegationTotalRewardsResponse
	info, err := rpc.Get(url, "rewards", &response)
	if err != nil {
		return nil, info, err
	}

	return &response, info, nil
}

func (rpc *RPC) GetCommission(address string) (*distributionTypes.QueryValidatorCommissionResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/distribution/v1beta1/validators/" + address + "/commission"

	var response distributionTypes.QueryValidatorCommissionResponse
	info, err := rpc.Get(url, "commission", &response)
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

func (rpc *RPC) GetDelegations(address string) (*types.DelegationsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/delegations/" + address + "?pagination.limit=1000"

	var response *types.DelegationsResponse
	info, err := rpc.GetOld(url, "delegations", &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.DelegationsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetRedelegations(address string) (*types.RedelegationsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/delegators/" + address + "/redelegations?pagination.limit=1000"

	var response *types.RedelegationsResponse
	info, err := rpc.GetOld(url, "commission", &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.RedelegationsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetUnbonds(address string) (*types.UnbondsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/delegators/" + address + "/unbonding_delegations?pagination.limit=1000"

	var response *types.UnbondsResponse
	info, err := rpc.GetOld(url, "unbonds", &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.UnbondsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetPool() (*types.PoolResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/pool"

	var response *types.PoolResponse
	info, err := rpc.GetOld(url, "pool", &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.PoolResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetSupply() (*types.SupplyResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/bank/v1beta1/supply?pagination.limit=10000&pagination.offset=0"

	var response *types.SupplyResponse
	info, err := rpc.GetOld(url, "supply", &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.SupplyResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetCommunityPool() (*types.CommunityPoolResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/distribution/v1beta1/community_pool?pagination.limit=10000&pagination.offset=0"

	var response *types.CommunityPoolResponse
	info, err := rpc.GetOld(url, "community_pool", &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.CommunityPoolResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetBlockTime() (time.Duration, error) {
	var newerBlock *types.BlockResponse
	_, err := rpc.GetOld(rpc.Chain.LCDEndpoint+"/cosmos/base/tendermint/v1beta1/blocks/latest", "block", &newerBlock)
	if err != nil {
		return 0, err
	}

	if newerBlock.Code != 0 {
		return 0, fmt.Errorf("expected code 0, but got %d: %s", newerBlock.Code, newerBlock.Message)
	}

	newerHeight := newerBlock.Block.Header.Height - 1000

	var olderBlock *types.BlockResponse
	_, err = rpc.GetOld(
		rpc.Chain.LCDEndpoint+"/cosmos/base/tendermint/v1beta1/blocks/"+strconv.FormatInt(newerHeight, 10),
		"block",
		&olderBlock,
	)
	if err != nil {
		return 0, err
	}

	if olderBlock.Code != 0 {
		return 0, fmt.Errorf("expected code 0, but got %d: %s", olderBlock.Code, olderBlock.Message)
	}

	timeDiff := olderBlock.Block.Header.Time.Sub(newerBlock.Block.Header.Time)
	heightDiff := olderBlock.Block.Header.Height - newerBlock.Block.Header.Height

	return time.Duration(float64(timeDiff.Nanoseconds()) / float64(heightDiff)), nil
}

func (rpc *RPC) GetActiveProposals() ([]types.Proposal, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/gov/v1/proposals?pagination.limit=1000&proposal_status=PROPOSAL_STATUS_VOTING_PERIOD"

	var response *types.ProposalsV1Response
	info, err := rpc.GetOld(url, "proposals_v1", &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code == 12 { // Not implemented, falling back to v1beta1
		rpc.Logger.Warn().Msg("v1 proposals are not supported, falling back to v1")

		url = rpc.Chain.LCDEndpoint + "/cosmos/gov/v1beta1/proposals?pagination.limit=1000&proposal_status=2"

		var response *types.ProposalsV1Beta1Response
		info, err := rpc.GetOld(url, "proposals_v1beta1", &response)
		if err != nil {
			return nil, info, err
		}

		if response.Code != 0 {
			return []types.Proposal{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
		}

		return utils.Map(response.Proposals, func(p types.ProposalV1Beta1) types.Proposal {
			return p.ToProposal()
		}), info, nil
	}

	if response.Code != 0 {
		return []types.Proposal{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return utils.Map(response.Proposals, func(p types.ProposalV1) types.Proposal {
		return p.ToProposal()
	}), info, nil
}

func (rpc *RPC) GetSingleProposal(proposalID string) (*types.Proposal, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/gov/v1/proposals/" + proposalID

	var response *types.ProposalV1Response
	info, err := rpc.GetOld(url, "proposal_v1", &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code == 5 { // proposal xxx doesn't exist
		return nil, info, nil
	}

	if response.Code == 12 { // Not implemented, falling back to v1beta1
		rpc.Logger.Warn().Msg("v1 proposal are not supported, falling back to v1")

		url = rpc.Chain.LCDEndpoint + "/cosmos/gov/v1beta1/proposals/" + proposalID

		var response *types.ProposalV1Beta1Response
		info, err := rpc.GetOld(url, "proposal_v1beta1", &response)
		if err != nil {
			return nil, info, err
		}

		if response.Code == 5 { // proposal xxx doesn't exist
			return nil, info, nil
		}

		if response.Code != 0 {
			return nil, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
		}

		proposal := response.Proposal.ToProposal()
		return &proposal, info, nil
	}

	if response.Code != 0 {
		return nil, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	proposal := response.Proposal.ToProposal()
	return &proposal, info, nil
}

func (rpc *RPC) GetOld(
	url string,
	queryName string,
	target interface{},
) (types.QueryInfo, error) {
	info, err := rpc.Client.Get(url, queryName, target)
	rpc.MetricsManager.LogQueryInfo(info)

	if err != nil {
		return info, err
	}

	return info, err
}

func (rpc *RPC) Get(
	url string,
	queryName string,
	target proto.Message,
) (types.QueryInfo, error) {
	bytes, queryInfo, err := rpc.Client.GetPlain(
		url,
		queryName,
	)
	if err != nil {
		rpc.Logger.Warn().Str("url", url).Err(err).Msg("LCD request failed")
		return queryInfo, err
	}

	// check whether the response is error first
	var errorResponse types.LCDError
	if err := json.Unmarshal(bytes, &errorResponse); err == nil {
		// if we successfully unmarshalled it into LCDError, so err == nil,
		// that means the response is indeed an error.
		if errorResponse.Code != 0 {
			rpc.Logger.Warn().Str("url", url).
				Err(err).
				Int("code", errorResponse.Code).
				Str("message", errorResponse.Message).
				Msg("LCD request returned an error")
			return queryInfo, errors.New(errorResponse.Message)
		}
	}

	if decodeErr := rpc.parseCodec.UnmarshalJSON(bytes, target); decodeErr != nil {
		rpc.Logger.Warn().Str("url", url).Err(decodeErr).Msg("JSON unmarshalling failed")
		return queryInfo, decodeErr
	}

	return queryInfo, nil
}
