package prolog

import (
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestResolveHeaderInfo(t *testing.T) {
	Convey("Given a context with explicit HeaderInfo", t, func() {
		headerTime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
		ctx := newSDKContext(cmtproto.Header{
			Height:  7,
			Time:    headerTime,
			ChainID: "legacy-chain",
			AppHash: []byte{0xAA},
		}).WithHeaderHash([]byte{0x55}).
			WithHeaderInfo(coreheader.Info{
				Height:  42,
				Time:    time.Date(2030, 6, 7, 8, 9, 10, 0, time.UTC),
				ChainID: "header-info-chain",
				Hash:    []byte{0x01, 0x02},
				AppHash: []byte{0x03, 0x04},
			})

		Convey("When resolving header info", func() {
			got := ResolveHeaderInfo(ctx)

			Convey("Then HeaderInfo values should be used", func() {
				So(got.Height, ShouldEqual, 42)
				So(got.Time, ShouldEqual, time.Date(2030, 6, 7, 8, 9, 10, 0, time.UTC))
				So(got.ChainID, ShouldEqual, "header-info-chain")
				So(got.Hash, ShouldResemble, []byte{0x01, 0x02})
				So(got.AppHash, ShouldResemble, []byte{0x03, 0x04})
			})

			Convey("And slices should be copied", func() {
				got.Hash[0] = 0xFF
				got.AppHash[0] = 0xEE

				headerInfo := ctx.HeaderInfo()
				So(headerInfo.Hash, ShouldResemble, []byte{0x01, 0x02})
				So(headerInfo.AppHash, ShouldResemble, []byte{0x03, 0x04})
			})
		})
	})

	Convey("Given a context without HeaderInfo", t, func() {
		ctx := newSDKContext(cmtproto.Header{
			Height:  100,
			Time:    time.Date(2024, 3, 4, 11, 3, 36, 0, time.UTC),
			ChainID: "axone-testchain-1",
			AppHash: []byte{0x0A, 0x0B},
		}).WithHeaderHash([]byte{0xCA, 0xFE})

		Convey("When resolving header info", func() {
			got := ResolveHeaderInfo(ctx)

			Convey("Then no fallback should be applied", func() {
				So(got.Height, ShouldEqual, 0)
				So(got.Time.IsZero(), ShouldBeTrue)
				So(got.ChainID, ShouldEqual, "")
				So(got.Hash, ShouldBeNil)
				So(got.AppHash, ShouldBeNil)
			})
		})
	})
}

func newSDKContext(header cmtproto.Header) sdk.Context {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	return sdk.NewContext(stateStore, header, false, log.NewNopLogger())
}
