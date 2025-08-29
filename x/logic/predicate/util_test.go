package predicate

import (
	"errors"
	"fmt"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/x/evidence"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/axone-protocol/axoned/v12/x/logic/testutil"
	"github.com/axone-protocol/axoned/v12/x/logic/types"
)

//nolint:gocognit
func TestAccounts(t *testing.T) {
	Convey("Under a mocked environment", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Convey("Given test cases", func() {
			cases := []struct {
				addresses                   []string
				authQueryServiceKeeperError string
				interfaceRegistryError      string
			}{
				{
					addresses: []string{},
				},
				{
					addresses: []string{
						"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						"axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
					},
				},
				{
					addresses: []string{
						"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						"axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
					},
					interfaceRegistryError: "can't unpack",
				},
				{
					addresses: []string{
						"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						"axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
					},
					authQueryServiceKeeperError: "i/o error",
				},
			}
			for nc, tc := range cases {
				Convey(fmt.Sprintf("Given the test case #%d", nc), func() {
					Convey("and a context", func() {
						sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")
						encCfg := moduletestutil.MakeTestEncodingConfig(evidence.AppModuleBasic{})

						authQueryServiceKeeper := testutil.NewMockAuthQueryService(ctrl)
						if tc.authQueryServiceKeeperError == "" {
							testutil.MockAuthQueryServiceWithAddresses(authQueryServiceKeeper, tc.addresses)
						} else {
							testutil.MockAuthQueryServiceWithError(authQueryServiceKeeper, errors.New(tc.authQueryServiceKeeperError))
						}

						interfaceRegistry := testutil.NewMockInterfaceRegistry(ctrl)

						interfaceRegistry.
							EXPECT().
							UnpackAny(gomock.Any(), gomock.Any()).
							AnyTimes().
							DoAndReturn(func(v *cdctypes.Any, iface any) error {
								if tc.interfaceRegistryError != "" {
									return errors.New(tc.interfaceRegistryError)
								}
								return encCfg.InterfaceRegistry.UnpackAny(v, iface)
							})

						db := dbm.NewMemDB()
						stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

						ctx := sdk.
							NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger()).
							WithValue(types.AuthQueryServiceContextKey, authQueryServiceKeeper).
							WithValue(types.InterfaceRegistryContextKey, interfaceRegistry)

						Convey("When Accounts is called", func() {
							next := Accounts(ctx, authQueryServiceKeeper, interfaceRegistry)

							Convey("Then next() behave as expected", func() {
								for _, wantAddr := range tc.addresses {
									result, ok := next()

									So(ok, ShouldBeTrue)
									if tc.interfaceRegistryError == "" && tc.authQueryServiceKeeperError == "" {
										So(result.A, ShouldNotBeNil)
										So(result.A.GetAddress().String(), ShouldEqual, wantAddr)
										So(result.B, ShouldBeNil)
									} else {
										So(result.A, ShouldBeNil)
										if tc.authQueryServiceKeeperError != "" {
											So(result.B, ShouldBeError, tc.authQueryServiceKeeperError)
											break
										}

										So(result.B, ShouldBeError, tc.interfaceRegistryError)
									}
								}

								for range 5 {
									result, ok := next()
									So(ok, ShouldBeFalse)
									So(result.A, ShouldBeNil)
									So(result.B, ShouldBeNil)
								}
							})
						})
					})
				})
			}
		})
	})
}
