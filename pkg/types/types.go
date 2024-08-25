package types

import (
	"fmt"
	"main/pkg/constants"
	"time"

	"cosmossdk.io/math"
)

type QueryInfo struct {
	Chain    string
	URL      string
	Duration time.Duration
	Success  bool
}

type Amount struct {
	Amount    math.LegacyDec
	Denom     string
	BaseDenom string
	PriceUSD  *math.LegacyDec
}

type Proposal struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	VotingStartTime time.Time `json:"voting_start_time"`
	VotingEndTime   time.Time `json:"voting_end_time"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
}

func (p Proposal) FormatStatus() string {
	switch p.Status {
	case "PROPOSAL_STATUS_VOTING_PERIOD":
		return "📥In voting"
	case "PROPOSAL_STATUS_PASSED":
		return "🏁Passed"
	case "PROPOSAL_STATUS_REJECTED":
		return "☠️Rejected"
	default:
		return p.Status
	}
}

type ChainsParams struct {
	Error  error
	Params map[string]*ChainParams
}

type ChainParams struct {
	Chain               *Chain
	StakingParams       StakingParams
	StakingParamsError  error
	SlashingParams      SlashingParams
	SlashingParamsError error

	VotingParams       VotingParams
	VotingParamsError  error
	DepositParams      DepositParams
	DepositParamsError error
	TallyParams        TallyParams
	TallyParamsError   error

	BlockTime      time.Duration
	BlockTimeError error

	MintParams      MintParams
	MintParamsError error

	Inflation      float64
	InflationError error
}

type ActiveProposals struct {
	Error     error
	Proposals map[string]*ChainActiveProposals
}

type ChainActiveProposals struct {
	Chain          *Chain
	Explorers      Explorers
	Proposals      []Proposal
	ProposalsError error
}

type SingleProposal struct {
	Chain     *Chain
	Explorers Explorers
	Proposal  *Proposal
	Error     error
}

type ValidatorsInfo struct {
	Error  error
	Chains map[string]ChainValidatorsInfo
}

type ChainValidatorsInfo struct {
	Chain      *Chain
	Explorers  Explorers
	Error      error
	Validators []ValidatorInfo
}

type ValidatorInfo struct {
	OperatorAddress         string
	Jailed                  bool
	Status                  string
	Tokens                  *Amount
	Moniker                 string
	Details                 string
	Identity                string
	Website                 string
	SecurityContact         string
	Commission              float64
	CommissionMax           float64
	CommissionMaxChangeRate float64
	VotingPowerPercent      float64
	Rank                    int
}

func (i ValidatorInfo) Active() bool {
	return i.Status == constants.ValidatorStatusBonded
}

func (i ValidatorInfo) FormatCommission() string {
	return fmt.Sprintf("%.2f", i.Commission*100)
}

func (i ValidatorInfo) GetVotingPowerPercent() string {
	return fmt.Sprintf("%.2f", i.VotingPowerPercent*100)
}

type ValidatorAddressWithMoniker struct {
	Chain   *Chain
	Address string
	Moniker string
}

func (v *ValidatorAddressWithMoniker) GetName() string {
	if v.Moniker != "" {
		return v.Moniker
	}

	return v.Address
}

type Delegation struct {
	Amount    *Amount
	Validator *ValidatorAddressWithMoniker
}

type Redelegation struct {
	SrcValidator   *ValidatorAddressWithMoniker
	DstValidator   *ValidatorAddressWithMoniker
	Amount         *Amount
	CompletionTime time.Time
}

type Unbond struct {
	Validator      *ValidatorAddressWithMoniker
	Amount         *Amount
	CompletionTime time.Time
}

type WalletsBalancesInfo struct {
	Error error
	Infos map[string]ChainWalletsBalancesInfo
}

type ChainWalletsBalancesInfo struct {
	Chain        *Chain
	Explorers    Explorers
	BalancesInfo map[string]*WalletBalancesInfo
}

type WalletBalancesInfo struct {
	Address            *WalletLink
	Balances           []*Amount
	BalancesError      error
	Rewards            []*Amount
	RewardsError       error
	Commissions        []*Amount
	CommissionsError   error
	Delegations        []*Delegation
	DelegationsError   error
	Redelegations      []*Redelegation
	RedelegationsError error
	Unbonds            []*Unbond
	UnbondsError       error
}
