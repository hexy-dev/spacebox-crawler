package raw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/jsonpb"

	"github.com/bro-n-bro/spacebox-crawler/v2/types"
)

var (
	bufPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	marshler = &jsonpb.Marshaler{
		EmitDefaults: false, // Set to false if you don't want to include default values in the JSON output
	}
)

func (m *Module) HandleTx(ctx context.Context, tx *types.Tx) error {
	rawTx := struct {
		Signer     string          `json:"signer"`
		TxResponse json.RawMessage `json:"tx_response"`
	}{
		Signer: tx.Signer,
	}

	b := bufPool.Get().(*bytes.Buffer) //nolint:forcetypeassert
	b.Reset()
	defer bufPool.Put(b)

	if err := marshler.Marshal(b, tx.TxResponse); err != nil {
		return fmt.Errorf("failed to marshal tx response: %w", err)
	}

	rawTx.TxResponse = append(rawTx.TxResponse, b.Bytes()...)

	return m.broker.PublishRawTransaction(ctx, rawTx)
}
