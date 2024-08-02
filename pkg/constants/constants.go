package constants

import "errors"

type FetcherName string

const (
	ValidatorStatusBonded = "BOND_STATUS_BONDED"
)

var (
	ErrWrongInvocation = errors.New("wrong invocation")
)
