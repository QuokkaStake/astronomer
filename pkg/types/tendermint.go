package types

import (
	"time"

	"cosmossdk.io/math"
)

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

type PoolResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Pool    struct {
		BondedTokens    math.LegacyDec `json:"bonded_tokens"`
		NotBondedTokens math.LegacyDec `json:"not_bonded_tokens"`
	} `json:"pool"`
}

type SupplyResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Supply  []SdkAmount `json:"supply"`
}

type CommunityPoolResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Pool    []SdkAmount `json:"pool"`
}

type LCDError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
