package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

	logicfs "github.com/axone-protocol/axoned/v14/x/logic/fs"
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
				So(err, ShouldBeNil)
				if err != nil {
					return
				}

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

func TestConsumeRequestIOGas(t *testing.T) {
	Convey("Given a user source request", t, func() {
		request := &types.QueryAskRequest{
			Program: "foo. bar.",
			Query:   "foo.",
		}

		Convey("When consuming request I/O gas with a zero coefficient", func() {
			gasMeter := storetypes.NewGasMeter(100)
			consumeRequestIOGas(gasMeter, request, 0)

			Convey("Then the I/O coefficient should default to one", func() {
				So(gasMeter.GasConsumed(), ShouldEqual, 13)
			})
		})

		Convey("When consuming request I/O gas for an empty request", func() {
			gasMeter := storetypes.NewGasMeter(100)
			consumeRequestIOGas(gasMeter, &types.QueryAskRequest{}, 1)

			Convey("Then no gas should be consumed", func() {
				So(gasMeter.GasConsumed(), ShouldEqual, 0)
			})
		})

		Convey("When consuming request I/O gas beyond the available limit", func() {
			gasMeter := storetypes.NewGasMeter(1010)
			var recovered any

			func() {
				defer func() {
					recovered = recover()
				}()
				consumeRequestIOGas(gasMeter, request, 100)
			}()

			Convey("Then it should fail as I/O gas consumption", func() {
				So(recovered, ShouldNotBeNil)
				gasErr, ok := recovered.(storetypes.ErrorOutOfGas)
				So(ok, ShouldBeTrue)
				So(gasErr.Descriptor, ShouldEqual, "IO")
			})
		})

		Convey("When request I/O gas multiplication overflows uint64", func() {
			gasMeter := storetypes.NewInfiniteGasMeter()
			overflowRequest := &types.QueryAskRequest{Program: "ab"}

			consumeRequestIOGas(gasMeter, overflowRequest, math.MaxUint64/2+1)

			Convey("Then the charge should saturate to MaxUint64", func() {
				So(gasMeter.GasConsumed(), ShouldEqual, uint64(math.MaxUint64))
			})
		})
	})
}

func TestNewInterpreterConsultsPublishedUserProgram(t *testing.T) {
	Convey("Given a keeper with a published user program", t, func() {
		encCfg := moduletestutil.MakeTestEncodingConfig()
		key := storetypes.NewKVStoreKey(types.StoreKey)
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

		var logicKeeper *Keeper
		logicKeeper = NewKeeper(
			encCfg.Codec,
			encCfg.InterfaceRegistry,
			key,
			key,
			authtypes.NewModuleAddress(govtypes.ModuleName),
			nil,
			nil,
			nil,
			func(ctx context.Context) (fs.FS, error) {
				return logicfs.NewVFS(ctx, nil, logicKeeper)
			},
		)

		source := "published_fact(alice)."
		programID := sha256.Sum256([]byte(source))
		publisher := authtypes.NewModuleAddress("publisher-a")
		err := logicKeeper.SetStoredProgram(testCtx.Ctx, programID[:], types.StoredProgram{
			Source:     source,
			CreatedAt:  10,
			SourceSize: uint64(len(source)),
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetProgramPublication(testCtx.Ctx, publisher, programID[:], types.ProgramPublication{
			PublishedAt: 11,
		})
		So(err, ShouldBeNil)

		Convey("when consulting it from the user-scoped path", func() {
			interpreter, _, err := logicKeeper.newInterpreter(testCtx.Ctx, types.DefaultParams())
			So(err, ShouldBeNil)

			path := "/v1/var/lib/logic/users/" + publisher.String() + "/programs/" + hex.EncodeToString(programID[:]) + ".pl"
			err = interpreter.ExecContext(testCtx.Ctx, ":- consult('"+path+"').")
			answer, queryErr := util.QueryInterpreter(testCtx.Ctx, interpreter, "published_fact(Who).", 1)

			Convey("then the consulted predicates should be available", func() {
				So(err, ShouldBeNil)
				So(queryErr, ShouldBeNil)
				So(answer.Results, ShouldHaveLength, 1)
				So(answer.Results[0].Substitutions, ShouldHaveLength, 1)
				So(answer.Results[0].Substitutions[0].Variable, ShouldEqual, "Who")
				So(answer.Results[0].Substitutions[0].Expression, ShouldEqual, "alice")
			})
		})
	})
}
