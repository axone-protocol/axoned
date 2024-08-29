//nolint:gocognit
package predicate

import (
	"fmt"
	"strings"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/golang/mock/gomock"
	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/axone-protocol/axoned/v10/x/logic/testutil"
	"github.com/axone-protocol/axoned/v10/x/logic/types"
)

func TestBank(t *testing.T) {
	Convey("Under a mocked environment", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			balances       []bank.Balance
			spendableCoins []bank.Balance
			lockedCoins    []bank.Balance
			program        string
			query          string
			wantResult     []testutil.TermResults
			wantError      error
		}{
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(100))),
					},
				},
				query:      `bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).`,
				wantResult: []testutil.TermResults{{"X": "[uaxone-100]"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))),
					},
				},
				query:      `bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [X]).`,
				wantResult: []testutil.TermResults{{"X": "uatom-100"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [X, Y]).`,
				wantResult: []testutil.TermResults{{"X": "uatom-589", "Y": "uaxone-420"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [-(D, A) | _]).`,
				wantResult: []testutil.TermResults{{"D": "uatom", "A": "589"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				query:      `bank_balances('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', [_, X]).`,
				wantResult: []testutil.TermResults{{"X": "uaxone-589"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				program:    `bank_balances_has_coin(A, D, V) :- bank_balances(A, R), member(D-V, R).`,
				query:      `bank_balances_has_coin('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', 'uaxone', V).`,
				wantResult: []testutil.TermResults{{"V": "589"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins: sdk.NewCoins(
							sdk.NewCoin("uaxone", math.NewInt(589)),
							sdk.NewCoin("uatom", math.NewInt(693)),
							sdk.NewCoin("uband", math.NewInt(4282)),
							sdk.NewCoin("uakt", math.NewInt(4099)),
							sdk.NewCoin("ukava", math.NewInt(836)),
							sdk.NewCoin("uscrt", math.NewInt(599)),
						),
					},
				},
				program:    `bank_balances_has_sufficient_coin(A, C, S) :- bank_balances(A, R), member(C, R), -(_, V) = C, compare(>, V, S).`,
				query:      `bank_balances_has_sufficient_coin('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', C, 600).`,
				wantResult: []testutil.TermResults{{"C": "uakt-4099"}, {"C": "uatom-693"}, {"C": "uband-4282"}, {"C": "ukava-836"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query: `bank_balances(Accounts, Balances).`,
				wantResult: []testutil.TermResults{
					{"Accounts": "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa", "Balances": "[uaxone-420]"},
					{"Accounts": "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep", "Balances": "[uatom-589]"},
				},
			},
			{
				balances:   []bank.Balance{},
				query:      `bank_balances('foo', X).`,
				wantResult: []testutil.TermResults{{"X": "[uaxone-100]"}},
				wantError: fmt.Errorf("error(resource_error(resource_module(bank)),[%s],bank_balances/2)",
					strings.Join(strings.Split("decoding bech32 failed: invalid bech32 string length 3", ""), ",")),
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(1000))),
					},
				},
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(100))),
					},
				},
				query:      `bank_spendable_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).`,
				wantResult: []testutil.TermResults{{"X": "[uaxone-100]"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))),
					},
				},
				query:      `bank_spendable_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [X]).`,
				wantResult: []testutil.TermResults{{"X": "uatom-100"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_spendable_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [X, Y]).`,
				wantResult: []testutil.TermResults{{"X": "uatom-589", "Y": "uaxone-420"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_spendable_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [-(D, A) | _]).`,
				wantResult: []testutil.TermResults{{"D": "uatom", "A": "589"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				query:      `bank_spendable_balances('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', [_, X]).`,
				wantResult: []testutil.TermResults{{"X": "uaxone-589"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				program:    `bank_spendable_has_coin(A, D, V) :- bank_spendable_balances(A, R), member(D-V, R).`,
				query:      `bank_spendable_has_coin('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', 'uaxone', V).`,
				wantResult: []testutil.TermResults{{"V": "589"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins: sdk.NewCoins(
							sdk.NewCoin("uaxone", math.NewInt(589)),
							sdk.NewCoin("uatom", math.NewInt(693)),
							sdk.NewCoin("uband", math.NewInt(4282)),
							sdk.NewCoin("uakt", math.NewInt(4099)),
							sdk.NewCoin("ukava", math.NewInt(836)),
							sdk.NewCoin("uscrt", math.NewInt(599)),
						),
					},
				},
				program: `bank_spendable_has_sufficient_coin(A, C, S) :- bank_spendable_balances(A, R), member(C, R),
-(_, V) = C, compare(>, V, S).`,
				query:      `bank_spendable_has_sufficient_coin('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', C, 600).`,
				wantResult: []testutil.TermResults{{"C": "uakt-4099"}, {"C": "uatom-693"}, {"C": "uband-4282"}, {"C": "ukava-836"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(1220))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(8000))),
					},
				},
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query: `bank_spendable_balances(Accounts, SpendableCoins).`,
				wantResult: []testutil.TermResults{
					{"Accounts": "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa", "SpendableCoins": "[uaxone-420]"},
					{"Accounts": "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep", "SpendableCoins": "[uatom-589]"},
				},
			},
			{
				spendableCoins: []bank.Balance{},
				query:          `bank_spendable_balances('foo', X).`,
				wantResult:     []testutil.TermResults{{"X": "[uaxone-100]"}},
				wantError: fmt.Errorf("error(resource_error(resource_module(bank)),[%s],bank_spendable_balances/2)",
					strings.Join(strings.Split("decoding bech32 failed: invalid bech32 string length 3", ""), ",")),
			},

			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(1000))),
					},
				},
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(100))),
					},
				},
				lockedCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(900))),
					},
				},
				query:      `bank_locked_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).`,
				wantResult: []testutil.TermResults{{"X": "[uaxone-900]"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))),
					},
				},
				query:      `bank_locked_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [X]).`,
				wantResult: []testutil.TermResults{{"X": "uatom-100"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_locked_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [X, Y]).`,
				wantResult: []testutil.TermResults{{"X": "uatom-589", "Y": "uaxone-420"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_locked_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [-(D, A) | _]).`,
				wantResult: []testutil.TermResults{{"D": "uatom", "A": "589"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				query:      `bank_locked_balances('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', [_, X]).`,
				wantResult: []testutil.TermResults{{"X": "uaxone-589"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				program:    `bank_locked_has_coin(A, D, V) :- bank_locked_balances(A, R), member(D-V, R).`,
				query:      `bank_locked_has_coin('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', 'uaxone', V).`,
				wantResult: []testutil.TermResults{{"V": "589"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins: sdk.NewCoins(
							sdk.NewCoin("uaxone", math.NewInt(589)),
							sdk.NewCoin("uatom", math.NewInt(693)),
							sdk.NewCoin("uband", math.NewInt(4282)),
							sdk.NewCoin("uakt", math.NewInt(4099)),
							sdk.NewCoin("ukava", math.NewInt(836)),
							sdk.NewCoin("uscrt", math.NewInt(599)),
						),
					},
				},
				program: `bank_locked_has_sufficient_coin(A, C, S) :- bank_locked_balances(A, R), member(C, R),
-(_, V) = C, compare(>, V, S).`,
				query:      `bank_locked_has_sufficient_coin('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', C, 600).`,
				wantResult: []testutil.TermResults{{"C": "uakt-4099"}, {"C": "uatom-693"}, {"C": "uband-4282"}, {"C": "ukava-836"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(1220))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(8000))),
					},
				},
				spendableCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(420))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				lockedCoins: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(800))),
					},
					{
						Address: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(7411))),
					},
				},
				query: `bank_locked_balances(Accounts, LockedCoins).`,
				wantResult: []testutil.TermResults{
					{"Accounts": "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa", "LockedCoins": "[uaxone-800]"},
					{"Accounts": "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep", "LockedCoins": "[uatom-7411]"},
				},
			},
			{
				lockedCoins: []bank.Balance{},
				query:       `bank_locked_balances('foo', X).`,
				wantResult:  []testutil.TermResults{{"X": "[uaxone-100]"}},
				wantError: fmt.Errorf("error(resource_error(resource_module(bank)),[%s],bank_locked_balances/2)",
					strings.Join(strings.Split("decoding bech32 failed: invalid bech32 string length 3", ""), ",")),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := dbm.NewMemDB()
					stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
					bankKeeper := testutil.NewMockBankKeeper(ctrl)
					ctx := sdk.
						NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger()).
						WithValue(types.BankKeeperContextKey, bankKeeper)
					sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")

					Convey("and a bank keeper initialized with the preconfigured balances", func() {
						for _, balance := range tc.balances {
							bankKeeper.
								EXPECT().
								GetAllBalances(ctx, sdk.MustAccAddressFromBech32(balance.Address)).
								AnyTimes().
								Return(balance.Coins)
						}
						for _, balance := range tc.spendableCoins {
							bankKeeper.
								EXPECT().
								SpendableCoins(ctx, sdk.MustAccAddressFromBech32(balance.Address)).
								AnyTimes().
								Return(balance.Coins)
						}
						for _, balance := range tc.lockedCoins {
							bankKeeper.
								EXPECT().
								LockedCoins(ctx, sdk.MustAccAddressFromBech32(balance.Address)).
								AnyTimes().
								Return(balance.Coins)
						}
						bankKeeper.
							EXPECT().
							GetAccountsBalances(ctx).
							AnyTimes().
							Return(tc.balances)

						Convey("and a vm", func() {
							interpreter := testutil.NewLightInterpreterMust(ctx)
							interpreter.Register2(engine.NewAtom("bank_balances"), BankBalances)
							interpreter.Register2(engine.NewAtom("bank_spendable_balances"), BankSpendableBalances)
							interpreter.Register2(engine.NewAtom("bank_locked_balances"), BankLockedBalances)

							err := interpreter.Compile(ctx, tc.program)
							So(err, ShouldBeNil)

							Convey("When the predicate is called", func() {
								sols, err := interpreter.QueryContext(ctx, tc.query)

								Convey("Then the error should be nil", func() {
									So(err, ShouldBeNil)
									So(sols, ShouldNotBeNil)

									Convey("and the bindings should be as expected", func() {
										var got []testutil.TermResults
										for sols.Next() {
											m := testutil.TermResults{}
											err := sols.Scan(m)
											So(err, ShouldBeNil)

											got = append(got, m)
										}
										if tc.wantError != nil {
											So(sols.Err(), ShouldNotBeNil)
											So(sols.Err().Error(), ShouldEqual, tc.wantError.Error())
										} else {
											So(sols.Err(), ShouldBeNil)
											So(got, ShouldResemble, tc.wantResult)
										}
									})
								})
							})
						})
					})
				})
			})
		}
	})
}
