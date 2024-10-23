package app

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
)

var upgrades = []string{
	"v11.0.0",
}

// registerUpgradeHandlers registers the chain upgrade handlers.
func (app *App) registerUpgradeHandlers() {
	for _, upgrade := range upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade,
			func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
				return app.ModuleManager.RunMigrations(ctx, app.configurator, vm)
			},
		)
	}

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
}
