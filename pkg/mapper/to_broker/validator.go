package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapValidator(val *types.Validator) model.Validator {
	return model.Validator{
		ConsensusAddress: val.ConsAddr,
		ConsensusPubKey:  val.ConsPubKey,
	}
}

func (tb ToBroker) MapValidators(vals []*types.Validator) []model.Validator {
	res := make([]model.Validator, len(vals))
	for i, val := range vals {
		res[i] = tb.MapValidator(val)
	}
	return res
}
