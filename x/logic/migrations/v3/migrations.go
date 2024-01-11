package v2

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/okp4/okp4d/x/logic/exported"
	v2types "github.com/okp4/okp4d/x/logic/migrations/v2/types"
	"github.com/okp4/okp4d/x/logic/types"
)

func v2ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&v2types.Params{})
}

// MigrateStore migrates the x/logic module state from the consensus version 2 to
// version 3.
// Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/logic
// module state.
// Then, `RegisteredPredicates` has been renamed to `PredicatesWhitelist` and add new
// `PredicatesBlacklist`.
func MigrateStore(ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	legacySubspace exported.Subspace,
) error {
	logger := ctx.Logger().
		With("module", "logic").
		With("migration", "v3")

	logger.Debug("starting module migration")

	logger.Debug("migrate logic params")

	store := ctx.KVStore(storeKey)

	var oldParams v2types.Params
	legacySubspace.WithKeyTable(v2ParamKeyTable()).
		GetParamSet(ctx, &oldParams)

	newParams := types.Params{
		Interpreter: types.Interpreter{
			Bootstrap: oldParams.Interpreter.Bootstrap,
			PredicatesFilter: types.Filter{
				Whitelist: oldParams.Interpreter.RegisteredPredicates,
				Blacklist: []string{},
			},
		},
		Limits: types.Limits{
			MaxGas:         oldParams.Limits.MaxGas,
			MaxSize:        oldParams.Limits.MaxSize,
			MaxResultCount: oldParams.Limits.MaxResultCount,
		},
	}

	if err := newParams.Validate(); err != nil {
		return err
	}

	bz, err := cdc.Marshal(&newParams)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	logger.Debug("module migration done")

	return nil
}
