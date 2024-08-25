package grpc

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/cosmos/cosmos-sdk/types/tx"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bro-n-bro/spacebox-crawler/v2/adapter/storage/model"
)

// Txs queries for all the transactions in a block. Transactions are returned
// in sdk.TxResponse format which internally contains a sdk.Tx. An error is
// returned if any query fails.
func (c *Client) Txs(ctx context.Context, height int64, txs tmtypes.Txs) ([]*tx.GetTxResponse, error) {
	txResponses := make([]*tx.GetTxResponse, 0, len(txs))

	for _, tmTx := range txs {
		hash := hex.EncodeToString(tmTx.Hash())

		respPb, err := c.TxService.GetTx(ctx, &tx.GetTxRequest{Hash: hash})
		if err != nil {
			_ = c.storage.InsertErrorTx(ctx, model.Tx{
				Created:      time.Now(),
				ErrorMessage: err.Error(),
				Hash:         hash,
				Height:       height,
			})

			c.log.Warn().Err(err).Int64("height", height).Msg("GetTx error")
			continue
		}

		txResponses = append(txResponses, &tx.GetTxResponse{Tx: respPb.Tx, TxResponse: respPb.TxResponse})
	}

	return txResponses, nil
}
