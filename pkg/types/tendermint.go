package types

import (
	"main/pkg/constants"
	"time"

	"cosmossdk.io/math"
)

type ValidatorResponse struct {
	Code      int       `json:"code"`
	Validator Validator `json:"validator"`
}

type ValidatorDescription struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type ValidatorCommission struct {
	CommissionRates ValidatorCommissionRates `json:"commission_rates"`
	UpdateTime      time.Time                `json:"update_time"`
}

type ValidatorCommissionRates struct {
	Rate          math.LegacyDec `json:"rate"`
	MaxRate       math.LegacyDec `json:"max_rate"`
	MaxChangeRate math.LegacyDec `json:"max_change_rate"`
}

type Validator struct {
	OperatorAddress   string               `json:"operator_address"`
	ConsensusPubkey   ConsensusPubkey      `json:"consensus_pubkey"`
	Jailed            bool                 `json:"jailed"`
	Status            string               `json:"status"`
	Tokens            string               `json:"tokens"`
	DelegatorShares   math.LegacyDec       `json:"delegator_shares"`
	Description       ValidatorDescription `json:"description"`
	UnbondingHeight   string               `json:"unbonding_height"`
	UnbondingTime     time.Time            `json:"unbonding_time"`
	Commission        ValidatorCommission  `json:"commission"`
	MinSelfDelegation string               `json:"min_self_delegation"`
}

func (v Validator) Active() bool {
	return v.Status == constants.ValidatorStatusBonded
}

type ConsensusPubkey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type PaginationResponse struct {
	Code       int        `json:"code"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total uint64 `json:"total,string"`
}

type ValidatorsResponse struct {
	Code       int          `json:"code"`
	Validators []*Validator `json:"validators"`
}

type ValidatorInfo struct {
	Chain     *Chain
	Validator *Validator
	Error     error
}
