package app

import (
	"fmt"

	v7 "github.com/axone/axoned/v7/app/upgrades/v7"
)

// RegisterUpgradeHandlers registers the chain upgrade handlers.
func (app *App) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v7.UpgradeName,
		v7.CreateUpgradeHandler(app.ModuleManager, app.configurator),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
}
