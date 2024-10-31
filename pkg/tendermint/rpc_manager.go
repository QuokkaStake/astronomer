package tendermint

import (
	"main/pkg/constants"
	converterPkg "main/pkg/converter"
	databasePkg "main/pkg/database"
	"main/pkg/metrics"
	"main/pkg/types"
	"sync"
	"time"

	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	govV1beta1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/rs/zerolog"
)

type NodeManager struct {
	Logger         zerolog.Logger
	Database       *databasePkg.Database
	Converter      *converterPkg.Converter
	MetricsManager *metrics.Manager
	RPCs           map[string]*RPC

	mutex sync.Mutex
}

func NewNodeManager(
	logger *zerolog.Logger,
	database *databasePkg.Database,
	converter *converterPkg.Converter,
	metricsManager *metrics.Manager,
) *NodeManager {
	return &NodeManager{
		Logger:         logger.With().Str("component", "node_manager").Logger(),
		Database:       database,
		Converter:      converter,
		MetricsManager: metricsManager,
		RPCs:           map[string]*RPC{},
	}
}

func (manager *NodeManager) GetRPC(chain *types.Chain) *RPC {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if rpc, ok := manager.RPCs[chain.Name]; ok {
		return rpc
	}

	rpc := NewRPC(
		chain,
		constants.RPCQueryTimeout,
		&manager.Logger,
		manager.Converter,
		manager.MetricsManager,
	)
	manager.RPCs[chain.Name] = rpc
	return rpc
}

func (manager *NodeManager) GetAllValidators(chain *types.Chain) (*stakingTypes.QueryValidatorsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetAllValidators(hosts)
	return response, err
}

func (manager *NodeManager) GetAllSigningInfos(chain *types.Chain) (*slashingTypes.QuerySigningInfosResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetAllSigningInfos(hosts)
	return response, err
}

func (manager *NodeManager) GetValidator(chain *types.Chain, address string) (*stakingTypes.QueryValidatorResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetValidator(address, hosts)
	return response, err
}

func (manager *NodeManager) GetStakingParams(chain *types.Chain) (*stakingTypes.QueryParamsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetStakingParams(hosts)
	return response, err
}

func (manager *NodeManager) GetSlashingParams(chain *types.Chain) (*slashingTypes.QueryParamsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetSlashingParams(hosts)
	return response, err
}

func (manager *NodeManager) GetGovParams(chain *types.Chain, paramsType string) (*govV1beta1Types.QueryParamsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetGovParams(paramsType, hosts)
	return response, err
}

func (manager *NodeManager) GetMintParams(chain *types.Chain) (*mintTypes.QueryParamsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetMintParams(hosts)
	return response, err
}

func (manager *NodeManager) GetInflation(chain *types.Chain) (*mintTypes.QueryInflationResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetInflation(hosts)
	return response, err
}

func (manager *NodeManager) GetBalance(chain *types.Chain, address string) (*bankTypes.QueryAllBalancesResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetBalance(address, hosts)
	return response, err
}

func (manager *NodeManager) GetRewards(chain *types.Chain, address string) (*distributionTypes.QueryDelegationTotalRewardsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetRewards(address, hosts)
	return response, err
}

func (manager *NodeManager) GetCommission(chain *types.Chain, address string) (*distributionTypes.QueryValidatorCommissionResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetCommission(address, hosts)
	return response, err
}

func (manager *NodeManager) GetDelegations(chain *types.Chain, address string) (*stakingTypes.QueryDelegatorDelegationsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetDelegations(address, hosts)
	return response, err
}

func (manager *NodeManager) GetRedelegations(chain *types.Chain, address string) (*stakingTypes.QueryRedelegationsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetRedelegations(address, hosts)
	return response, err
}

func (manager *NodeManager) GetUnbonds(chain *types.Chain, address string) (*stakingTypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetUnbonds(address, hosts)
	return response, err
}

func (manager *NodeManager) GetPool(chain *types.Chain) (*stakingTypes.QueryPoolResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetPool(hosts)
	return response, err
}

func (manager *NodeManager) GetSupply(chain *types.Chain) (*bankTypes.QueryTotalSupplyResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetSupply(hosts)
	return response, err
}

func (manager *NodeManager) GetCommunityPool(chain *types.Chain) (*distributionTypes.QueryCommunityPoolResponse, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetCommunityPool(hosts)
	return response, err
}

func (manager *NodeManager) GetBlockTime(chain *types.Chain) (time.Duration, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return 0, err
	}

	rpc := manager.GetRPC(chain)
	response, err := rpc.GetBlockTime(hosts)
	return response, err
}

func (manager *NodeManager) GetActiveProposals(chain *types.Chain) ([]*types.Proposal, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return []*types.Proposal{}, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetActiveProposals(hosts)
	return response, err
}

func (manager *NodeManager) GetSingleProposal(chain *types.Chain, id string) (*types.Proposal, error) {
	hosts, err := manager.Database.GetLCDHosts(chain)
	if err != nil {
		return nil, err
	}

	rpc := manager.GetRPC(chain)
	response, _, err := rpc.GetSingleProposal(id, hosts)
	return response, err
}
