package keeper_test

import (
	"context"
	"encoding/hex"
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

	"github.com/axone-protocol/axoned/v14/x/logic"
	"github.com/axone-protocol/axoned/v14/x/logic/keeper"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

func TestQueryProgram(t *testing.T) {
	Convey("Given a keeper with a stored program", t, func() {
		logicKeeper, sdkCtx := newQueryKeeper(t)
		programID := mustDecodeHex(t, "01")

		err := logicKeeper.SetStoredProgram(sdkCtx, programID, types.StoredProgram{
			Source:     "foo.",
			CreatedAt:  42,
			SourceSize: 4,
		})
		So(err, ShouldBeNil)

		Convey("When querying the program", func() {
			res, err := logicKeeper.Program(sdkCtx, &types.QueryProgramRequest{ProgramId: "01"})

			Convey("Then it should return the program metadata", func() {
				So(err, ShouldBeNil)
				So(res.Program, ShouldResemble, types.ProgramMetadata{
					ProgramId:  "01",
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
		programID := mustDecodeHex(t, "02")

		err := logicKeeper.SetStoredProgram(sdkCtx, programID, types.StoredProgram{
			Source:     "bar.",
			CreatedAt:  84,
			SourceSize: 4,
		})
		So(err, ShouldBeNil)

		Convey("When querying the program source", func() {
			res, err := logicKeeper.ProgramSource(sdkCtx, &types.QueryProgramSourceRequest{ProgramId: "02"})

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

		err := logicKeeper.SetStoredProgram(sdkCtx, mustDecodeHex(t, "01"), types.StoredProgram{
			Source:     "alpha.",
			CreatedAt:  10,
			SourceSize: 6,
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetStoredProgram(sdkCtx, mustDecodeHex(t, "02"), types.StoredProgram{
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
					ProgramId:  "01",
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

		err := logicKeeper.SetStoredProgram(sdkCtx, mustDecodeHex(t, "01"), types.StoredProgram{
			Source:     "alpha.",
			CreatedAt:  10,
			SourceSize: 6,
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetStoredProgram(sdkCtx, mustDecodeHex(t, "02"), types.StoredProgram{
			Source:     "beta.",
			CreatedAt:  20,
			SourceSize: 5,
		})
		So(err, ShouldBeNil)

		publisherAddr, err := sdk.AccAddressFromBech32(publisher)
		So(err, ShouldBeNil)
		otherPublisherAddr, err := sdk.AccAddressFromBech32(otherPublisher)
		So(err, ShouldBeNil)

		err = logicKeeper.SetProgramPublication(sdkCtx, publisherAddr, mustDecodeHex(t, "01"), types.ProgramPublication{
			PublishedAt: 100,
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetProgramPublication(sdkCtx, publisherAddr, mustDecodeHex(t, "02"), types.ProgramPublication{
			PublishedAt: 200,
		})
		So(err, ShouldBeNil)
		err = logicKeeper.SetProgramPublication(sdkCtx, otherPublisherAddr, mustDecodeHex(t, "02"), types.ProgramPublication{
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
						ProgramId:  "01",
						CreatedAt:  10,
						SourceSize: 6,
					},
					Publication: types.ProgramPublication{PublishedAt: 100},
				})
				So(res.Programs[1], ShouldResemble, types.PublishedProgram{
					Program: types.ProgramMetadata{
						ProgramId:  "02",
						CreatedAt:  20,
						SourceSize: 5,
					},
					Publication: types.ProgramPublication{PublishedAt: 200},
				})
				So(res.Pagination, ShouldNotBeNil)
				So(res.Pagination.Total, ShouldEqual, 2)
			})
		})
	})
}

func TestQueryProgramErrors(t *testing.T) {
	Convey("Given a query keeper", t, func() {
		logicKeeper, sdkCtx := newQueryKeeper(t)

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
			_, err := logicKeeper.Program(sdkCtx, &types.QueryProgramRequest{ProgramId: "01"})

			Convey("Then it should return NotFound error", func() {
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.NotFound)
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

func newQueryKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
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

	return logicKeeper, testCtx.Ctx
}

func mustDecodeHex(t *testing.T, value string) []byte {
	t.Helper()

	bz, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("failed to decode hex: %v", err)
	}

	return bz
}

var _ context.Context = sdk.Context{}
