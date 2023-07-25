package liquidity

import (
	"context"
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox-crawler/types/liquidity"
	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	errWrongDenomsLength   = errors.New("wrong denoms length")
	errWrongBalancesLength = errors.New("wrong balances length")
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	var poolID uint64
	switch msg := cosmosMsg.(type) {
	case *liquidity.MsgCreatePool:
		event, err := tx.FindEventByType(index, liquidity.TypeMsgCreatePool)
		if err != nil {
			return err
		}
		poolIDStr, err := tx.FindAttributeByKey(event, liquidity.AttributeValuePoolID)
		if err != nil {
			return err
		}

		poolID, err = strconv.ParseUint(poolIDStr, 10, 64)
		if err != nil {
			return err
		}

	case *liquidity.MsgDepositWithinBatch:
		poolID = msg.PoolId
	case *liquidity.MsgWithdrawWithinBatch:
		poolID = msg.PoolId
	case *liquidity.MsgSwapWithinBatch:
		poolID = msg.PoolId
	}

	return m.updateLiquidityPool(ctx, poolID)
}

func (m *Module) updateLiquidityPool(ctx context.Context, poolID uint64) error {
	resp, err := m.client.LiquidityQueryClient.LiquidityPool(ctx, &liquidity.QueryLiquidityPoolRequest{PoolId: poolID})
	if err != nil {
		return err
	}

	coinA, coinB, err := m.parseReverseCoins(ctx, resp.Pool.ReserveAccountAddress, resp.Pool.ReserveCoinDenoms)
	if err != nil {
		return err
	}

	return m.broker.PublishLiquidityPool(ctx, model.LiquidityPool{
		PoolID:                poolID,
		ReserveAccountAddress: resp.Pool.ReserveAccountAddress,
		ADenom:                coinA.Denom,
		BDenom:                coinB.Denom,
		PoolCoinDenom:         resp.Pool.PoolCoinDenom,
		LiquidityA:            coinA,
		LiquidityB:            coinB,
	})
}

func (m *Module) parseReverseCoins(
	ctx context.Context,
	address string,
	denoms []string,
) (coinA model.Coin, coinB model.Coin, err error) {

	if len(denoms) != 2 {
		err = errWrongDenomsLength
		return
	}

	resp, err := m.client.BankQueryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
		Address: address,
	})
	if err != nil {
		return
	}

	if len(resp.Balances) != 2 {
		err = errWrongBalancesLength
		return
	}

	coins := m.tbM.MapCoins(types.NewCoinsFromCdk(resp.Balances))
	coinA = coins[0]
	coinB = coins[1]

	return
}
