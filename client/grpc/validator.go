package grpc

import (
	"context"

	"github.com/cometbft/cometbft/crypto/ed25519"
	cometbftcoretypes "github.com/cometbft/cometbft/rpc/core/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/types/query"
)

const (
	defaultLimit = 100
)

func (c *Client) Validators(ctx context.Context, height int64) (*cometbftcoretypes.ResultValidators, error) {
	vals := &cometbftcoretypes.ResultValidators{
		BlockHeight: height,
	}

	var offset uint64

	for {
		respPb, err := c.TmsService.GetValidatorSetByHeight(ctx, &cmtservice.GetValidatorSetByHeightRequest{
			Height: height,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      defaultLimit,
				CountTotal: true,
			},
		})
		if err != nil {
			return nil, err
		}

		if offset == 0 { // first iteration
			vals.Validators = make([]*cometbfttypes.Validator, 0, vals.Total)
		}

		for _, val := range respPb.Validators {
			vals.Validators = append(vals.Validators, convertValidator(val))
		}

		vals.Total = int(respPb.Pagination.Total)

		if len(respPb.Validators) < defaultLimit {
			break
		}

		offset += defaultLimit
	}

	vals.Count = len(vals.Validators)

	return vals, nil
}

func convertValidator(c *cmtservice.Validator) *cometbfttypes.Validator {
	pk := ed25519.PubKey(c.PubKey.Value)

	return &cometbfttypes.Validator{
		Address:          cometbfttypes.Address(c.Address),
		PubKey:           &pk,
		VotingPower:      c.VotingPower,
		ProposerPriority: c.ProposerPriority,
	}
}
