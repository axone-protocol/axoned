package app

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var upgrades = []string{
	"v11.0.0",
	"v13.0.0",
}

// registerUpgradeHandlers registers the chain upgrade handlers.
func (app *App) registerUpgradeHandlers() {
	for _, upgrade := range upgrades {
		switch upgrade {
		case "v13.0.0":
			app.UpgradeKeeper.SetUpgradeHandler(
				upgrade,
				func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
					if err := migrateVestingAccountTypeURLs(ctx, app.AccountKeeper, app.appCodec, app.keys[authtypes.StoreKey]); err != nil {
						return vm, fmt.Errorf("failed to migrate vesting accounts: %w", err)
					}
					return app.ModuleManager.RunMigrations(ctx, app.configurator, vm)
				},
			)
		default:
			app.UpgradeKeeper.SetUpgradeHandler(
				upgrade,
				func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
					return app.ModuleManager.RunMigrations(ctx, app.configurator, vm)
				},
			)
		}
	}

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
}

// migrateVestingAccountTypeURLs migrates vesting account type URLs from old format (/vesting.v1beta1.*)
// to new format (/cosmos.vesting.v1beta1.*) by working directly with raw protobuf data.
//
//nolint:funlen
func migrateVestingAccountTypeURLs(
	ctx context.Context, _ authkeeper.AccountKeeper, cdc codec.Codec, storeKey *storetypes.KVStoreKey,
) error {
	typeURLMappings := map[string]string{
		"/vesting.v1beta1.ContinuousVestingAccount": "/cosmos.vesting.v1beta1.ContinuousVestingAccount",
		"/vesting.v1beta1.DelayedVestingAccount":    "/cosmos.vesting.v1beta1.DelayedVestingAccount",
		"/vesting.v1beta1.PeriodicVestingAccount":   "/cosmos.vesting.v1beta1.PeriodicVestingAccount",
		"/vesting.v1beta1.PermanentLockedAccount":   "/cosmos.vesting.v1beta1.PermanentLockedAccount",
	}
	unsupportedTypeURLs := []string{
		"/vesting.v1beta1.CliffVestingAccount",
	}

	logger := sdk.UnwrapSDKContext(ctx).Logger()
	logger.Info("Starting vesting account type URL migration from /vesting.v1beta1.* to /cosmos.vesting.v1beta1.*")

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(storeKey)
	accountStore := prefix.NewStore(store, authtypes.AddressStoreKeyPrefix)

	migratedCount := 0

	iterator := accountStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		accountBytes := iterator.Value()

		var accountAny types.Any
		if err := cdc.Unmarshal(accountBytes, &accountAny); err != nil {
			return fmt.Errorf("failed to unmarshal account %s bytes: %w", sdk.AccAddress(iterator.Key()), err)
		}

		oldTypeURL := accountAny.TypeUrl

		if lo.Contains(unsupportedTypeURLs, oldTypeURL) {
			return fmt.Errorf("cannot migrate account %s with unsupported type %s - this vesting account type migration is not supported",
				sdk.AccAddress(iterator.Key()), oldTypeURL)
		}

		if newTypeURL, needsMigration := typeURLMappings[oldTypeURL]; needsMigration {
			logger.Debug("Migrating account type URL",
				"old_type", oldTypeURL,
				"new_type", newTypeURL,
				"address", sdk.AccAddress(iterator.Key()))

			accountAny.TypeUrl = newTypeURL

			newAccountBytes, err := cdc.Marshal(&accountAny)
			if err != nil {
				return fmt.Errorf("failed to marshal migrated %s account: %w", sdk.AccAddress(iterator.Key()), err)
			}

			accountStore.Set(iterator.Key(), newAccountBytes)
			migratedCount++

			logger.Debug("Successfully migrated account type URL",
				"address", sdk.AccAddress(iterator.Key()),
				"from", oldTypeURL,
				"to", newTypeURL)
		}
	}

	logger.Info("Vesting account type URL migration completed", "migrated_accounts", migratedCount)
	return nil
}
