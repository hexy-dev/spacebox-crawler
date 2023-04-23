package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishFeeAllowance(ctx context.Context, feeAllowance model.FeeAllowance) error {
	data, err := jsoniter.Marshal(feeAllowance)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(FeeAllowance, data)
}
