package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/logic/exported"
	v2 "github.com/okp4/okp4d/x/logic/migrations/v2"
	v3 "github.com/okp4/okp4d/x/logic/migrations/v3"
)

type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

func NewMigrator(keeper Keeper, legacySubspace exported.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
	}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.legacySubspace, m.keeper.cdc)
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.legacySubspace)
}
