package pricefetcher

import (
	"fmt"
	"main/pkg/constants"
	"main/pkg/http"
	"main/pkg/metrics"
	"main/pkg/types"
	"main/pkg/utils"
	"strings"

	"github.com/rs/zerolog"
)

type CoingeckoPriceFetcher struct {
	Client         *http.Client
	Logger         zerolog.Logger
	MetricsManager *metrics.Manager
}

func NewCoingeckoPriceFetcher(
	logger *zerolog.Logger,
	metricsManager *metrics.Manager,
) *CoingeckoPriceFetcher {
	return &CoingeckoPriceFetcher{
		Client:         http.NewClient(logger, "coingecko"),
		Logger:         logger.With().Str("component", "coingecko_price_fetcher").Logger(),
		MetricsManager: metricsManager,
	}
}

func (c *CoingeckoPriceFetcher) GetPrices(denomInfos []*types.Denom) (Prices, error) {
	currenciesToFetch := utils.Map(denomInfos, func(denomInfo *types.Denom) string {
		return denomInfo.CoingeckoCurrency.String
	})

	var coingeckoResponse map[string]map[string]float64
	queryInfo, err := c.Client.Get(
		"https://api.coingecko.com",
		fmt.Sprintf(
			"/api/v3/simple/price?ids=%s&vs_currencies=%s",
			strings.Join(currenciesToFetch, ","),
			constants.CoingeckoBaseCurrency,
		),
		"fetch_prices",
		&coingeckoResponse,
	)

	c.MetricsManager.LogQueryInfo(queryInfo)

	if err != nil {
		c.Logger.Error().
			Err(err).
			Strs("currencies", currenciesToFetch).
			Msg("Could not get rates, probably rate-limiting")
		return map[string]map[string]float64{}, err
	}

	result := make(map[string]map[string]float64)

	for _, denomInfo := range denomInfos {
		if _, ok := result[denomInfo.Chain]; !ok {
			result[denomInfo.Chain] = make(map[string]float64)
		}

		coinPrice, ok := coingeckoResponse[denomInfo.CoingeckoCurrency.String]
		if !ok {
			continue
		}

		if usdCoinPrice, ok := coinPrice[constants.CoingeckoBaseCurrency]; ok {
			result[denomInfo.Chain][denomInfo.Denom] = usdCoinPrice
		}
	}

	return result, nil
}

func (c *CoingeckoPriceFetcher) Name() string {
	return "coingecko"
}
