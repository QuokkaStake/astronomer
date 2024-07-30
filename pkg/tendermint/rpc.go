package tendermint

import (
	"fmt"
	"main/pkg/http"
	"main/pkg/types"
	"strconv"
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
	logger *zerolog.Logger,
) *RPC {
	return &RPC{
		Chain:   chain,
		Client:  http.NewClient(logger, chain.GetName()),
		Timeout: timeout,
		Logger: logger.With().
			Str("component", "rpc").
			Str("chain", chain.GetName()).
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
		return &types.ValidatorsResponse{}, info, fmt.Errorf("expected code 0, but got %d", response.Code)
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
		return &types.StakingParamsResponse{}, info, fmt.Errorf("expected code 0, but got %d", response.Code)
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
		return &types.SlashingParamsResponse{}, info, fmt.Errorf("expected code 0, but got %d", response.Code)
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
		return &types.GovParamsResponse{}, info, fmt.Errorf("expected code 0, but got %d", response.Code)
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
		return 0, fmt.Errorf("expected code 0, but got %d", newerBlock.Code)
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
		return 0, fmt.Errorf("expected code 0, but got %d", olderBlock.Code)
	}

	timeDiff := olderBlock.Block.Header.Time.Sub(newerBlock.Block.Header.Time)
	heightDiff := olderBlock.Block.Header.Height - newerBlock.Block.Header.Height

	return time.Duration(float64(timeDiff.Nanoseconds()) / float64(heightDiff)), nil
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
