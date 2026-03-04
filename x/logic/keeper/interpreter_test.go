package keeper

import (
	"context"
	"io/fs"
	"math"
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

func TestConsumeSourceGas(t *testing.T) {
	Convey("Given a user source request", t, func() {
		request := &types.QueryServiceAskRequest{
			Program: "foo. bar.",
			Query:   "foo.",
		}

		Convey("When consuming source gas with a zero coefficient", func() {
			gasMeter := storetypes.NewGasMeter(100)
			consumeSourceGas(gasMeter, request, 0)

			Convey("Then the source coefficient should default to one", func() {
				So(gasMeter.GasConsumed(), ShouldEqual, 13)
			})
		})

		Convey("When consuming source gas for an empty request", func() {
			gasMeter := storetypes.NewGasMeter(100)
			consumeSourceGas(gasMeter, &types.QueryServiceAskRequest{}, 1)

			Convey("Then no gas should be consumed", func() {
				So(gasMeter.GasConsumed(), ShouldEqual, 0)
			})
		})

		Convey("When consuming source gas beyond the available limit", func() {
			gasMeter := storetypes.NewGasMeter(1010)
			var recovered any

			func() {
				defer func() {
					recovered = recover()
				}()
				consumeSourceGas(gasMeter, request, 100)
			}()

			Convey("Then it should fail as source gas consumption", func() {
				So(recovered, ShouldNotBeNil)
				gasErr, ok := recovered.(storetypes.ErrorOutOfGas)
				So(ok, ShouldBeTrue)
				So(gasErr.Descriptor, ShouldEqual, "Source")
			})
		})

		Convey("When source gas multiplication overflows uint64", func() {
			gasMeter := storetypes.NewInfiniteGasMeter()
			overflowRequest := &types.QueryServiceAskRequest{Program: "ab"}

			consumeSourceGas(gasMeter, overflowRequest, math.MaxUint64/2+1)

			Convey("Then the charge should saturate to MaxUint64", func() {
				So(gasMeter.GasConsumed(), ShouldEqual, uint64(math.MaxUint64))
			})
		})
	})
}
