package rpc

import (
	"context"
	"encoding/base64"

	"github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/types"
	"golang.org/x/sync/errgroup"
)

//nolint:gosec
func (c *Client) Genesis(ctx context.Context) (*types.GenesisDoc, error) {
	chunk, err := c.RPCClient.GenesisChunked(ctx, 0)
	if err != nil {
		return nil, err
	}

	decodedData, err := base64.StdEncoding.DecodeString(chunk.Data)
	if err != nil {
		return nil, err
	}

	chunks := make([][]byte, chunk.TotalChunks)
	chunks[0] = decodedData

	g, ctx2 := errgroup.WithContext(ctx)
	for i := uint(1); i < uint(chunk.TotalChunks); i++ {
		func(index uint) {
			g.Go(func() error {
				var ch []byte
				ch, err = c.getGenesisChunk(ctx2, index)
				if err != nil {
					return err
				}

				chunks[index] = ch

				return nil
			})
		}(i)
	}

	if err = g.Wait(); err != nil {
		return nil, err
	}

	var totalData []byte
	for _, ch := range chunks {
		totalData = append(totalData, ch...)
	}

	var resp *types.GenesisDoc
	if err = json.Unmarshal(totalData, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) getGenesisChunk(ctx context.Context, id uint) ([]byte, error) {
	resp, err := c.RPCClient.GenesisChunked(ctx, id)
	if err != nil {
		return nil, err
	}

	decodedData, err := base64.StdEncoding.DecodeString(resp.Data)
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}
