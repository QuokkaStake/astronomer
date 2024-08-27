package types

import (
	"main/pkg/constants"
	"main/pkg/utils"
	"sort"
	"time"

	"cosmossdk.io/math"
)

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
	Message    string     `json:"message"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total uint64 `json:"total,string"`
}

type ValidatorsResponse struct {
	Code       int          `json:"code"`
	Message    string       `json:"message"`
	Validators []*Validator `json:"validators"`
}

type ValidatorResponse struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Validator Validator `json:"validator"`
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

type StakingParamsResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Params  StakingParams `json:"params"`
}

type StakingParams struct {
	UnbondingTime Duration `json:"unbonding_time"`
	Message       string   `json:"message"`
	MaxValidators int      `json:"max_validators"`
}

type SlashingParamsResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Params  SlashingParams `json:"params"`
}

type SlashingParams struct {
	SignedBlocksWindow      int      `json:"signed_blocks_window,string"`
	MinSignedPerWindow      float64  `json:"min_signed_per_window,string"`
	DowntimeJailDuration    Duration `json:"downtime_jail_duration"`
	SlashFractionDowntime   float64  `json:"slash_fraction_downtime,string"`
	SlashFractionDoubleSign float64  `json:"slash_fraction_double_sign,string"`
}

type GovParamsResponse struct {
	Code          int           `json:"code"`
	Message       string        `json:"message"`
	VotingParams  VotingParams  `json:"voting_params"`
	DepositParams DepositParams `json:"deposit_params"`
	TallyParams   TallyParams   `json:"tally_params"`
}

type VotingParams struct {
	VotingPeriod Duration `json:"voting_period"`
}

type DepositParams struct {
	MaxDepositPeriod Duration `json:"max_deposit_period"`
}

type TallyParams struct {
	Quorum        float64 `json:"quorum,string"`
	Threshold     float64 `json:"threshold,string"`
	VetoThreshold float64 `json:"veto_threshold,string"`
}

type BlockResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Block   struct {
		Header struct {
			Time   time.Time `json:"time"`
			Height int64     `json:"height,string"`
		} `json:"header"`
	} `json:"block"`
}

type MintParamsResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Params  MintParams `json:"params"`
}

type MintParams struct {
	InflationRateChange float64 `json:"inflation_rate_change,string"`
	InflationMax        float64 `json:"inflation_max,string"`
	InflationMin        float64 `json:"inflation_min,string"`
	GoalBonded          float64 `json:"goal_bonded,string"`
	BlocksPerYear       int64   `json:"blocks_per_year,string"`
}

type InflationResponse struct {
	Code      int     `json:"code"`
	Message   string  `json:"message"`
	Inflation float64 `json:"inflation,string"`
}

type ProposalsV1Response struct {
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	Proposals []ProposalV1 `json:"proposals"`
}

type ProposalV1Response struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Proposal *ProposalV1 `json:"proposal"`
}

type ProposalV1 struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	VotingStartTime time.Time `json:"voting_start_time"`
	VotingEndTime   time.Time `json:"voting_end_time"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
}

func (p *ProposalV1) ToProposal() Proposal {
	return Proposal{
		ID:              p.ID,
		Status:          p.Status,
		VotingStartTime: p.VotingStartTime,
		VotingEndTime:   p.VotingEndTime,
		Title:           p.Title,
		Summary:         p.Summary,
	}
}

type ProposalsV1Beta1Response struct {
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	Proposals []ProposalV1Beta1 `json:"proposals"`
}

type ProposalV1Beta1Response struct {
	Code     int              `json:"code"`
	Message  string           `json:"message"`
	Proposal *ProposalV1Beta1 `json:"proposal"`
}

type ProposalV1Beta1 struct {
	ProposalID      string    `json:"proposal_id"`
	Status          string    `json:"status"`
	VotingStartTime time.Time `json:"voting_start_time"`
	VotingEndTime   time.Time `json:"voting_end_time"`
	Content         struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"content"`
}

func (p *ProposalV1Beta1) ToProposal() Proposal {
	return Proposal{
		ID:              p.ProposalID,
		Status:          p.Status,
		VotingStartTime: p.VotingStartTime,
		VotingEndTime:   p.VotingEndTime,
		Title:           p.Content.Title,
		Summary:         p.Content.Description,
	}
}

type SdkAmount struct {
	Amount math.LegacyDec `json:"amount"`
	Denom  string         `json:"denom"`
}

func (s SdkAmount) ToAmount() *Amount {
	return &Amount{
		Amount: s.Amount,
		Denom:  s.Denom,
	}
}

type BalancesResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Balances []SdkAmount `json:"balances"`
}

type RewardsResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Total   []SdkAmount `json:"total"`
}

type SdkDelegation struct {
	Delegation struct {
		ValidatorAddress string `json:"validator_address"`
	} `json:"delegation"`
	Balance SdkAmount `json:"balance"`
}

type DelegationsResponse struct {
	Code        int             `json:"code"`
	Message     string          `json:"message"`
	Delegations []SdkDelegation `json:"delegation_responses"`
}

type SdkUnbondEntry struct {
	CompletionTime time.Time      `json:"completion_time"`
	Balance        math.LegacyDec `json:"balance"`
}

type SdkUnbond struct {
	ValidatorAddress string           `json:"validator_address"`
	Entries          []SdkUnbondEntry `json:"entries"`
}

type UnbondsResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Unbonds []SdkUnbond `json:"unbonding_responses"`
}

type SdkRedelegationEntry struct {
	Entry struct {
		CompletionTime time.Time `json:"completion_time"`
	} `json:"redelegation_entry"`
	Balance math.LegacyDec `json:"balance"`
}

type SdkRedelegation struct {
	Redelegation struct {
		ValidatorSrcAddress string `json:"validator_src_address"`
		ValidatorDstAddress string `json:"validator_dst_address"`
	} `json:"redelegation"`
	Entries []SdkRedelegationEntry `json:"entries"`
}

type RedelegationsResponse struct {
	Code          int               `json:"code"`
	Message       string            `json:"message"`
	Redelegations []SdkRedelegation `json:"redelegation_responses"`
}

type SdkCommission struct {
	Commission []SdkAmount `json:"commission"`
}

type CommissionsResponse struct {
	Code       int           `json:"code"`
	Message    string        `json:"message"`
	Commission SdkCommission `json:"commission"`
}

type PoolResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Pool    struct {
		BondedTokens    math.LegacyDec `json:"bonded_tokens"`
		NotBondedTokens math.LegacyDec `json:"not_bonded_tokens"`
	} `json:"pool"`
}
