package types

import (
	"fmt"
	"main/pkg/constants"
	"main/pkg/utils"
	"sort"
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

func (v Validator) FormatCommission() string {
	return fmt.Sprintf("%.2f", v.Commission.CommissionRates.Rate.MustFloat64()*100)
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

func (v ValidatorsResponse) GetTotalVP() math.LegacyDec {
	vp := math.LegacyNewDec(0)

	for _, validator := range v.Validators {
		if validator.Active() {
			vp = vp.Add(validator.DelegatorShares)
		}
	}

	return vp
}

func (v ValidatorsResponse) FindValidatorRank(valoper string) int {
	valsActive := utils.Filter(v.Validators, func(v *Validator) bool {
		return v.Active()
	})

	sort.SliceStable(valsActive, func(i, j int) bool {
		return valsActive[i].DelegatorShares.GT(valsActive[j].DelegatorShares)
	})

	for index, validator := range valsActive {
		if validator.OperatorAddress == valoper {
			return index + 1
		}
	}

	return 0
}

type ValidatorInfo struct {
	Chain              *Chain
	Validator          *Validator
	VotingPowerPercent float64
	Rank               int
	Error              error
}

func (i ValidatorInfo) GetVotingPowerPercent() string {
	return fmt.Sprintf("%.2f", i.VotingPowerPercent*100)

}
