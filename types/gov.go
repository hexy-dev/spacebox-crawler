package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/gogo/protobuf/proto"
)

const (
	ProposalStatusInvalid = "PROPOSAL_STATUS_INVALID"
)

type (
	// DepositParams contains the data of the deposit parameters of the x/gov module
	DepositParams struct {
		MinDeposit       Coins `json:"min_deposit,omitempty" yaml:"min_deposit"`
		MaxDepositPeriod int64 `json:"max_deposit_period,omitempty" yaml:"max_deposit_period"`
	}

	// VotingParams contains the voting parameters of the x/gov module
	VotingParams struct {
		VotingPeriod int64 `json:"voting_period,omitempty" yaml:"voting_period"`
	}

	// TallyParams contains the tally parameters of the x/gov module
	TallyParams struct {
		Quorum        sdk.Dec `json:"quorum,omitempty"`
		Threshold     sdk.Dec `json:"threshold,omitempty"`
		VetoThreshold sdk.Dec `json:"veto_threshold,omitempty" yaml:"veto_threshold"`
	}

	// GovParams contains the data of the x/gov module parameters
	GovParams struct {
		TallyParams   TallyParams   `json:"tally_params" yaml:"tally_params"`
		DepositParams DepositParams `json:"deposit_params" yaml:"deposit_params"`
		VotingParams  VotingParams  `json:"voting_params" yaml:"voting_params"`
		Height        int64         `json:"height" ymal:"height"`
	}

	// Proposal represents a single governance proposal
	Proposal struct {
		DepositEndTime  time.Time
		VotingStartTime time.Time
		VotingEndTime   time.Time
		SubmitTime      time.Time
		Content         govtypes.Content
		Status          string
		ProposalRoute   string
		ProposalType    string
		Proposer        string
		ProposalID      uint64
	}

	// ProposalUpdate contains the data that should be used when updating a governance proposal
	ProposalUpdate struct {
		VotingStartTime time.Time
		VotingEndTime   time.Time
		Status          string
		ProposalID      uint64
	}

	// ProposalVoteMessage contains the data of a single proposal vote
	ProposalVoteMessage struct {
		Voter      string
		ProposalID uint64
		Height     int64
		Option     govtypes.VoteOption
	}

	// ProposalDeposit contains the data of a single deposit made towards a proposal
	ProposalDeposit struct {
		Depositor  string
		Amount     Coins
		ProposalID uint64
		Height     int64
	}

	// TallyResult contains the data about the final results of a proposal
	TallyResult struct {
		ProposalID uint64
		Yes        int64
		Abstain    int64
		No         int64
		NoWithVeto int64
		Height     int64
	}

	// ProposalValidatorStatusSnapshot represents a single snapshot of the status of a validator associated
	// with a single proposal
	ProposalValidatorStatusSnapshot struct {
		ValidatorConsAddress string
		ProposalID           uint64
		ValidatorVotingPower int64
		ValidatorStatus      int
		Height               int64
		ValidatorJailed      bool
	}
)

// NewDepositParam allows to build a new DepositParams
func NewDepositParam(d govtypes.DepositParams) DepositParams {
	return DepositParams{
		MinDeposit:       NewCoinsFromCdk(d.MinDeposit),
		MaxDepositPeriod: d.MaxDepositPeriod.Nanoseconds(),
	}
}

// NewVotingParams allows to build a new VotingParams instance
func NewVotingParams(v govtypes.VotingParams) VotingParams {
	return VotingParams{
		VotingPeriod: v.VotingPeriod.Nanoseconds(),
	}
}

// NewTallyParams allows to build a new TallyParams instance
func NewTallyParams(t govtypes.TallyParams) TallyParams {
	return TallyParams{
		Quorum:        t.Quorum,
		Threshold:     t.Threshold,
		VetoThreshold: t.VetoThreshold,
	}
}

// NewGovParams allows to build a new GovParams instance
func NewGovParams(votingParams VotingParams, depositParams DepositParams, tallyParams TallyParams, height int64) *GovParams {
	return &GovParams{
		DepositParams: depositParams,
		VotingParams:  votingParams,
		TallyParams:   tallyParams,
		Height:        height,
	}
}

// NewProposal return a new Proposal instance
func NewProposal(proposalID uint64, proposalRoute, proposalType, proposer, status string, content govtypes.Content,
	submitTime, depositEndTime, votingStartTime, votingEndTime time.Time) Proposal {
	return Proposal{
		Content:         content,
		ProposalRoute:   proposalRoute,
		ProposalType:    proposalType,
		ProposalID:      proposalID,
		Status:          status,
		SubmitTime:      submitTime,
		DepositEndTime:  depositEndTime,
		VotingStartTime: votingStartTime,
		VotingEndTime:   votingEndTime,
		Proposer:        proposer,
	}
}

// Equal tells whether p and other contain the same data
func (p Proposal) Equal(other Proposal) bool {
	return p.ProposalRoute == other.ProposalRoute &&
		p.ProposalType == other.ProposalType &&
		p.ProposalID == other.ProposalID &&
		p.Content.String() == other.Content.String() &&
		p.Status == other.Status &&
		p.SubmitTime.Equal(other.SubmitTime) &&
		p.DepositEndTime.Equal(other.DepositEndTime) &&
		p.VotingStartTime.Equal(other.VotingStartTime) &&
		p.VotingEndTime.Equal(other.VotingEndTime) &&
		p.Proposer == other.Proposer
}

// NewProposalUpdate allows to build a new ProposalUpdate instance
func NewProposalUpdate(proposalID uint64, status string, votingStartTime, votingEndTime time.Time) ProposalUpdate {
	return ProposalUpdate{
		ProposalID:      proposalID,
		Status:          status,
		VotingStartTime: votingStartTime,
		VotingEndTime:   votingEndTime,
	}
}

// NewProposalDeposit return a new ProposalDeposit instance
func NewProposalDeposit(proposalID uint64, depositor string, amount sdk.Coins, height int64) ProposalDeposit {
	return ProposalDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     NewCoinsFromCdk(amount),
		Height:     height,
	}
}

// NewProposalVoteMessage return a new ProposalVoteMessage instance
func NewProposalVoteMessage(proposalID uint64, voter string, option govtypes.VoteOption, height int64,
) ProposalVoteMessage {
	return ProposalVoteMessage{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
		Height:     height,
	}
}

// NewTallyResult return a new TallyResult instance
func NewTallyResult(proposalID uint64, yes, abstain, no, noWithVeto, height int64) TallyResult {
	return TallyResult{
		ProposalID: proposalID,
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
		Height:     height,
	}
}

//// ProposalStakingPoolSnapshot contains the data about a single staking pool snapshot to be associated with a proposal
// type ProposalStakingPoolSnapshot struct {
//	ProposalID uint64
//	Pool       *Pool
// }
//
//// NewProposalStakingPoolSnapshot returns a new ProposalStakingPoolSnapshot instance
// func NewProposalStakingPoolSnapshot(proposalID uint64, pool *Pool) ProposalStakingPoolSnapshot {
//	return ProposalStakingPoolSnapshot{
//		ProposalID: proposalID,
//		Pool:       pool,
//	}
// }

// NewProposalValidatorStatusSnapshot returns a new ProposalValidatorStatusSnapshot instance
func NewProposalValidatorStatusSnapshot(proposalID uint64, validatorConsAddr string, validatorStatus int,
	validatorJailed bool, validatorVotingPower, height int64) ProposalValidatorStatusSnapshot {
	return ProposalValidatorStatusSnapshot{
		ProposalID:           proposalID,
		ValidatorStatus:      validatorStatus,
		ValidatorConsAddress: validatorConsAddr,
		ValidatorVotingPower: validatorVotingPower,
		ValidatorJailed:      validatorJailed,
		Height:               height,
	}
}

func GetProposalContentBytes(content govtypes.Content, cdc codec.Codec) ([]byte, error) {
	// Encode the content properly
	protoContent, ok := content.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("invalid proposal content types: %T", content)
	}

	anyContent, err := codectypes.NewAnyWithValue(protoContent)
	if err != nil {
		return nil, err
	}

	contentBz, err := cdc.MarshalJSON(anyContent)
	if err != nil {
		return nil, err
	}
	return contentBz, nil
}
