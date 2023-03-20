package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v2types "github.com/okp4/okp4d/x/logic/migrations/v2/types"
	"github.com/okp4/okp4d/x/logic/types"
)

func MigrateStore(ctx sdk.Context, paramstore paramtypes.Subspace, cdc codec.BinaryCodec) error {
	logger := ctx.Logger().
		With("module", "logic").
		With("migration", "v3")

	logger.Debug("starting module migration")

	logger.Debug("migrate logic params")
	var oldParams v2types.Params
	paramstore.GetParamSet(ctx, &oldParams)

	newParams := types.Params{
		Interpreter: types.Interpreter{
			Bootstrap:           oldParams.Interpreter.Bootstrap,
			PredicatesWhitelist: oldParams.Interpreter.RegisteredPredicates,
			PredicatesBlacklist: []string{},
		},
		Limits: types.Limits{
			MaxGas:         oldParams.Limits.MaxGas,
			MaxSize:        oldParams.Limits.MaxSize,
			MaxResultCount: oldParams.Limits.MaxResultCount,
		},
	}
	paramstore.SetParamSet(ctx, &newParams)

	logger.Debug("module migration done")

	return nil
}
