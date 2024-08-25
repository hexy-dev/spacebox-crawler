package types

import (
	"context"
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"
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

	GenesisHandler interface {
		Module
		// HandleGenesis handles a genesis state.
		HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error
	}
)
