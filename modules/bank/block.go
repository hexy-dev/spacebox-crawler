package bank

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	resp, err := m.client.BankQueryClient.TotalSupply(
		ctx,
		&banktypes.QueryTotalSupplyRequest{},
		grpcClient.GetHeightRequestHeader(block.Height))
	if err != nil {
		return err
	}

	// TODO: test it
	err = m.broker.PublishSupply(ctx, model.NewSupply(block.Height, m.tbM.MapCoins(types.NewCoinsFromCdk(resp.Supply))))
	if err != nil {
		return err
	}
	return nil
}
