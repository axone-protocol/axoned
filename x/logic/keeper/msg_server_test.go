package keeper_test

import (
	gocontext "context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/cosmos/gogoproto/proto"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v14/x/logic"
	logicfs "github.com/axone-protocol/axoned/v14/x/logic/fs"
	"github.com/axone-protocol/axoned/v14/x/logic/keeper"
	logictestutil "github.com/axone-protocol/axoned/v14/x/logic/testutil"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

type failingMarshalCodec struct {
	codec.BinaryCodec
}

func (f failingMarshalCodec) Marshal(msg proto.Message) ([]byte, error) {
	if _, ok := msg.(*types.StoredProgram); ok {
		return nil, errors.New("marshal failed")
	}

	return f.BinaryCodec.Marshal(msg)
}

func TestUpdateParams(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			name      string
			request   *types.MsgUpdateParams
			expectErr bool
		}{
			{
				name: "set invalid authority",
				request: &types.MsgUpdateParams{
					Authority: "foo",
				},
				expectErr: true,
			},
			{
				name: "set full valid params",
				request: &types.MsgUpdateParams{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    types.DefaultParams(),
				},
				expectErr: false,
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d: %v, with request: %v", nc, tc.name, tc.request), func() {
					encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
					key := storetypes.NewKVStoreKey(types.StoreKey)
					testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

					// gomock initializations
					ctrl := gomock.NewController(t)
					accountKeeper := logictestutil.NewMockAccountKeeper(ctrl)
					authQueryService := logictestutil.NewMockAuthQueryService(ctrl)
					bankKeeper := logictestutil.NewMockBankKeeper(ctrl)
					fsProvider := logictestutil.NewMockFS(ctrl)

					logicKeeper := keeper.NewKeeper(
						encCfg.Codec,
						encCfg.InterfaceRegistry,
						key,
						key,
						authtypes.NewModuleAddress(govtypes.ModuleName),
						accountKeeper,
						authQueryService,
						bankKeeper,
						func(_ gocontext.Context) (fs.FS, error) {
							return fsProvider, nil
						})

					msgServer := keeper.NewMsgServerImpl(*logicKeeper)

					Convey("when call msg server to update params", func() {
						res, err := msgServer.UpdateParams(testCtx.Ctx, tc.request)

						Convey("then it should return the expected result", func() {
							if tc.expectErr {
								So(err, ShouldNotBeNil)
								So(res, ShouldBeNil)
							} else {
								So(err, ShouldBeNil)
								So(res, ShouldNotBeNil)
							}
						})
					})
				})
		}
	})
}

func TestStoreProgram(t *testing.T) {
	Convey("Given a msg server", t, func() {
		encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
		key := storetypes.NewKVStoreKey(types.StoreKey)
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
		testCtx.Ctx = testCtx.Ctx.WithBlockTime(time.Date(2026, time.March, 9, 12, 0, 0, 0, time.UTC))

		// gomock initializations
		ctrl := gomock.NewController(t)
		accountKeeper := logictestutil.NewMockAccountKeeper(ctrl)
		authQueryService := logictestutil.NewMockAuthQueryService(ctrl)
		bankKeeper := logictestutil.NewMockBankKeeper(ctrl)
		fsProvider := logictestutil.NewMockFS(ctrl)

		logicKeeper := keeper.NewKeeper(
			encCfg.Codec,
			encCfg.InterfaceRegistry,
			key,
			key,
			authtypes.NewModuleAddress(govtypes.ModuleName),
			accountKeeper,
			authQueryService,
			bankKeeper,
			func(_ gocontext.Context) (fs.FS, error) {
				return fsProvider, nil
			},
		)
		if err := logicKeeper.SetParams(testCtx.Ctx, types.DefaultParams()); err != nil {
			t.Fatal(err)
		}
		msgServer := keeper.NewMsgServerImpl(*logicKeeper)

		source := "father(bob, alice)."
		sourceHash := sha256.Sum256([]byte(source))
		expectedProgramID := hex.EncodeToString(sourceHash[:])

		publisherA := authtypes.NewModuleAddress("publisher-a").String()
		publisherB := authtypes.NewModuleAddress("publisher-b").String()
		publisherAAddr, err := sdk.AccAddressFromBech32(publisherA)
		if err != nil {
			t.Fatal(err)
		}
		publisherBAddr, err := sdk.AccAddressFromBech32(publisherB)
		if err != nil {
			t.Fatal(err)
		}

		Convey("when request is nil", func() {
			res, err := msgServer.StoreProgram(testCtx.Ctx, nil)

			Convey("then it should fail", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldBeNil)
			})
		})

		Convey("when publisher is invalid", func() {
			res, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
				Publisher: "foo",
				Source:    source,
			})

			Convey("then it should fail", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldBeNil)
			})
		})

		Convey("when a new source exceeds MaxSize", func() {
			params := types.DefaultParams()
			params.Limits.MaxSize = uint64(len(source)) - 1
			err := logicKeeper.SetParams(testCtx.Ctx, params)
			So(err, ShouldBeNil)

			res, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
				Publisher: publisherA,
				Source:    source,
			})

			Convey("then it should fail before storing", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldBeNil)
				So(testCtx.Ctx.KVStore(key).Get(types.StoredProgramKey(sourceHash[:])),
					ShouldBeNil)
			})
		})

		Convey("when the stored artifact bytes are invalid", func() {
			testCtx.Ctx.KVStore(key).Set(types.StoredProgramKey(sourceHash[:]), []byte("invalid"))

			res, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
				Publisher: publisherA,
				Source:    source,
			})

			Convey("then it should fail", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldBeNil)
			})
		})

		Convey("when an existing artifact under the same id has a different source", func() {
			err := logicKeeper.SetStoredProgram(testCtx.Ctx, sourceHash[:], types.StoredProgram{
				Source:     "different_source.",
				CreatedAt:  testCtx.Ctx.BlockTime().Unix(),
				SourceSize: uint64(len("different_source.")),
			})
			So(err, ShouldBeNil)

			res, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
				Publisher: publisherA,
				Source:    source,
			})

			Convey("then it should fail with a hash collision", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldBeNil)
			})
		})

		Convey("when publication bytes are invalid", func() {
			err := logicKeeper.SetStoredProgram(testCtx.Ctx, sourceHash[:], types.StoredProgram{
				Source:     source,
				CreatedAt:  testCtx.Ctx.BlockTime().Unix(),
				SourceSize: uint64(len(source)),
			})
			So(err, ShouldBeNil)
			testCtx.Ctx.KVStore(key).Set(types.ProgramPublicationKey(publisherAAddr, sourceHash[:]), []byte("invalid"))

			res, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
				Publisher: publisherA,
				Source:    source,
			})

			Convey("then it should fail", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldBeNil)
			})
		})

		Convey("when source is invalid", func() {
			res, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
				Publisher: publisherA,
				Source:    "father(bob, alice)",
			})

			Convey("then it should fail", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldBeNil)
			})
		})

		Convey("when source is valid", func() {
			res, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
				Publisher: publisherA,
				Source:    source,
			})

			Convey("then it should store canonical artifact and publisher publication", func() {
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.ProgramId, ShouldEqual, expectedProgramID)

				bz := testCtx.Ctx.KVStore(key).Get(types.StoredProgramKey(sourceHash[:]))
				So(bz, ShouldNotBeNil)

				var stored types.StoredProgram
				err = encCfg.Codec.Unmarshal(bz, &stored)
				So(err, ShouldBeNil)
				So(stored.Source, ShouldEqual, source)
				So(stored.CreatedAt, ShouldEqual, testCtx.Ctx.BlockTime().Unix())
				So(stored.SourceSize, ShouldEqual, uint64(len(source)))

				pubAKey := types.ProgramPublicationKey(publisherAAddr, sourceHash[:])
				pubABz := testCtx.Ctx.KVStore(key).Get(pubAKey)
				So(pubABz, ShouldNotBeNil)

				var pubA types.ProgramPublication
				err = encCfg.Codec.Unmarshal(pubABz, &pubA)
				So(err, ShouldBeNil)
				So(pubA.PublishedAt, ShouldEqual, testCtx.Ctx.BlockTime().Unix())

				vfs, err := logicfs.NewVFS(testCtx.Ctx, nil, logicKeeper)
				So(err, ShouldBeNil)
				content, err := fs.ReadFile(vfs, "/v1/usr/share/logic/"+publisherA+"/"+expectedProgramID+".pl")
				So(err, ShouldBeNil)
				So(string(content), ShouldEqual, source)
			})

			Convey("and storing the same source with another publisher", func() {
				testCtx.Ctx = testCtx.Ctx.WithBlockTime(testCtx.Ctx.BlockTime().Add(1 * time.Minute))
				res2, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
					Publisher: publisherB,
					Source:    source,
				})

				Convey("then it should keep a single artifact and add publisher publication", func() {
					So(err, ShouldBeNil)
					So(res2, ShouldNotBeNil)
					So(res2.ProgramId, ShouldEqual, expectedProgramID)

					bz := testCtx.Ctx.KVStore(key).Get(types.StoredProgramKey(sourceHash[:]))
					So(bz, ShouldNotBeNil)

					var stored types.StoredProgram
					err = encCfg.Codec.Unmarshal(bz, &stored)
					So(err, ShouldBeNil)
					So(stored.Source, ShouldEqual, source)

					pubAKey := types.ProgramPublicationKey(publisherAAddr, sourceHash[:])
					pubABz := testCtx.Ctx.KVStore(key).Get(pubAKey)
					So(pubABz, ShouldNotBeNil)
					var pubA types.ProgramPublication
					err = encCfg.Codec.Unmarshal(pubABz, &pubA)
					So(err, ShouldBeNil)
					So(pubA.PublishedAt, ShouldEqual, testCtx.Ctx.BlockTime().Add(-1*time.Minute).Unix())

					pubBKey := types.ProgramPublicationKey(publisherBAddr, sourceHash[:])
					pubBBz := testCtx.Ctx.KVStore(key).Get(pubBKey)
					So(pubBBz, ShouldNotBeNil)
					var pubB types.ProgramPublication
					err = encCfg.Codec.Unmarshal(pubBBz, &pubB)
					So(err, ShouldBeNil)
					So(pubB.PublishedAt, ShouldEqual, testCtx.Ctx.BlockTime().Unix())
				})

				Convey("and publishing again from the same publisher", func() {
					testCtx.Ctx = testCtx.Ctx.WithBlockTime(testCtx.Ctx.BlockTime().Add(1 * time.Minute))
					res3, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
						Publisher: publisherB,
						Source:    source,
					})

					Convey("then it should be idempotent for publication too", func() {
						So(err, ShouldBeNil)
						So(res3, ShouldNotBeNil)
						So(res3.ProgramId, ShouldEqual, expectedProgramID)

						pubBKey := types.ProgramPublicationKey(publisherBAddr, sourceHash[:])
						pubBBz := testCtx.Ctx.KVStore(key).Get(pubBKey)
						So(pubBBz, ShouldNotBeNil)

						var pubB types.ProgramPublication
						err = encCfg.Codec.Unmarshal(pubBBz, &pubB)
						So(err, ShouldBeNil)
						So(pubB.PublishedAt, ShouldEqual, testCtx.Ctx.BlockTime().Add(-1*time.Minute).Unix())
					})
				})
			})

			Convey("and republishing after tightening limits", func() {
				params := types.DefaultParams()
				params.Limits.MaxSize = uint64(len(source)) - 1
				err := logicKeeper.SetParams(testCtx.Ctx, params)
				So(err, ShouldBeNil)

				testCtx.Ctx = testCtx.Ctx.WithBlockTime(testCtx.Ctx.BlockTime().Add(2 * time.Minute))
				res4, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
					Publisher: publisherA,
					Source:    source,
				})

				Convey("then it should remain idempotent for an already stored artifact", func() {
					So(err, ShouldBeNil)
					So(res4, ShouldNotBeNil)
					So(res4.ProgramId, ShouldEqual, expectedProgramID)
				})
			})
		})
	})
}

func TestStoreProgramMarshalFailure(t *testing.T) {
	Convey("Given a msg server with a codec failing on marshal", t, func() {
		encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
		key := storetypes.NewKVStoreKey(types.StoreKey)
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

		ctrl := gomock.NewController(t)
		accountKeeper := logictestutil.NewMockAccountKeeper(ctrl)
		authQueryService := logictestutil.NewMockAuthQueryService(ctrl)
		bankKeeper := logictestutil.NewMockBankKeeper(ctrl)
		fsProvider := logictestutil.NewMockFS(ctrl)

		logicKeeper := keeper.NewKeeper(
			failingMarshalCodec{BinaryCodec: encCfg.Codec},
			encCfg.InterfaceRegistry,
			key,
			key,
			authtypes.NewModuleAddress(govtypes.ModuleName),
			accountKeeper,
			authQueryService,
			bankKeeper,
			func(_ gocontext.Context) (fs.FS, error) {
				return fsProvider, nil
			},
		)
		err := logicKeeper.SetParams(testCtx.Ctx, types.DefaultParams())
		So(err, ShouldBeNil)

		msgServer := keeper.NewMsgServerImpl(*logicKeeper)
		publisher := authtypes.NewModuleAddress("publisher-a").String()

		Convey("when storing a valid source", func() {
			res, err := msgServer.StoreProgram(testCtx.Ctx, &types.MsgStoreProgram{
				Publisher: publisher,
				Source:    "father(bob, alice).",
			})

			Convey("then it should surface the storage error", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldBeNil)
			})
		})
	})
}
