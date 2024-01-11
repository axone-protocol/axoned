//nolint:gocognit
package predicate

import (
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"

	"cosmossdk.io/store/metrics"

	dbm "github.com/cosmos/cosmos-db"

	"github.com/golang/mock/gomock"
	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
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
			wantResult     []types.TermResults
			wantError      error
		}{
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(100))),
					},
				},
				query:      `bank_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).`,
				wantResult: []types.TermResults{{"X": "[uknow-100]"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))),
					},
				},
				query:      `bank_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [X]).`,
				wantResult: []types.TermResults{{"X": "uatom-100"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [X, Y]).`,
				wantResult: []types.TermResults{{"X": "uatom-589", "Y": "uknow-420"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A) | _]).`,
				wantResult: []types.TermResults{{"D": "uatom", "A": "589"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				query:      `bank_balances('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', [_, X]).`,
				wantResult: []types.TermResults{{"X": "uknow-589"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				program:    `bank_balances_has_coin(A, D, V) :- bank_balances(A, R), member(D-V, R).`,
				query:      `bank_balances_has_coin('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', 'uknow', V).`,
				wantResult: []types.TermResults{{"V": "589"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins: sdk.NewCoins(
							sdk.NewCoin("uknow", math.NewInt(589)),
							sdk.NewCoin("uatom", math.NewInt(693)),
							sdk.NewCoin("uband", math.NewInt(4282)),
							sdk.NewCoin("uakt", math.NewInt(4099)),
							sdk.NewCoin("ukava", math.NewInt(836)),
							sdk.NewCoin("uscrt", math.NewInt(599)),
						),
					},
				},
				program:    `bank_balances_has_sufficient_coin(A, C, S) :- bank_balances(A, R), member(C, R), -(_, V) = C, compare(>, V, S).`,
				query:      `bank_balances_has_sufficient_coin('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', C, 600).`,
				wantResult: []types.TermResults{{"C": "uakt-4099"}, {"C": "uatom-693"}, {"C": "uband-4282"}, {"C": "ukava-836"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query: `bank_balances(Accounts, Balances).`,
				wantResult: []types.TermResults{
					{"Accounts": "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm", "Balances": "[uknow-420]"},
					{"Accounts": "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38", "Balances": "[uatom-589]"},
				},
			},
			{
				balances:   []bank.Balance{},
				query:      `bank_balances('foo', X).`,
				wantResult: []types.TermResults{{"X": "[uknow-100]"}},
				wantError: fmt.Errorf("error(resource_error(resource_module(bank)),[%s],unknown)",
					strings.Join(strings.Split("decoding bech32 failed: invalid bech32 string length 3", ""), ",")),
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(1000))),
					},
				},
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(100))),
					},
				},
				query:      `bank_spendable_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).`,
				wantResult: []types.TermResults{{"X": "[uknow-100]"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))),
					},
				},
				query:      `bank_spendable_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [X]).`,
				wantResult: []types.TermResults{{"X": "uatom-100"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_spendable_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [X, Y]).`,
				wantResult: []types.TermResults{{"X": "uatom-589", "Y": "uknow-420"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_spendable_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A) | _]).`,
				wantResult: []types.TermResults{{"D": "uatom", "A": "589"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				query:      `bank_spendable_balances('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', [_, X]).`,
				wantResult: []types.TermResults{{"X": "uknow-589"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				program:    `bank_spendable_has_coin(A, D, V) :- bank_spendable_balances(A, R), member(D-V, R).`,
				query:      `bank_spendable_has_coin('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', 'uknow', V).`,
				wantResult: []types.TermResults{{"V": "589"}},
			},
			{
				spendableCoins: []bank.Balance{
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins: sdk.NewCoins(
							sdk.NewCoin("uknow", math.NewInt(589)),
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
				query:      `bank_spendable_has_sufficient_coin('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', C, 600).`,
				wantResult: []types.TermResults{{"C": "uakt-4099"}, {"C": "uatom-693"}, {"C": "uband-4282"}, {"C": "ukava-836"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(1220))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(8000))),
					},
				},
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query: `bank_spendable_balances(Accounts, SpendableCoins).`,
				wantResult: []types.TermResults{
					{"Accounts": "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm", "SpendableCoins": "[uknow-420]"},
					{"Accounts": "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38", "SpendableCoins": "[uatom-589]"},
				},
			},
			{
				spendableCoins: []bank.Balance{},
				query:          `bank_spendable_balances('foo', X).`,
				wantResult:     []types.TermResults{{"X": "[uknow-100]"}},
				wantError: fmt.Errorf("error(resource_error(resource_module(bank)),[%s],unknown)",
					strings.Join(strings.Split("decoding bech32 failed: invalid bech32 string length 3", ""), ",")),
			},

			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(1000))),
					},
				},
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(100))),
					},
				},
				lockedCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(900))),
					},
				},
				query:      `bank_locked_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).`,
				wantResult: []types.TermResults{{"X": "[uknow-900]"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))),
					},
				},
				query:      `bank_locked_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [X]).`,
				wantResult: []types.TermResults{{"X": "uatom-100"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_locked_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [X, Y]).`,
				wantResult: []types.TermResults{{"X": "uatom-589", "Y": "uknow-420"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				query:      `bank_locked_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A) | _]).`,
				wantResult: []types.TermResults{{"D": "uatom", "A": "589"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				query:      `bank_locked_balances('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', [_, X]).`,
				wantResult: []types.TermResults{{"X": "uknow-589"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420)), sdk.NewCoin("uatom", math.NewInt(493))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(589)), sdk.NewCoin("uatom", math.NewInt(693))),
					},
				},
				program:    `bank_locked_has_coin(A, D, V) :- bank_locked_balances(A, R), member(D-V, R).`,
				query:      `bank_locked_has_coin('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', 'uknow', V).`,
				wantResult: []types.TermResults{{"V": "589"}},
			},
			{
				lockedCoins: []bank.Balance{
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins: sdk.NewCoins(
							sdk.NewCoin("uknow", math.NewInt(589)),
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
				query:      `bank_locked_has_sufficient_coin('okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38', C, 600).`,
				wantResult: []types.TermResults{{"C": "uakt-4099"}, {"C": "uatom-693"}, {"C": "uband-4282"}, {"C": "ukava-836"}},
			},
			{
				balances: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(1220))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(8000))),
					},
				},
				spendableCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(420))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(589))),
					},
				},
				lockedCoins: []bank.Balance{
					{
						Address: "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm",
						Coins:   sdk.NewCoins(sdk.NewCoin("uknow", math.NewInt(800))),
					},
					{
						Address: "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38",
						Coins:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(7411))),
					},
				},
				query: `bank_locked_balances(Accounts, LockedCoins).`,
				wantResult: []types.TermResults{
					{"Accounts": "okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm", "LockedCoins": "[uknow-800]"},
					{"Accounts": "okp41wze8mn5nsgl9qrgazq6a92fvh7m5e6pslyrz38", "LockedCoins": "[uatom-7411]"},
				},
			},
			{
				lockedCoins: []bank.Balance{},
				query:       `bank_locked_balances('foo', X).`,
				wantResult:  []types.TermResults{{"X": "[uknow-100]"}},
				wantError: fmt.Errorf("error(resource_error(resource_module(bank)),[%s],unknown)",
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
					sdk.GetConfig().SetBech32PrefixForAccount("okp4", "okp4pub")

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
							So(err, ShouldEqual, nil)

							Convey("When the predicate is called", func() {
								sols, err := interpreter.QueryContext(ctx, tc.query)

								Convey("Then the error should be nil", func() {
									So(err, ShouldEqual, nil)
									So(sols, ShouldNotBeNil)

									Convey("and the bindings should be as expected", func() {
										var got []types.TermResults
										for sols.Next() {
											m := types.TermResults{}
											err := sols.Scan(m)
											So(err, ShouldEqual, nil)

											got = append(got, m)
										}
										if tc.wantError != nil {
											So(sols.Err(), ShouldNotEqual, nil)
											So(sols.Err().Error(), ShouldEqual, tc.wantError.Error())
										} else {
											So(sols.Err(), ShouldEqual, nil)
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
