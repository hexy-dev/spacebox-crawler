package types

import (
	"fmt"
	"time"

	"github.com/cometbft/cometbft/crypto"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pkg/errors"
)

var (
	ErrNoEventFound     = errors.New("no event found")
	ErrNoAttributeFound = errors.New("no event with attribute")
)

type (
	PubKey interface {
		Bytes() []byte
	}

	ValidatorPreCommit struct {
		Timestamp        time.Time
		ValidatorAddress string
		BlockIDFlag      uint64
		VotingPower      int64
	}

	Block struct {
		rb *coretypes.ResultBlock

		Timestamp           time.Time
		Hash                string
		ProposerAddress     string
		ValidatorPreCommits []ValidatorPreCommit
		Evidence            tmtypes.EvidenceData
		TxNum               int
		TotalGas            uint64
		Height              int64
	}

	// Txs - slice of transactions
	Txs []*Tx

	// Tx represents an already existing blockchain transaction
	Tx struct {
		*sdktx.Tx
		*sdk.TxResponse
		Signer string
	}

	Validators []*Validator

	Validator struct {
		ConsAddr   string
		ConsPubkey string
	}
)

func NewBlock(height int64, hash, proposerAddress string, txNum int, totalGas uint64, timestamp time.Time,
	evidence tmtypes.EvidenceData) *Block {

	return &Block{
		Height:          height,
		Hash:            hash,
		TxNum:           txNum,
		ProposerAddress: proposerAddress,
		Timestamp:       timestamp,
		Evidence:        evidence,
		TotalGas:        totalGas,
	}
}

// NewBlockFromTmBlock builds a new Block instance from a given ResultBlock object
func NewBlockFromTmBlock(blk *coretypes.ResultBlock, totalGas uint64) *Block {
	res := NewBlock(
		blk.Block.Height,
		blk.Block.Hash().String(),
		sdk.ConsAddress(blk.Block.ProposerAddress).String(),
		len(blk.Block.Txs),
		totalGas,
		blk.Block.Time,
		blk.Block.Evidence,
	)

	if blk.Block.LastCommit != nil {
		res.ValidatorPreCommits = NewValidatorPreCommitsFromTmSignatures(blk.Block.LastCommit.Signatures)
	}

	res.rb = blk

	return res
}

func NewValidatorPreCommitsFromTmSignatures(sigs []tmtypes.CommitSig) []ValidatorPreCommit {
	res := make([]ValidatorPreCommit, 0, len(sigs))
	for _, sig := range sigs {
		if len(sig.Signature) == 0 {
			continue
		}

		res = append(res, ValidatorPreCommit{
			ValidatorAddress: sdk.ConsAddress(sig.ValidatorAddress).String(),
			BlockIDFlag:      uint64(sig.BlockIDFlag),
			Timestamp:        sig.Timestamp,
		})
	}

	return res
}

func NewTxsFromTmTxs(txs []*sdktx.GetTxResponse, cdc codec.Codec) Txs {
	res := make(Txs, len(txs))
	for i, tx := range txs {
		var signer string
		if tx.Tx.AuthInfo != nil {
			if len(tx.Tx.AuthInfo.SignerInfos) > 0 && tx.Tx.AuthInfo.SignerInfos[0].PublicKey != nil {
				var pk cryptotypes.PubKey
				if err := cdc.UnpackAny(tx.Tx.AuthInfo.SignerInfos[0].PublicKey, &pk); err == nil {
					signer, _ = ConvertAddressToBech32String(pk.Address())
				}
			}
		}

		res[i] = &Tx{
			Tx:         tx.Tx,
			TxResponse: tx.TxResponse,
			Signer:     signer,
		}
	}

	return res
}

func NewValidatorsFromTmValidator(tmVals *coretypes.ResultValidators) Validators {
	res := make(Validators, 0, len(tmVals.Validators))
	for _, val := range tmVals.Validators {
		consAddr := sdk.ConsAddress(val.Address).String()
		if consPubKey, err := ConvertValidatorPubKeyToBech32String(val.PubKey); err == nil {
			// we need only valid validators
			res = append(res, NewValidator(consAddr, consPubKey))
		}
	}

	return res
}

// NewValidator allows to build a new Validator instance
func NewValidator(consAddr string, consPubKey string) *Validator {
	return &Validator{
		ConsAddr:   consAddr,
		ConsPubkey: consPubKey,
	}
}

func ConvertAddressToBech32String(address cryptotypes.Address) (string, error) {
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	return bech32.ConvertAndEncode(bech32Prefix, address)
}

// ConvertValidatorPubKeyToBech32String converts the given pubKey to Bech32 string
func ConvertValidatorPubKeyToBech32String(pubKey crypto.PubKey) (string, error) {
	bech32Prefix := sdk.GetConfig().GetBech32ConsensusPubPrefix()
	return bech32.ConvertAndEncode(bech32Prefix, pubKey.Bytes())
}

// TotalGas calculates and returns total used gas of all transactions
func (txs Txs) TotalGas() (totalGas uint64) {
	for _, tx := range txs {
		totalGas += uint64(tx.GasUsed) //nolint:gosec
	}

	return totalGas
}

// FindEventByType searches inside the given tx events for the message having the specified index, in order
// to find the event having the given type, and returns it.
// If no such event is found, returns an error instead.
func (tx Tx) FindEventByType(index int, eventType string) (sdk.StringEvent, error) {
	for _, ev := range tx.Logs[index].Events {
		if ev.Type == eventType {
			return ev, nil
		}
	}

	return sdk.StringEvent{}, fmt.Errorf("%w: %s inside tx with hash %s", ErrNoEventFound,
		eventType, tx.TxHash)
}

// FindAttributeByKey searches inside the specified event of the given tx to find the attribute having the given key.
// If the specified event does not contain a such attribute, returns an error instead.
func (tx Tx) FindAttributeByKey(event sdk.StringEvent, attrKey string) (string, error) {
	for _, attr := range event.Attributes {
		if attr.Key == attrKey {
			return attr.Value, nil
		}
	}

	return "", fmt.Errorf("%w: %s found inside tx with hash %s", ErrNoAttributeFound, attrKey, tx.TxHash)
}

// Successful tells whether this tx is successful or not
func (tx Tx) Successful() bool {
	return tx.TxResponse.Code == 0
}

func (b Block) Raw() *coretypes.ResultBlock { return b.rb }
