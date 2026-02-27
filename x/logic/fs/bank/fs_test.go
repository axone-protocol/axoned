package bank

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/testutil"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

func TestVFS(t *testing.T) {
	addr := "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa"
	Convey("Given a bank VFS", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")

		mockBankKeeper := testutil.NewMockBankKeeper(ctrl)
		sdkCtx := newTestContext()
		ctx := context.WithValue(sdkCtx, types.BankKeeperContextKey, mockBankKeeper)

		bankFS := NewFS(ctx)

		Convey("When reading a valid balances file", func() {
			coins := sdk.NewCoins(
				sdk.NewCoin("uatom", math.NewInt(100)),
				sdk.NewCoin("uaxone", math.NewInt(200)),
			)

			mockBankKeeper.EXPECT().
				GetAllBalances(gomock.Any(), sdk.MustAccAddressFromBech32(addr)).
				Return(coins)

			data, err := bankFS.ReadFile(addr + "/balances/@")

			Convey("Then it should return the balances as Prolog terms", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldContainSubstring, "-(uatom,100)")
				So(string(data), ShouldContainSubstring, "-(uaxone,200)")
			})
		})

		Convey("When reading a valid spendable balances file", func() {
			coins := sdk.NewCoins(
				sdk.NewCoin("uatom", math.NewInt(80)),
				sdk.NewCoin("uaxone", math.NewInt(150)),
			)

			mockBankKeeper.EXPECT().
				SpendableCoins(gomock.Any(), sdk.MustAccAddressFromBech32(addr)).
				Return(coins)

			data, err := bankFS.ReadFile(addr + "/spendable/@")

			Convey("Then it should return the spendable balances as Prolog terms", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldContainSubstring, "-(uatom,80)")
				So(string(data), ShouldContainSubstring, "-(uaxone,150)")
			})
		})

		Convey("When reading with invalid paths", func() {
			cases := []struct {
				name    string
				path    string
				message string
			}{
				{name: "non-bank path", path: "invalid/path", message: "does not exist"},
				{name: "unknown balances collection", path: addr + "/unknown/@", message: "does not exist"},
				{name: "invalid bech32 address", path: "invalid_addr/balances/@", message: "does not exist"},
				{name: "missing address segment", path: "/balances/@", message: "does not exist"},
			}

			for _, tc := range cases {
				Convey("Then it should fail for "+tc.name, func() {
					_, err := bankFS.ReadFile(tc.path)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, tc.message)
				})
			}
		})

		Convey("When reading balances with an IBC denomination", func() {
			coins := sdk.NewCoins(
				sdk.NewCoin("ibc/0123456789ABCDEF", math.NewInt(100)),
			)

			mockBankKeeper.EXPECT().
				GetAllBalances(gomock.Any(), sdk.MustAccAddressFromBech32(addr)).
				Return(coins)

			data, err := bankFS.ReadFile(addr + "/balances/@")

			Convey("Then it should quote the denomination to keep the term readable", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldContainSubstring, "-('ibc/0123456789ABCDEF',100)")
			})
		})

		Convey("When reading balances with amount larger than int64", func() {
			amount, ok := math.NewIntFromString("9223372036854775808")
			So(ok, ShouldBeTrue)
			coins := sdk.NewCoins(
				sdk.NewCoin("uaxone", amount),
			)

			mockBankKeeper.EXPECT().
				GetAllBalances(gomock.Any(), sdk.MustAccAddressFromBech32(addr)).
				Return(coins)

			data, err := bankFS.ReadFile(addr + "/balances/@")

			Convey("Then it should serialize the amount as an atom to preserve precision", func() {
				So(err, ShouldBeNil)
				So(strings.TrimSpace(string(data)), ShouldEqual, "-(uaxone,'9223372036854775808').")
			})
		})

		Convey("When opening a file", func() {
			file, err := bankFS.Open(addr + "/balances/@")

			Convey("Then it should return a valid file", func() {
				So(err, ShouldBeNil)
				So(file, ShouldNotBeNil)
				Reset(func() {
					file.Close()
				})

				stat, err := file.Stat()
				So(err, ShouldBeNil)
				So(stat.Name(), ShouldEqual, "@")
			})
		})

		Convey("When the bank keeper is missing from context", func() {
			bankFS := NewFS(newTestContext())

			_, err := bankFS.Open(addr + "/balances/@")

			Convey("Then it should surface a term error", func() {
				So(err, ShouldNotBeNil)
				So(errors.Is(err, errVFSUnavailable), ShouldBeTrue)
			})
		})
	})
}

func newTestContext() sdk.Context {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	headerInfo := coreheader.Info{
		Height: 42,
		Time:   time.Date(2024, 4, 10, 10, 44, 27, 0, time.UTC),
	}

	return sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger()).WithHeaderInfo(headerInfo)
}
