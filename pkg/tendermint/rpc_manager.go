package tendermint

import (
	"main/pkg/constants"
	converterPkg "main/pkg/converter"
	databasePkg "main/pkg/database"
	"main/pkg/metrics"
	"main/pkg/types"
	"sync"

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
	logger zerolog.Logger,
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

func (manager *NodeManager) GetRPC(chain *types.Chain) (*RPC, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if rpc, ok := manager.RPCs[chain.Name]; ok {
		return rpc, nil
	}

	rpc := NewRPC(
		chain,
		constants.RPCQueryTimeout,
		manager.Logger,
		manager.Converter,
		manager.MetricsManager,
	)
	manager.RPCs[chain.Name] = rpc
	return rpc, nil
}
