package tendermint

import (
	"fmt"
	"main/pkg/http"
	"main/pkg/types"
	"main/pkg/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type RPC struct {
	Chain   *types.Chain
	Client  *http.Client
	Timeout int
	Logger  zerolog.Logger

	LastHeight map[string]int64
	Mutex      sync.Mutex
}

func NewRPC(
	chain *types.Chain,
	timeout int,
	logger zerolog.Logger,
) *RPC {
	return &RPC{
		Chain:   chain,
		Client:  http.NewClient(logger, chain.Name),
		Timeout: timeout,
		Logger: logger.With().
			Str("component", "rpc").
			Str("chain", chain.Name).
			Logger(),
		LastHeight: map[string]int64{},
	}
}

func (rpc *RPC) GetAllValidators() (*types.ValidatorsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/validators?pagination.count_total=true&pagination.limit=1000"

	var response *types.ValidatorsResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.ValidatorsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetValidator(address string) (*types.ValidatorResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/validators/" + address

	var response *types.ValidatorResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.ValidatorResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetStakingParams() (*types.StakingParamsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/params"

	var response *types.StakingParamsResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.StakingParamsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetSlashingParams() (*types.SlashingParamsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/slashing/v1beta1/params"

	var response *types.SlashingParamsResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.SlashingParamsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetGovParams(paramsType string) (*types.GovParamsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/gov/v1beta1/params/" + paramsType

	var response *types.GovParamsResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.GovParamsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetMintParams() (*types.MintParamsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/mint/v1beta1/params"

	var response *types.MintParamsResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.MintParamsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetInflation() (*types.InflationResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/mint/v1beta1/inflation"

	var response *types.InflationResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.InflationResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetBalance(address string) (*types.BalancesResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/bank/v1beta1/balances/" + address

	var response *types.BalancesResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.BalancesResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetRewards(address string) (*types.RewardsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/distribution/v1beta1/delegators/" + address + "/rewards"

	var response *types.RewardsResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.RewardsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetCommission(address string) (*types.CommissionsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/distribution/v1beta1/validators/" + address + "/commission"

	var response *types.CommissionsResponse
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		// not being a validator is acceptable
		if strings.Contains(response.Message, "validator does not exist") {
			return &types.CommissionsResponse{
				Commission: types.SdkCommission{
					Commission: make([]types.SdkAmount, 0),
				},
			}, info, nil
		}

		return &types.CommissionsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetDelegations(address string) (*types.DelegationsResponse, types.QueryInfo, error) {
	url := rpc.Chain.LCDEndpoint + "/cosmos/staking/v1beta1/delegations/" + address + "?pagination.limit=1000"

	var response *types.DelegationsResponse
	info, err := rpc.Get(url, &response)
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
	info, err := rpc.Get(url, &response)
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
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code != 0 {
		return &types.UnbondsResponse{}, info, fmt.Errorf("expected code 0, but got %d: %s", response.Code, response.Message)
	}

	return response, info, nil
}

func (rpc *RPC) GetBlockTime() (time.Duration, error) {
	var newerBlock *types.BlockResponse
	_, err := rpc.Get(rpc.Chain.LCDEndpoint+"/cosmos/base/tendermint/v1beta1/blocks/latest", &newerBlock)
	if err != nil {
		return 0, err
	}

	if newerBlock.Code != 0 {
		return 0, fmt.Errorf("expected code 0, but got %d: %s", newerBlock.Code, newerBlock.Message)
	}

	newerHeight := newerBlock.Block.Header.Height - 1000

	var olderBlock *types.BlockResponse
	_, err = rpc.Get(
		rpc.Chain.LCDEndpoint+"/cosmos/base/tendermint/v1beta1/blocks/"+strconv.FormatInt(newerHeight, 10),
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
	info, err := rpc.Get(url, &response)
	if err != nil {
		return nil, info, err
	}

	if response.Code == 12 { // Not implemented, falling back to v1beta1
		rpc.Logger.Warn().Msg("v1 proposals are not supported, falling back to v1")

		url = rpc.Chain.LCDEndpoint + "/cosmos/gov/v1beta1/proposals?pagination.limit=1000&proposal_status=2"

		var response *types.ProposalsV1Beta1Response
		info, err := rpc.Get(url, &response)
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
	info, err := rpc.Get(url, &response)
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
		info, err := rpc.Get(url, &response)
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

func (rpc *RPC) Get(
	url string,
	target interface{},
) (types.QueryInfo, error) {
	info, err := rpc.Client.Get(url, target)

	if err != nil {
		return info, err
	}

	return info, err
}
