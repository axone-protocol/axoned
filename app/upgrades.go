package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	v4 "github.com/okp4/okp4d/app/upgrades/v4"
	v41 "github.com/okp4/okp4d/app/upgrades/v41"
	v5 "github.com/okp4/okp4d/app/upgrades/v5"
	v6 "github.com/okp4/okp4d/app/upgrades/v6"
)

func (app *App) setupUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v4.UpgradeName,
		v4.CreateUpgradeHandler(app.mm, app.configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v41.UpgradeName,
		v41.CreateUpgradeHandler(app.mm, app.configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v5.UpgradeName,
		v5.CreateUpgradeHandler(app.ParamsKeeper, &app.ConsensusParamsKeeper, app.mm, app.configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v6.UpgradeName,
		v6.CreateUpgradeHandler(app.mm, app.configurator),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades
	switch upgradeInfo.Name {
	case v4.UpgradeName:
		storeUpgrades = v4.StoreUpgrades
	case v41.UpgradeName:
		storeUpgrades = v41.StoreUpgrades
	case v5.UpgradeName:
		storeUpgrades = v5.StoreUpgrades
	case v6.UpgradeName:
		storeUpgrades = v6.StoreUpgrades
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
