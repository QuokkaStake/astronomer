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

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
)

type DataFetcher struct {
	Logger         zerolog.Logger
	Database       *database.Database
	MetricsManager *metrics.Manager
	PriceFetchers  map[constants.PriceFetcherName]priceFetcher.PriceFetcher
	Cache          *cache.Cache
	RPCs           map[string]*tendermint.RPC
	registry       codecTypes.InterfaceRegistry
	parseCodec     *codec.ProtoCodec
}

func NewDataFetcher(
	logger zerolog.Logger,
	database *database.Database,
	metricsManager *metrics.Manager,
) *DataFetcher {
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	parseCodec := codec.NewProtoCodec(interfaceRegistry)

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
		registry:       interfaceRegistry,
		parseCodec:     parseCodec,
	}
}

func (f *DataFetcher) GetRPC(chain *types.Chain) *tendermint.RPC {
	if rpc, ok := f.RPCs[chain.Name]; ok {
		return rpc
	}

	f.RPCs[chain.Name] = tendermint.NewRPC(chain, constants.RPCQueryTimeout, f.Logger, f.MetricsManager)
	return f.RPCs[chain.Name]
}
