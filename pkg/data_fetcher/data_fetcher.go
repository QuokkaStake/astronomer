package datafetcher

import (
	"main/pkg/cache"
	"main/pkg/constants"
	"main/pkg/database"
	priceFetcher "main/pkg/price_fetcher"

	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Logger        zerolog.Logger
	Database      *database.Database
	PriceFetchers map[constants.PriceFetcherName]priceFetcher.PriceFetcher
	Cache         *cache.Cache
}

func NewDataFetcher(logger zerolog.Logger, database *database.Database) *DataFetcher {
	priceFetchers := map[constants.PriceFetcherName]priceFetcher.PriceFetcher{
		constants.PriceFetcherNameCoingecko: priceFetcher.NewCoingeckoPriceFetcher(logger),
	}

	return &DataFetcher{
		Logger:        logger.With().Str("component", "data_fetcher").Logger(),
		Database:      database,
		PriceFetchers: priceFetchers,
		Cache:         cache.NewCache(),
	}
}
