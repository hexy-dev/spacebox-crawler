package rpc

import (
	"context"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
)

func (c *Client) GetBlockResults(ctx context.Context, height int64) (*coretypes.ResultBlockResults, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	result, err := c.RPCClient.BlockResults(ctx, &height)
	if err != nil {
		return nil, err
	}

	return result, nil
}
