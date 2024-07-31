package types

import (
	"time"
)

type QueryInfo struct {
	Chain    string
	URL      string
	Duration time.Duration
	Success  bool
}

type Amount struct {
	Amount float64
	Denom  string
}

type Proposal struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	VotingStartTime time.Time `json:"voting_start_time"`
	VotingEndTime   time.Time `json:"voting_end_time"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
}

func (p Proposal) VotingTimeLeft() time.Duration {
	return time.Until(p.VotingEndTime)
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
	Chain          *Chain
	Proposals      []Proposal
	ProposalsError error
}
