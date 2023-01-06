package bank

import (
	"context"
	"encoding/json"

	"github.com/hexy-dev/spacebox/broker/model"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/hexy-dev/spacebox-crawler/modules/utils"
	"github.com/hexy-dev/spacebox-crawler/types"
)

// HandleGenesis handles the genesis state of the x/bank module in order to store the initial values
// of the different account balances.
func (m *Module) HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var bankState banktypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[banktypes.ModuleName], &bankState); err != nil {
		return err
	}

	// Store the balances
	accounts, err := utils.GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return err
	}

	uniqueAddresses := make(map[string]struct{})
	for _, acc := range accounts {
		uniqueAddresses[acc.Address] = struct{}{}
	}

	for _, balance := range bankState.Balances {
		_, ok := uniqueAddresses[balance.Address]
		if !ok {
			// skip already published
			continue
		}

		ab := model.NewAccountBalance(balance.Address, doc.InitialHeight, m.tbM.MapCoins(types.NewCoinsFromCdk(balance.Coins)))
		// TODO: test it
		if err = m.broker.PublishAccountBalance(ctx, ab); err != nil {
			return err
		}
	}

	return nil
}
