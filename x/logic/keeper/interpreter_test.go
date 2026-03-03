package keeper

import (
	"context"
	"io/fs"
	"testing"
	"testing/fstest"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v14/x/logic/types"
	"github.com/axone-protocol/axoned/v14/x/logic/util"
)

func TestNewInterpreterBootstrapIsFree(t *testing.T) {
	Convey("Given a keeper with a finite gas meter", t, func() {
		encCfg := moduletestutil.MakeTestEncodingConfig()
		key := storetypes.NewKVStoreKey(types.StoreKey)
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
		testCtx.Ctx = testCtx.Ctx.WithGasMeter(storetypes.NewGasMeter(100_000))

		logicKeeper := NewKeeper(
			encCfg.Codec,
			encCfg.InterfaceRegistry,
			key,
			key,
			authtypes.NewModuleAddress(govtypes.ModuleName),
			nil,
			nil,
			nil,
			func(context.Context) (fs.FS, error) {
				return fstest.MapFS{}, nil
			},
		)

		Convey("When creating a new interpreter", func() {
			interpreter, _, err := logicKeeper.newInterpreter(testCtx.Ctx, types.DefaultParams())

			Convey("Then the kernel bootstrap should not consume gas", func() {
				So(err, ShouldBeNil)
				So(testCtx.Ctx.GasMeter().GasConsumed(), ShouldEqual, 0)
			})

			Convey("And when executing user-space logic", func() {
				err = interpreter.ExecContext(testCtx.Ctx, "foo.")
				_, queryErr := util.QueryInterpreter(testCtx.Ctx, interpreter, "foo.", 1)

				Convey("Then gas should be consumed", func() {
					So(err, ShouldBeNil)
					So(queryErr, ShouldBeNil)
					So(testCtx.Ctx.GasMeter().GasConsumed(), ShouldBeGreaterThan, 0)
				})
			})
		})
	})
}
