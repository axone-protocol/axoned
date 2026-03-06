package mint_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/axone-protocol/axoned/v14/x/mint/testutil"
	"github.com/axone-protocol/axoned/v14/x/mint/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	Convey("Given the module initialization", t, func() {
		var accountKeeper authkeeper.AccountKeeper

		app, err := simtestutil.SetupAtGenesis(
			depinject.Configs(
				testutil.AppConfig,
				depinject.Supply(log.NewNopLogger()),
			), &accountKeeper)
		So(err, ShouldBeNil)

		ctx := app.NewContext(false)
		acc := accountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
		So(acc, ShouldNotBeNil)
	})
}
