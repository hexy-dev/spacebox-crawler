package rep

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/tx"
	tmcoretypes "github.com/tendermint/tendermint/rpc/coretypes"
	tmtypes "github.com/tendermint/tendermint/types"
)

type (
	GrpcClient interface {
		Block(ctx context.Context, height int64) (*tmcoretypes.ResultBlock, error)

		Txs(ctx context.Context, height int64, txs tmtypes.Txs) ([]*tx.GetTxResponse, error)
	}

	RPCClient interface {
		SubscribeNewBlocks(ctx context.Context) (<-chan tmcoretypes.ResultEvent, error)
		Genesis(ctx context.Context) (*tmtypes.GenesisDoc, error)
		GetLastBlockHeight(ctx context.Context) (int64, error)
	}
)
