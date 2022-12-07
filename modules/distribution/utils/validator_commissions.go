package utils

import (
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
	"context"
	"fmt"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// UpdateValidatorsCommissionAmounts updates the validators commissions amounts
func UpdateValidatorsCommissionAmounts(height int64, client distrtypes.QueryClient) {
	//log.Debug().Str("module", "distribution").
	//	Int64("height", height).
	//	Msg("updating validators commissions")

	//validators, err := db.GetValidators()
	//if err != nil {
	//log.Error().Str("module", "distribution").Err(err).
	//	Int64("height", height).
	//	Msg("error while getting validators")
	//return
	//}

	var validators []types.StakingValidator

	if len(validators) == 0 {
		// No validators, just skip
		return
	}

	// Get all the commissions
	for _, validator := range validators {
		go updateValidatorCommission(height, client, validator)
	}
}

func updateValidatorCommission(height int64, distrClient distrtypes.QueryClient, validator types.StakingValidator) {
	commission, err := GetValidatorCommissionAmount(height, validator.GetOperator(), validator.GetSelfDelegateAddress(),
		distrClient)
	if err != nil {
		//log.Error().Str("module", "distribution").Err(err).
		//	Int64("height", height).
		//	Str("validator", validator.GetOperator()).
		//	Send()
	}

	// TODO:
	_ = commission
	//err = db.SaveValidatorCommissionAmount(commission)
	//if err != nil {
	//	log.Error().Str("module", "distribution").Err(err).
	//		Int64("height", height).
	//		Str("validator", validator.GetOperator()).
	//		Msg("error while saving validator commission amounts")
	//}
}

// GetValidatorCommissionAmount returns the amount of the validator commission for the given validator
func GetValidatorCommissionAmount(
	height int64, operatorAddress, selfDelegateAddress string, distrClient distrtypes.QueryClient,
) (types.ValidatorCommissionAmount, error) {
	res, err := distrClient.ValidatorCommission(
		context.Background(),
		&distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: operatorAddress},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return types.ValidatorCommissionAmount{}, fmt.Errorf("error while getting validator commission: %s", err)
	}

	return types.NewValidatorCommissionAmount(
		operatorAddress,
		selfDelegateAddress,
		res.Commission.Commission,
		height,
	), nil
}
