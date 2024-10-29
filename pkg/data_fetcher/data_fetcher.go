package datafetcher

import (
	"main/pkg/cache"
	"main/pkg/constants"
	converterPkg "main/pkg/converter"
	"main/pkg/database"
	"main/pkg/metrics"
	priceFetcher "main/pkg/price_fetcher"
	"main/pkg/tendermint"
	"sync"

	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Logger         zerolog.Logger
	Database       *database.Database
	Converter      *converterPkg.Converter
	MetricsManager *metrics.Manager
	PriceFetchers  map[constants.PriceFetcherName]priceFetcher.PriceFetcher
	Cache          *cache.Cache
	RPCs           map[string]*tendermint.RPC
	NodesManager   *tendermint.NodeManager

	mutex sync.Mutex
}

func NewDataFetcher(
	logger *zerolog.Logger,
	database *database.Database,
	converter *converterPkg.Converter,
	metricsManager *metrics.Manager,
	nodesManager *tendermint.NodeManager,
) *DataFetcher {
	priceFetchers := map[constants.PriceFetcherName]priceFetcher.PriceFetcher{
		constants.PriceFetcherNameCoingecko: priceFetcher.NewCoingeckoPriceFetcher(logger),
	}

	return &DataFetcher{
		Logger:         logger.With().Str("component", "data_fetcher").Logger(),
		Database:       database,
		Converter:      converter,
		MetricsManager: metricsManager,
		PriceFetchers:  priceFetchers,
		Cache:          cache.NewCache(),
		RPCs:           map[string]*tendermint.RPC{},
		NodesManager:   nodesManager,
	}
}
