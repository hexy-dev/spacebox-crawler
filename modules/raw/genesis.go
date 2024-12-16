package raw

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cometbft/cometbft/types"
	jsoniter "github.com/json-iterator/go"
)

func (m *Module) HandleGenesis(ctx context.Context, doc *types.GenesisDoc, _ map[string]json.RawMessage) error {
	rawGenesis := struct {
		GenesisTime     string          `json:"genesis_time"`
		ChainID         string          `json:"chain_id"`
		AppHash         string          `json:"app_hash"`
		ConsensusParams json.RawMessage `json:"consensus_params"`
		AppState        json.RawMessage `json:"app_state"`
		InitialHeight   int64           `json:"initial_height"`
	}{
		GenesisTime:   doc.GenesisTime.String(),
		ChainID:       doc.ChainID,
		InitialHeight: doc.InitialHeight,
		AppHash:       doc.AppHash.String(),
		AppState:      doc.AppState,
	}

	var err error
	rawGenesis.ConsensusParams, err = jsoniter.Marshal(doc.ConsensusParams)
	if err != nil {
		return fmt.Errorf("failed to marshal consensus params: %w", err)
	}

	return m.broker.PublishRawGenesis(ctx, rawGenesis)
}
