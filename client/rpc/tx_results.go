package rpc

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
)

func (c *Client) GetTxResults(ctx context.Context, height int64) ([]*abci.ExecTxResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	result, err := c.RPCClient.BlockResults(ctx, &height)
	if err != nil {
		return nil, err
	}

	return result.TxsResults, nil
}
