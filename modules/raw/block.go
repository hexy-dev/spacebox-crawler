package raw

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/bro-n-bro/spacebox-crawler/v2/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	rawBlock := struct {
		Hash            string          `json:"hash"`
		ProposerAddress string          `json:"proposer_address"`
		Block           json.RawMessage `json:"block"`
		TotalGas        uint64          `json:"total_gas"`
		NumTxs          uint16          `json:"num_txs"`
	}{
		TotalGas:        block.TotalGas,
		Hash:            block.Hash,
		ProposerAddress: block.ProposerAddress,
		NumTxs:          uint16(block.TxNum), //nolint:gosec
	}

	var err error
	rawBlock.Block, err = jsoniter.Marshal(block.Raw().Block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}

	if err = m.broker.PublishRawBlock(ctx, rawBlock); err != nil {
		return fmt.Errorf("failed to publish raw block: %w", err)
	}

	return m.publishBlockResults(ctx, block.Height, block.Timestamp)
}

func (m *Module) publishBlockResults(ctx context.Context, height int64, timestamp time.Time) error {
	brResp, err := m.rpcClient.GetBlockResults(ctx, height)
	if err != nil {
		return fmt.Errorf("failed to get block results: %w", err)
	}

	rawBR := struct {
		*coretypes.ResultBlockResults
		Timestamp time.Time `json:"timestamp"`
	}{
		ResultBlockResults: brResp,
		Timestamp:          timestamp,
	}

	return m.broker.PublishRawBlockResults(ctx, rawBR)
}
