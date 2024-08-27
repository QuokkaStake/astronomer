package datafetcher

import (
	"main/pkg/cache"
	"main/pkg/constants"
	"main/pkg/database"
	"main/pkg/metrics"
	priceFetcher "main/pkg/price_fetcher"
	"main/pkg/tendermint"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Logger         zerolog.Logger
	Database       *database.Database
	MetricsManager *metrics.Manager
	PriceFetchers  map[constants.PriceFetcherName]priceFetcher.PriceFetcher
	Cache          *cache.Cache
	RPCs           map[string]*tendermint.RPC
}

func NewDataFetcher(
	logger zerolog.Logger,
	database *database.Database,
	metricsManager *metrics.Manager,
) *DataFetcher {
	priceFetchers := map[constants.PriceFetcherName]priceFetcher.PriceFetcher{
		constants.PriceFetcherNameCoingecko: priceFetcher.NewCoingeckoPriceFetcher(logger),
	}

	return &DataFetcher{
		Logger:         logger.With().Str("component", "data_fetcher").Logger(),
		Database:       database,
		MetricsManager: metricsManager,
		PriceFetchers:  priceFetchers,
		Cache:          cache.NewCache(),
		RPCs:           map[string]*tendermint.RPC{},
	}
}

func (f *DataFetcher) GetRPC(chain *types.Chain) *tendermint.RPC {
	if rpc, ok := f.RPCs[chain.Name]; ok {
		return rpc
	}

	f.RPCs[chain.Name] = tendermint.NewRPC(chain, constants.RPCQueryTimeout, f.Logger, f.MetricsManager)
	return f.RPCs[chain.Name]
}
