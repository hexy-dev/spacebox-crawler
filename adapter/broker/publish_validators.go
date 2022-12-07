package broker

import (
	"context"

	"github.com/pkg/errors"

	"bro-n-bro-osmosis/adapter/broker/model"

	jsoniter "github.com/json-iterator/go"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (b *Broker) PublishValidators(ctx context.Context, vals []model.Validator) error {
	return nil

	for i := 0; i < len(vals); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		data, err := jsoniter.Marshal(vals[i]) // FIXME: maybe user another way to encode data
		if err != nil {
			return errors.Wrap(err, MsgErrJsonMarshalFail)
		}
		err = b.p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: Validator, Partition: kafka.PartitionAny},
			Value:          data,
			//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
		}, nil)
		if err != nil {
			return errors.Wrap(err, "produce account fail")
		}
	}
	return nil
}
