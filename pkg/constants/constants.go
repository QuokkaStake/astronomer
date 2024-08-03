package constants

import (
	"errors"
	"fmt"
)

type FetcherName string

const (
	ValidatorStatusBonded = "BOND_STATUS_BONDED"
)

var (
	ErrWrongInvocation = errors.New("wrong invocation")
	ErrChainNotFound   = fmt.Errorf("chain not found")
)
