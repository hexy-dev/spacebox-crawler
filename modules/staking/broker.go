package staking

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
)

type broker interface {
	PublishAccounts(context.Context, []model.Account) error // FIXME: method from auth module

	PublishUnbondingDelegation(context.Context, model.UnbondingDelegation) error
	PublishUnbondingDelegationMessage(context.Context, model.UnbondingDelegationMessage) error
	PublishStakingParams(ctx context.Context, sp model.StakingParams) error
	PublishDelegation(ctx context.Context, d model.Delegation) error
	PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error
	PublishStakingPool(ctx context.Context, sp model.StakingPool) error
	PublishValidators(ctx context.Context, vals []model.Validator) error
	PublishValidatorsInfo(ctx context.Context, infos []model.ValidatorInfo) error
	PublishValidatorsStatuses(ctx context.Context, statuses []model.ValidatorStatus) error
	PublishRedelegation(context.Context, model.Redelegation) error
	PublishRedelegationMessage(context.Context, model.RedelegationMessage) error
}
