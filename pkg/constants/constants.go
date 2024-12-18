package constants

import (
	"errors"
	"fmt"
)

type FetcherName string
type PriceFetcherName string

const (
	ValidatorStatusBonded = "BOND_STATUS_BONDED"

	CoingeckoBaseCurrency = "usd"

	PriceFetcherNameCoingecko = "coingecko"

	PrometheusMetricsPrefix = "astronomer_"

	RPCQueryTimeout = 10
	RetriesCount    = 3
)

var (
	ErrWrongInvocation = errors.New("wrong invocation")
	ErrChainNotFound   = fmt.Errorf("chain not found")
	ErrChainNotBound   = fmt.Errorf("chain not bound to this chat")
	ErrLCDNotFound     = fmt.Errorf("chain LCD host not found")
)
