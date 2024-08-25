package rpc

import (
	"context"

	tmcoretypes "github.com/tendermint/tendermint/rpc/coretypes"
)

func (c *Client) SubscribeNewBlocks(ctx context.Context) (<-chan tmcoretypes.ResultEvent, error) {
	return c.RPCClient.Subscribe(ctx, "", "tm.event = 'NewBlock'")
}
