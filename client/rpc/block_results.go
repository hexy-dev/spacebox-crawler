package rpc

import (
	"context"

	tmcoretypes "github.com/tendermint/tendermint/rpc/coretypes"
)

func (c *Client) GetBlockResults(ctx context.Context, height int64) (*tmcoretypes.ResultBlockResults, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	result, err := c.RPCClient.BlockResults(ctx, &height)
	if err != nil {
		return nil, err
	}

	return result, nil
}
