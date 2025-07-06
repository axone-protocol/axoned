package keeper

import (
	"github.com/jinzhu/copier"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	v1beta2types "github.com/axone-protocol/axoned/v12/x/logic/legacy/v1beta2/types"
	"github.com/axone-protocol/axoned/v12/x/logic/types"
)

func MigrateStoreV3ToV4(k Keeper) module.MigrationHandler {
	getParams := func(ctx sdk.Context) (params v1beta2types.Params, err error) {
		store := ctx.KVStore(k.storeKey)
		bz := store.Get(types.ParamsKey)
		if bz == nil {
			return params, nil
		}
		err = k.cdc.Unmarshal(bz, &params)

		return params, err
	}

	return func(ctx sdk.Context) error {
		paramsFrom, err := getParams(ctx)
		if err != nil {
			return err
		}

		var paramsTo types.Params

		// Interpreter
		if err := copier.Copy(&paramsTo.Interpreter, paramsFrom.Interpreter); err != nil {
			return err
		}

		// Limits
		if v := paramsFrom.Limits.MaxSize; v != nil {
			paramsTo.Limits.MaxSize = v.Uint64()
		}
		if v := paramsFrom.Limits.MaxResultCount; v != nil {
			paramsTo.Limits.MaxResultCount = v.Uint64()
		}
		if v := paramsFrom.Limits.MaxUserOutputSize; v != nil {
			paramsTo.Limits.MaxUserOutputSize = v.Uint64()
		}
		if v := paramsFrom.Limits.MaxVariables; v != nil {
			paramsTo.Limits.MaxVariables = v.Uint64()
		}

		// GasPolicy
		if v := paramsFrom.GasPolicy.WeightingFactor; v != nil {
			paramsTo.GasPolicy.WeightingFactor = v.Uint64()
		}
		if v := paramsFrom.GasPolicy.DefaultPredicateCost; v != nil {
			paramsTo.GasPolicy.DefaultPredicateCost = v.Uint64()
		}
		if v := paramsFrom.GasPolicy.PredicateCosts; v != nil {
			for _, pc := range v {
				paramsTo.GasPolicy.PredicateCosts = append(paramsTo.GasPolicy.PredicateCosts, types.PredicateCost{
					Predicate: pc.Predicate,
					Cost:      pc.Cost.Uint64(),
				})
			}
		}

		return k.SetParams(ctx, paramsTo)
	}
}
