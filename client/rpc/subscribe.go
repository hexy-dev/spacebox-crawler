package rpc

import (
	"context"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
)

func (c *Client) SubscribeNewBlocks(ctx context.Context) (<-chan coretypes.ResultEvent, error) {
	return c.RPCClient.Subscribe(ctx, "", "tm.event = 'NewBlock'")
}
