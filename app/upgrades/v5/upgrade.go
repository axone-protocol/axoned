package v5

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/okp4/okp4d/app"
)

const UpgradeName = "v5.0.0"

var StoreUpgrades *storetypes.StoreUpgrades // No store root upgrade.

// CreateUpgradeHandler is the handler that will perform migration from v4.1.0 to v5.0.0.
// This migration include following update that need migration :
//   - Migrate Tendermint consensus parameters from x/params moduel to a dedicated
//     x/consensus module.
//     -
func CreateUpgradeHandler(
	app *app.App,
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {

	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())

	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		logger.Debug("migrate consensus params keeper")
		baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)

		logger.Debug("running module migrations...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
