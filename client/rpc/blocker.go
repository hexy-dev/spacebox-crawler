package rpc

import (
	"context"

	"github.com/bro-n-bro/spacebox-crawler/v2/types"
)

// GetBlockEvents returns begin block and end block events.
func (c *Client) GetBlockEvents(ctx context.Context, height int64) (begin, end types.BlockerEvents, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	result, err := c.RPCClient.BlockResults(ctx, &height)
	if err != nil {
		return nil, nil, err
	}

	begin = types.NewBlockerEventsAttributes(result.FinalizeBlockEvents)
	end = types.NewBlockerEventsAttributes(result.FinalizeBlockEvents)

	return
}
