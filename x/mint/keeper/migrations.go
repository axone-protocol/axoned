package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v3 "github.com/okp4/okp4d/x/mint/migrations/v3"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, ctx.KVStore(m.keeper.storeKey), m.keeper.cdc)
}
