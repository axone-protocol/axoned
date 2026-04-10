package keeper_test

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v15/x/logic"
	"github.com/axone-protocol/axoned/v15/x/logic/keeper"
	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

func TestQueryProgram(t *testing.T) {
	Convey("Given a keeper with a stored program", t, func() {
		logicKeeper, sdkCtx := newQueryKeeper(t)
		programIDHex := fullProgramIDHex("01")
		programID := mustDecodeHex(t, programIDHex)

		err := logicKeeper.SetStoredProgram(sdkCtx, programID, types.StoredProgram{
			Source:     "foo.",
			CreatedAt:  42,
			SourceSize: 4,
		})
		So(err, ShouldBeNil)

		Convey("When querying the program", func() {
			res, err := logicKeeper.Program(sdkCtx, &types.QueryProgramRequest{ProgramId: programIDHex})

			Convey("Then it should return the program metadata", func() {
				So(err, ShouldBeNil)
				So(res.Program, ShouldResemble, types.ProgramMetadata{
					ProgramId:  programIDHex,
					CreatedAt:  42,
					SourceSize: 4,
				})
			})
		})
	})
}

func TestQueryProgramSource(t *testing.T) {
	Convey("Given a keeper with a stored program", t, func() {
		logicKeeper, sdkCtx := newQueryKeeper(t)
		programIDHex := fullProgramIDHex("02")
		programID := mustDecodeHex(t, programIDHex)

		err := logicKeeper.SetStoredProgram(sdkCtx, programID, types.StoredProgram{
			Source:     "bar.",
			CreatedAt:  84,
			SourceSize: 4,
		})
		So(err, ShouldBeNil)

		Convey("When querying the program source", func() {
			res, err := logicKeeper.ProgramSource(sdkCtx, &types.QueryProgramSourceRequest{ProgramId: programIDHex})

			Convey("Then it should return the source", func() {
				So(err, ShouldBeNil)
				So(res.Source, ShouldEqual, "bar.")
			})
		})
	})
}

func TestQueryPrograms(t *testing.T) {
	Convey("Given a keeper with multiple stored programs", t, func() {
		logicKeeper, sdkCtx := newQueryKeeper(t)
		programID1Hex := fullProgramIDHex("01")
		programID2Hex := fullProgramIDHex("02")

		err := logicKeeper.SetStoredProgram(sdkCtx, mustDecodeHex(t, programID1Hex), types.StoredProgram{
			Source:     "alpha.",
			CreatedAt:  10,
			SourceSize: 6,
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetStoredProgram(sdkCtx, mustDecodeHex(t, programID2Hex), types.StoredProgram{
			Source:     "beta.",
			CreatedAt:  20,
			SourceSize: 5,
		})
		So(err, ShouldBeNil)

		Convey("When querying programs with pagination", func() {
			res, err := logicKeeper.Programs(sdkCtx, &types.QueryProgramsRequest{
				Pagination: &querytypes.PageRequest{Limit: 1, CountTotal: true},
			})

			Convey("Then it should return the first program with pagination info", func() {
				So(err, ShouldBeNil)
				So(len(res.Programs), ShouldEqual, 1)
				So(res.Programs[0], ShouldResemble, types.ProgramMetadata{
					ProgramId:  programID1Hex,
					CreatedAt:  10,
					SourceSize: 6,
				})
				So(res.Pagination, ShouldNotBeNil)
				So(res.Pagination.Total, ShouldEqual, 2)
				So(len(res.Pagination.NextKey), ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestQueryProgramsByPublisher(t *testing.T) {
	Convey("Given a keeper with programs published by different publishers", t, func() {
		logicKeeper, sdkCtx := newQueryKeeper(t)
		publisher := sdk.AccAddress([]byte("publisher-address-01")).String()
		otherPublisher := sdk.AccAddress([]byte("publisher-address-02")).String()
		programID1Hex := fullProgramIDHex("01")
		programID2Hex := fullProgramIDHex("02")

		err := logicKeeper.SetStoredProgram(sdkCtx, mustDecodeHex(t, programID1Hex), types.StoredProgram{
			Source:     "alpha.",
			CreatedAt:  10,
			SourceSize: 6,
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetStoredProgram(sdkCtx, mustDecodeHex(t, programID2Hex), types.StoredProgram{
			Source:     "beta.",
			CreatedAt:  20,
			SourceSize: 5,
		})
		So(err, ShouldBeNil)

		publisherAddr, err := sdk.AccAddressFromBech32(publisher)
		So(err, ShouldBeNil)
		otherPublisherAddr, err := sdk.AccAddressFromBech32(otherPublisher)
		So(err, ShouldBeNil)

		err = logicKeeper.SetProgramPublication(sdkCtx, publisherAddr, mustDecodeHex(t, programID1Hex), types.ProgramPublication{
			PublishedAt: 100,
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetProgramPublication(sdkCtx, publisherAddr, mustDecodeHex(t, programID2Hex), types.ProgramPublication{
			PublishedAt: 200,
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetProgramPublication(sdkCtx, otherPublisherAddr, mustDecodeHex(t, programID2Hex), types.ProgramPublication{
			PublishedAt: 300,
		})
		So(err, ShouldBeNil)

		Convey("When querying programs by publisher", func() {
			res, err := logicKeeper.ProgramsByPublisher(sdkCtx, &types.QueryProgramsByPublisherRequest{
				Publisher:  publisher,
				Pagination: &querytypes.PageRequest{Limit: 10, CountTotal: true},
			})

			Convey("Then it should return only that publisher's programs", func() {
				So(err, ShouldBeNil)
				So(len(res.Programs), ShouldEqual, 2)
				So(res.Programs[0], ShouldResemble, types.PublishedProgram{
					Program: types.ProgramMetadata{
						ProgramId:  programID1Hex,
						CreatedAt:  10,
						SourceSize: 6,
					},
					Publication: types.ProgramPublication{PublishedAt: 100},
				})
				So(res.Programs[1], ShouldResemble, types.PublishedProgram{
					Program: types.ProgramMetadata{
						ProgramId:  programID2Hex,
						CreatedAt:  20,
						SourceSize: 5,
					},
					Publication: types.ProgramPublication{PublishedAt: 200},
				})
				So(res.Pagination, ShouldNotBeNil)
				So(res.Pagination.Total, ShouldEqual, 2)
			})
		})

		Convey("When querying programs by publisher with pagination", func() {
			res, err := logicKeeper.ProgramsByPublisher(sdkCtx, &types.QueryProgramsByPublisherRequest{
				Publisher:  publisher,
				Pagination: &querytypes.PageRequest{Limit: 1, CountTotal: true},
			})

			Convey("Then it should return pagination info for the next page", func() {
				So(err, ShouldBeNil)
				So(len(res.Programs), ShouldEqual, 1)
				So(res.Programs[0], ShouldResemble, types.PublishedProgram{
					Program: types.ProgramMetadata{
						ProgramId:  programID1Hex,
						CreatedAt:  10,
						SourceSize: 6,
					},
					Publication: types.ProgramPublication{PublishedAt: 100},
				})
				So(res.Pagination, ShouldNotBeNil)
				So(res.Pagination.Total, ShouldEqual, 2)
				So(len(res.Pagination.NextKey), ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestQueryProgramErrors(t *testing.T) {
	Convey("Given a query keeper", t, func() {
		logicKeeper, sdkCtx, key := newQueryKeeperWithStoreKey(t)

		Convey("When querying with nil request", func() {
			_, err := logicKeeper.Program(sdkCtx, nil)

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying with invalid program ID", func() {
			_, err := logicKeeper.Program(sdkCtx, &types.QueryProgramRequest{ProgramId: "zz"})

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying non-existent program", func() {
			_, err := logicKeeper.Program(sdkCtx, &types.QueryProgramRequest{ProgramId: fullProgramIDHex("09")})

			Convey("Then it should return NotFound error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.NotFound)
			})
		})

		Convey("When querying with a short program ID", func() {
			_, err := logicKeeper.Program(sdkCtx, &types.QueryProgramRequest{ProgramId: "00"})

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying with an empty program ID", func() {
			_, err := logicKeeper.Program(sdkCtx, &types.QueryProgramRequest{ProgramId: ""})

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying a program with corrupted stored bytes", func() {
			programIDHex := fullProgramIDHex("01")
			sdkCtx.KVStore(key).Set(types.StoredProgramKey(mustDecodeHex(t, programIDHex)), []byte("invalid"))
			_, err := logicKeeper.Program(sdkCtx, &types.QueryProgramRequest{ProgramId: programIDHex})

			Convey("Then it should return Internal error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.Internal)
			})
		})

		Convey("When querying programs by invalid publisher", func() {
			_, err := logicKeeper.ProgramsByPublisher(sdkCtx, &types.QueryProgramsByPublisherRequest{Publisher: "invalid"})

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})
	})
}

func TestQueryProgramSourceErrors(t *testing.T) {
	Convey("Given a query keeper", t, func() {
		logicKeeper, sdkCtx, key := newQueryKeeperWithStoreKey(t)

		Convey("When querying the program source with nil request", func() {
			_, err := logicKeeper.ProgramSource(sdkCtx, nil)

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying the program source with an empty program ID", func() {
			_, err := logicKeeper.ProgramSource(sdkCtx, &types.QueryProgramSourceRequest{ProgramId: ""})

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying the source of a non-existent program", func() {
			_, err := logicKeeper.ProgramSource(sdkCtx, &types.QueryProgramSourceRequest{ProgramId: fullProgramIDHex("0a")})

			Convey("Then it should return NotFound error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.NotFound)
			})
		})

		Convey("When querying the program source with a short program ID", func() {
			_, err := logicKeeper.ProgramSource(sdkCtx, &types.QueryProgramSourceRequest{ProgramId: "00"})

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying a program source with corrupted stored bytes", func() {
			programIDHex := fullProgramIDHex("02")
			sdkCtx.KVStore(key).Set(types.StoredProgramKey(mustDecodeHex(t, programIDHex)), []byte("invalid"))
			_, err := logicKeeper.ProgramSource(sdkCtx, &types.QueryProgramSourceRequest{ProgramId: programIDHex})

			Convey("Then it should return Internal error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.Internal)
			})
		})
	})
}

func TestQueryProgramsErrors(t *testing.T) {
	Convey("Given a query keeper", t, func() {
		logicKeeper, sdkCtx, key := newQueryKeeperWithStoreKey(t)

		Convey("When querying programs with nil request", func() {
			_, err := logicKeeper.Programs(sdkCtx, nil)

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying programs with corrupted stored bytes", func() {
			sdkCtx.KVStore(key).Set(types.StoredProgramKey(mustDecodeHex(t, fullProgramIDHex("03"))), []byte("invalid"))
			_, err := logicKeeper.Programs(sdkCtx, &types.QueryProgramsRequest{
				Pagination: &querytypes.PageRequest{Limit: 10},
			})

			Convey("Then it should return Internal error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.Internal)
			})
		})
	})
}

func TestQueryProgramsByPublisherErrors(t *testing.T) {
	Convey("Given a query keeper", t, func() {
		logicKeeper, sdkCtx, key := newQueryKeeperWithStoreKey(t)
		publisher := sdk.AccAddress([]byte("publisher-address-01")).String()
		publisherAddr, err := sdk.AccAddressFromBech32(publisher)
		So(err, ShouldBeNil)

		Convey("When querying programs by publisher with nil request", func() {
			_, err := logicKeeper.ProgramsByPublisher(sdkCtx, nil)

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When querying programs by publisher with empty publisher", func() {
			_, err := logicKeeper.ProgramsByPublisher(sdkCtx, &types.QueryProgramsByPublisherRequest{})

			Convey("Then it should return InvalidArgument error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("When a publication contains corrupted bytes", func() {
			sdkCtx.KVStore(key).Set(
				types.ProgramPublicationKey(publisherAddr, mustDecodeHex(t, fullProgramIDHex("04"))),
				[]byte("invalid"),
			)
			_, err := logicKeeper.ProgramsByPublisher(sdkCtx, &types.QueryProgramsByPublisherRequest{
				Publisher:  publisher,
				Pagination: &querytypes.PageRequest{Limit: 10},
			})

			Convey("Then it should return Internal error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.Internal)
			})
		})

		Convey("When a publication references a missing stored program", func() {
			err := logicKeeper.SetProgramPublication(sdkCtx, publisherAddr, mustDecodeHex(t, fullProgramIDHex("05")), types.ProgramPublication{
				PublishedAt: 100,
			})
			So(err, ShouldBeNil)
			_, err = logicKeeper.ProgramsByPublisher(sdkCtx, &types.QueryProgramsByPublisherRequest{
				Publisher:  publisher,
				Pagination: &querytypes.PageRequest{Limit: 10},
			})

			Convey("Then it should return Internal error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.Internal)
			})
		})

		Convey("When a referenced stored program contains corrupted bytes", func() {
			programIDHex := fullProgramIDHex("06")
			err := logicKeeper.SetProgramPublication(sdkCtx, publisherAddr, mustDecodeHex(t, programIDHex), types.ProgramPublication{
				PublishedAt: 200,
			})
			So(err, ShouldBeNil)
			sdkCtx.KVStore(key).Set(types.StoredProgramKey(mustDecodeHex(t, programIDHex)), []byte("invalid"))
			_, err = logicKeeper.ProgramsByPublisher(sdkCtx, &types.QueryProgramsByPublisherRequest{
				Publisher:  publisher,
				Pagination: &querytypes.PageRequest{Limit: 10},
			})

			Convey("Then it should return Internal error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.Internal)
			})
		})
	})
}

func newQueryKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	logicKeeper, sdkCtx, _ := newQueryKeeperWithStoreKey(t)

	return logicKeeper, sdkCtx
}

func newQueryKeeperWithStoreKey(t *testing.T) (*keeper.Keeper, sdk.Context, *storetypes.KVStoreKey) {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

	logicKeeper := keeper.NewKeeper(
		encCfg.Codec,
		encCfg.InterfaceRegistry,
		key,
		key,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		nil,
		nil,
		nil,
		nil,
	)

	return logicKeeper, testCtx.Ctx, key
}

func mustDecodeHex(t *testing.T, value string) []byte {
	t.Helper()

	bz, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("failed to decode hex: %v", err)
	}

	return bz
}

func fullProgramIDHex(byteHex string) string {
	return strings.Repeat(byteHex, 32)
}

var _ context.Context = sdk.Context{}
