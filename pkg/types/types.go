package types

import (
	"fmt"
	"main/pkg/constants"
	"time"

	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"cosmossdk.io/math"
)

type QueryInfo struct {
	Chain    string
	URL      string
	Query    string
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

	SigningInfo *slashingTypes.ValidatorSigningInfo
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
	Infos map[string]*ChainWalletsBalancesInfo
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

func (w *WalletsBalancesInfo) SetChain(chain *Chain, explorers []*Explorer) {
	w.Infos[chain.Name] = &ChainWalletsBalancesInfo{
		Chain:        chain,
		Explorers:    explorers,
		BalancesInfo: map[string]*WalletBalancesInfo{},
	}
}

func (w *WalletsBalancesInfo) SetAddressInfo(chainName string, address *WalletLink) {
	if _, ok := w.Infos[chainName].BalancesInfo[address.Address]; !ok {
		w.Infos[chainName].BalancesInfo[address.Address] = &WalletBalancesInfo{
			Address: address,
		}
	}
}

func (w *WalletsBalancesInfo) SetBalancesError(chainName string, address *WalletLink, err error) {
	w.Infos[chainName].BalancesInfo[address.Address].BalancesError = err
}

func (w *WalletsBalancesInfo) SetBalances(chainName string, address *WalletLink, balances []*Amount) {
	w.Infos[chainName].BalancesInfo[address.Address].Balances = balances
}

func (w *WalletsBalancesInfo) SetRewardsError(chainName string, address *WalletLink, err error) {
	w.Infos[chainName].BalancesInfo[address.Address].RewardsError = err
}

func (w *WalletsBalancesInfo) SetRewards(chainName string, address *WalletLink, rewards []*Amount) {
	w.Infos[chainName].BalancesInfo[address.Address].Rewards = rewards
}

func (w *WalletsBalancesInfo) SetCommissionsError(chainName string, address *WalletLink, err error) {
	w.Infos[chainName].BalancesInfo[address.Address].CommissionsError = err
}

func (w *WalletsBalancesInfo) SetCommissions(chainName string, address *WalletLink, commissions []*Amount) {
	w.Infos[chainName].BalancesInfo[address.Address].Commissions = commissions
}

func (w *WalletsBalancesInfo) SetDelegationsError(chainName string, address *WalletLink, err error) {
	w.Infos[chainName].BalancesInfo[address.Address].DelegationsError = err
}

func (w *WalletsBalancesInfo) SetDelegations(chainName string, address *WalletLink, delegations []*Delegation) {
	w.Infos[chainName].BalancesInfo[address.Address].Delegations = delegations
}

func (w *WalletsBalancesInfo) SetRedelegationsError(chainName string, address *WalletLink, err error) {
	w.Infos[chainName].BalancesInfo[address.Address].RedelegationsError = err
}

func (w *WalletsBalancesInfo) SetRedelegations(chainName string, address *WalletLink, redelegations []*Redelegation) {
	w.Infos[chainName].BalancesInfo[address.Address].Redelegations = redelegations
}

func (w *WalletsBalancesInfo) SetUnbondsError(chainName string, address *WalletLink, err error) {
	w.Infos[chainName].BalancesInfo[address.Address].UnbondsError = err
}

func (w *WalletsBalancesInfo) SetUnbonds(chainName string, address *WalletLink, unbonds []*Unbond) {
	w.Infos[chainName].BalancesInfo[address.Address].Unbonds = unbonds
}

type SupplyInfo struct {
	Error    error
	Supplies map[string]*ChainSupply
}

type ChainSupply struct {
	Chain              *Chain
	PoolError          error
	BondedTokens       *Amount
	NotBondedTokens    *Amount
	SupplyError        error
	AllSupplies        map[string]*Amount
	CommunityPoolError error
	AllCommunityPool   map[string]*Amount
}

func (c ChainSupply) HasBondedSupply() bool {
	if c.AllSupplies == nil {
		return false
	}

	if c.BondedTokens == nil {
		return false
	}

	_, found := c.AllSupplies[c.Chain.BaseDenom]
	return found
}

func (c ChainSupply) BondedSupplyPercent() float64 {
	baseDenomSupply := c.AllSupplies[c.Chain.BaseDenom]
	return c.BondedTokens.Amount.MustFloat64() / baseDenomSupply.Amount.MustFloat64()
}

func (c ChainSupply) HasCommunityPoolSupply() bool {
	if c.AllCommunityPool == nil || c.AllSupplies == nil {
		return false
	}

	if _, found := c.AllCommunityPool[c.Chain.BaseDenom]; !found {
		return false
	}
	if _, found := c.AllSupplies[c.Chain.BaseDenom]; !found {
		return false
	}
	return true
}

func (c ChainSupply) CommunityPoolSupplyPercent() float64 {
	baseDenomSupply := c.AllSupplies[c.Chain.BaseDenom]
	baseDenomCommunityPool := c.AllCommunityPool[c.Chain.BaseDenom]
	return baseDenomCommunityPool.Amount.MustFloat64() / baseDenomSupply.Amount.MustFloat64()
}

type WalletsList struct {
	Error error
	Infos map[string]*ChainWalletsList
}

type ChainWalletsList struct {
	Chain     *Chain
	Explorers Explorers
	Wallets   []*WalletLink
}

type ChainWallet struct {
	Chain     *Chain
	Explorers Explorers
	Wallet    *WalletLink
}
