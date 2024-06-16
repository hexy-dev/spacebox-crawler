package grpc

import (
	"context"

	cometbftcoretypes "github.com/cometbft/cometbft/rpc/core/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
)

func (c *Client) Block(ctx context.Context, height int64) (*cometbftcoretypes.ResultBlock, error) {
	resp, err := c.TmsService.GetBlockByHeight(
		ctx,
		&cmtservice.GetBlockByHeightRequest{
			Height: height,
		},
	)
	if err != nil {
		return nil, err
	}

	block, err := cometbfttypes.BlockFromProto(resp.Block) // nolint:staticcheck
	if err != nil {
		return nil, err
	}

	blockID, err := cometbfttypes.BlockIDFromProto(resp.BlockId)
	if err != nil {
		return nil, err
	}

	return &cometbftcoretypes.ResultBlock{
		Block:   block,
		BlockID: *blockID,
	}, nil
}
