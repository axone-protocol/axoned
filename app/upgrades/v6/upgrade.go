package v5

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const UpgradeName = "v6.0.0"

var StoreUpgrades = &storetypes.StoreUpgrades{
	Added: []string{
		"feeibc",
	},
}

// CreateUpgradeHandler is the handler that will perform migration from v5.0.0 to v6.0.0.
// Migrate the mint module with new parameters.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		logger.Debug("running module migrations...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
