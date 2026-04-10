package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a Migrator instance for the state migration.
func NewMigrator(k Keeper) Migrator {
	return Migrator{keeper: k}
}

// Migrate4to5 rewrites the stored params using the v5 canonical schema.
// Removed interpreter fields and the legacy predicate-based gas policy are
// ignored on read and dropped from the rewritten state.
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	store := ctx.KVStore(m.keeper.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return nil
	}

	var params types.Params
	if err := m.keeper.cdc.Unmarshal(bz, &params); err != nil {
		return err
	}
	params.GasPolicy = types.CanonicalGasPolicy(params.GasPolicy)

	return m.keeper.SetParams(ctx, params)
}
