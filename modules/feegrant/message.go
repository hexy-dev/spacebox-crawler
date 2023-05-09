package feegrant

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *feegranttypes.MsgGrantAllowance:

		var (
			allowance  feegranttypes.FeeAllowanceI
			expiration time.Time
		)
		if err := m.cdc.UnpackAny(msg.Allowance, &allowance); err != nil {
			return err
		}

		ex, err := allowance.ExpiresAt()
		if err != nil {
			return err
		}

		if ex != nil {
			expiration = *ex
		}

		data, err := m.cdc.MarshalJSON(msg.Allowance)
		if err != nil {
			return err
		}

		if err = m.broker.PublishGrantAllowanceMessage(ctx, model.GrantAllowanceMessage{
			Height:     tx.Height,
			MsgIndex:   int64(index),
			TxHash:     tx.TxHash,
			Granter:    msg.Granter,
			Grantee:    msg.Grantee,
			Expiration: expiration,
			Allowance:  data,
		}); err != nil {
			m.log.Err(err).Int64("height", tx.Height).Msg("error while publishing grant allowance message")
			return err
		}

		if err = m.publishFeeAllowance(ctx, tx.Height, msg.Granter, msg.Grantee); err != nil {
			m.log.Err(err).
				Int64("height", tx.Height).
				Str("message", "MsgGrantAllowance").
				Msg("error while publishing fee allowance")
			return err
		}

	case *feegranttypes.MsgRevokeAllowance:
		if err := m.broker.PublishRevokeAllowanceMessage(ctx, model.RevokeAllowanceMessage{
			Height:   tx.Height,
			MsgIndex: int64(index),
			TxHash:   tx.TxHash,
			Granter:  msg.Granter,
			Grantee:  msg.Grantee,
		}); err != nil {
			m.log.Err(err).Int64("height", tx.Height).Msg("error while publishing grant allowance message")
			return err
		}

		if err := m.publishFeeAllowance(ctx, tx.Height, msg.Granter, msg.Grantee); err != nil {
			m.log.Err(err).
				Int64("height", tx.Height).
				Str("message", "MsgRevokeAllowance").
				Msg("error while publishing fee allowance")
			return err
		}
	}

	return nil
}

func (m *Module) publishFeeAllowance(ctx context.Context, height int64, granter, grantee string) error {
	respPb, err := m.client.FeegrantQueryClient.Allowance(ctx, &feegranttypes.QueryAllowanceRequest{
		Granter: granter,
		Grantee: grantee,
	})
	if err != nil {
		// set fee allowance to inactive if it was not found
		if err.Error() == "rpc error: code = Internal desc = fee-grant not found: unauthorized" { //nolint:misspell
			m.log.Debug().
				Str("granter", granter).
				Str("grantee", grantee).
				Int64("height", height).
				Msg("fee allowance not found, setting to inactive")

			return m.broker.PublishFeeAllowance(ctx, model.FeeAllowance{
				Granter: granter,
				Grantee: grantee,
				Height:  height,
			})
		}

		return errors.Wrap(err, "error while querying fee allowance")
	}

	allowanceBytes, err := m.cdc.MarshalJSON(respPb.Allowance)
	if err != nil {
		return err
	}

	var (
		allowance  feegranttypes.FeeAllowanceI
		expiration time.Time
	)
	if err = m.cdc.UnpackAny(respPb.Allowance.Allowance, &allowance); err != nil {
		return err
	}

	ex, err := allowance.ExpiresAt()
	if err != nil {
		return err
	}

	if ex != nil {
		expiration = *ex
	}

	return m.broker.PublishFeeAllowance(ctx, model.FeeAllowance{
		Granter:    granter,
		Grantee:    grantee,
		Allowance:  allowanceBytes,
		Expiration: expiration,
		Height:     height,
	})
}
