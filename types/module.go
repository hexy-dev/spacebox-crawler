package types

import (
	"context"
	"encoding/json"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Module interface {
		// Name base implementation of Module interface.
		Name() string
	}

	BlockHandler interface {
		Module
		// HandleBlock handles a single block in blockchain.
		HandleBlock(ctx context.Context, block *Block) error
	}

	TransactionHandler interface {
		Module
		// HandleTx handles a single transaction of block.
		HandleTx(ctx context.Context, tx *Tx) error
	}

	MessageHandler interface {
		Module
		// HandleMessage handles a single message of transaction.
		HandleMessage(ctx context.Context, index int, msg sdk.Msg, tx *Tx) error
	}

	RecursiveMessagesHandler interface {
		Module
		HandleMessageRecursive(ctx context.Context, index int, msg sdk.Msg, tx *Tx) ([]*types.Any, error)
	}

	ValidatorsHandler interface {
		Module
		// HandleValidators handles of all validators in blockchain.
		HandleValidators(ctx context.Context, vals *coretypes.ResultValidators) error
	}

	GenesisHandler interface {
		Module
		// HandleGenesis handles a genesis state.
		HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error
	}

	BeginBlockerHandler interface {
		Module
		// HandleBeginBlocker handles of beginblocker events.
		HandleBeginBlocker(ctx context.Context, eventsMap BlockerEvents, height int64) error
	}

	EndBlockerHandler interface {
		Module
		// HandleEndBlocker handles of endblocker events.
		HandleEndBlocker(ctx context.Context, eventsMap BlockerEvents, height int64) error
	}
)
