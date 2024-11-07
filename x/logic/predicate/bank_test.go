//nolint:gocognit
package predicate

import (
	"context"
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/engine"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/x/evidence"

	codecaddress "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/axone-protocol/axoned/v10/x/logic/testutil"
	"github.com/axone-protocol/axoned/v10/x/logic/types"
	"github.com/axone-protocol/axoned/v10/x/logic/util"
)

func TestBank(t *testing.T) {
	const (
		bench32DecodingFail = "d,e,c,o,d,i,n,g, ,b,e,c,h,3,2, ,f,a,i,l,e,d,:, ,i,n,v,a,l,i,d, ,b,e,c,h,3,2, ,s,t,r,i,n,g, ,l,e,n,g,t,h, ,3"
	)
	Convey("Under a mocked environment", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			ctx            context.Context
			balances       []bank.Balance
			spendableCoins []bank.Balance
			lockedCoins    []bank.Balance
			program        string
			query          string
			wantResult     []testutil.TermResults
			wantError      error
		}{
			{
				balances:   []bank.Balance{},
				query:      `bank_balances(X, Y).`,
				wantResult: nil,
			},
			{
				balances: []bank.Balance{
					{
						Address: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
						Coins:   sdk.NewCoins(sdk.NewCoin("uaxone", math.NewInt(100))),
					},
				},
				query:      `bank_balances('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', X).`,
				wantResult: nil,
			},
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
				wantError:  fmt.Errorf("error(domain_error(encoding(bech32),foo),[%s],bank_balances/2)", bench32DecodingFail),
			},
			{
				ctx:        context.Background(),
				balances:   []bank.Balance{},
				query:      `bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).`,
				wantResult: []testutil.TermResults{{"X": "[uaxone-100]"}},
				wantError:  fmt.Errorf("error(resource_error(resource_context(bankKeeper)),bank_balances/2)"),
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
				wantError:      fmt.Errorf("error(domain_error(encoding(bech32),foo),[%s],bank_spendable_balances/2)", bench32DecodingFail),
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
				wantError:   fmt.Errorf("error(domain_error(encoding(bech32),foo),[%s],bank_locked_balances/2)", bench32DecodingFail),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := dbm.NewMemDB()
					stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
					accountKeeper := testutil.NewMockAccountKeeper(ctrl)
					authQueryServiceKeeper := testutil.NewMockAuthQueryService(ctrl)
					bankKeeper := testutil.NewMockBankKeeper(ctrl)
					encCfg := moduletestutil.MakeTestEncodingConfig(evidence.AppModuleBasic{})
					sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")

					ctx := tc.ctx
					if ctx == nil {
						ctx = sdk.
							NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger()).
							WithValue(types.BankKeeperContextKey, bankKeeper).
							WithValue(types.AuthKeeperContextKey, accountKeeper).
							WithValue(types.AuthQueryServiceContextKey, authQueryServiceKeeper).
							WithValue(types.InterfaceRegistryContextKey, encCfg.InterfaceRegistry)
					}

					Convey("and a bank keeper initialized with the preconfigured balances", func(c C) {
						addresses := lo.Uniq(
							lo.Map(
								append(append(append(make([]bank.Balance, 0), tc.balances...), tc.spendableCoins...), tc.lockedCoins...),
								func(it bank.Balance, _ int) string {
									return it.Address
								}))

						accountKeeper.
							EXPECT().
							GetAccount(gomock.Any(), gomock.Any()).
							AnyTimes().
							DoAndReturn(func(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
								if _, ok := lo.Find(addresses, func(item string) bool {
									return addr.String() == item
								}); !ok {
									return nil
								}

								accAddr, err := codecaddress.NewBech32Codec("axone").StringToBytes(addr.String())
								c.So(err, ShouldBeNil)

								return authtypes.NewBaseAccountWithAddress(accAddr)
							})
						testutil.MockAuthQueryServiceWithAddresses(authQueryServiceKeeper, addresses)
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

//nolint:lll
func TestAccount(t *testing.T) {
	Convey("Under a mocked environment", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			ctx                         context.Context
			addresses                   []string
			authQueryServiceKeeperError bool
			program                     string
			query                       string
			wantAnswer                  *types.Answer
		}{
			{
				addresses: []string{
					"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
					"axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
				},
				program: `test_existence :- account(axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep), account(axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa).`,
				query:   `test_existence.`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{},
					Results: []types.Result{
						{
							Error:         "",
							Substitutions: []types.Substitution{},
						},
					},
				},
			},
			{
				addresses: []string{
					"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
				},
				query: `account(axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep).`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{},
					Results:   []types.Result{},
				},
			},
			{
				addresses: []string{
					"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
					"axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep",
				},
				query: `account(X).`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{"X"},
					Results: []types.Result{
						{
							Substitutions: []types.Substitution{
								{Variable: "X", Expression: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa"},
							},
						},
						{
							Substitutions: []types.Substitution{
								{Variable: "X", Expression: "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep"},
							},
						},
					},
				},
			},
			{
				addresses: []string{
					"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
					"axone1cchqkgswx4z35p65lpvet4wyz34058xr6nc73r",
					"axone1rnrm2ajmtwp80kw2tggm9rdn8lthvxhuy97xj4",
					"axone1p0t82zgu2gmklkc2wnwu8cg7d560uvh6nzvlnm",
					"axone13wtyx96tkhzge77kdquplpt8g7q3t3wgawrusn",
					"axone18kllsfrcp7pdymvg2tftvy4udwardfjx42rz60",
				},
				query: `account(X).`,
				wantAnswer: &types.Answer{
					HasMore:   true,
					Variables: []string{"X"},
					Results: []types.Result{
						{
							Substitutions: []types.Substitution{
								{Variable: "X", Expression: "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa"},
							},
						},
						{
							Substitutions: []types.Substitution{
								{Variable: "X", Expression: "axone1cchqkgswx4z35p65lpvet4wyz34058xr6nc73r"},
							},
						},
						{
							Substitutions: []types.Substitution{
								{Variable: "X", Expression: "axone1rnrm2ajmtwp80kw2tggm9rdn8lthvxhuy97xj4"},
							},
						},
						{
							Substitutions: []types.Substitution{
								{Variable: "X", Expression: "axone1p0t82zgu2gmklkc2wnwu8cg7d560uvh6nzvlnm"},
							},
						},
						{
							Substitutions: []types.Substitution{
								{Variable: "X", Expression: "axone13wtyx96tkhzge77kdquplpt8g7q3t3wgawrusn"},
							},
						},
					},
				},
			},
			{
				addresses: []string{
					"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
				},
				query: `account(unparseable).`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{},
					Results: []types.Result{
						{
							Error:         "error(domain_error(encoding(bech32),unparseable),[d,e,c,o,d,i,n,g, ,b,e,c,h,3,2, ,f,a,i,l,e,d,:, ,i,n,v,a,l,i,d, ,s,e,p,a,r,a,t,o,r, ,i,n,d,e,x, ,-,1],account/1)",
							Substitutions: nil,
						},
					},
				},
			},
			{
				addresses: []string{
					"axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa",
				},
				query: `account("wrong agument type").`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{},
					Results: []types.Result{
						{
							Error:         "error(type_error(atom,[w,r,o,n,g, ,a,g,u,m,e,n,t, ,t,y,p,e]),account/1)",
							Substitutions: nil,
						},
					},
				},
			},
			{
				ctx:       context.Background(),
				addresses: []string{},
				query:     `account(dont_care).`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{},
					Results: []types.Result{
						{
							Error:         "error(resource_error(resource_context(authKeeper)),account/1)",
							Substitutions: nil,
						},
					},
				},
			},
			{
				ctx:       context.WithValue(context.Background(), types.AuthKeeperContextKey, testutil.NewMockAccountKeeper(ctrl)),
				addresses: []string{},
				query:     `account(dont_care).`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{},
					Results: []types.Result{
						{
							Error:         "error(resource_error(resource_context(authQueryService)),account/1)",
							Substitutions: nil,
						},
					},
				},
			},
			{
				ctx: context.WithValue(
					context.WithValue(
						context.Background(),
						types.AuthKeeperContextKey, testutil.NewMockAccountKeeper(ctrl)),
					types.AuthQueryServiceContextKey, testutil.NewMockAuthQueryService(ctrl)),
				addresses: []string{},
				query:     `account(dont_care).`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{},
					Results: []types.Result{
						{
							Error:         "error(resource_error(resource_context(interfaceRegistry)),account/1)",
							Substitutions: nil,
						},
					},
				},
			},
			{
				addresses:                   []string{},
				authQueryServiceKeeperError: true,
				query:                       `account(X).`,
				wantAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{"X"},
					Results: []types.Result{
						{
							Error: "error(resource_error(resource_module(auth)),[r,p,c, ,e,r,r,o,r,:, ,c,o,d,e, ,=, ,P,e,r,m,i,s,s,i,o,n,D,e,n,i,e,d, ,d,e,s,c, ,=, ,n,o,t, ,a,l,l,o,w,e,d],account/1)",
						},
					},
				},
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					accountKeeper := testutil.NewMockAccountKeeper(ctrl)
					authQueryServiceKeeper := testutil.NewMockAuthQueryService(ctrl)
					encCfg := moduletestutil.MakeTestEncodingConfig(evidence.AppModuleBasic{})

					ctx := func() context.Context {
						if tc.ctx != nil {
							return tc.ctx
						}
						db := dbm.NewMemDB()
						stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

						return sdk.
							NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger()).
							WithValue(types.AuthKeeperContextKey, accountKeeper).
							WithValue(types.AuthQueryServiceContextKey, authQueryServiceKeeper).
							WithValue(types.InterfaceRegistryContextKey, encCfg.InterfaceRegistry)
					}()

					sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")

					Convey("and a set of existing addresses", func(c C) {
						accountKeeper.
							EXPECT().
							GetAccount(gomock.Any(), gomock.Any()).
							AnyTimes().
							DoAndReturn(func(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
								if _, ok := lo.Find(tc.addresses, func(item string) bool {
									return addr.String() == item
								}); !ok {
									return nil
								}

								accAddr, err := codecaddress.NewBech32Codec("axone").StringToBytes(addr.String())
								c.So(err, ShouldBeNil)

								return authtypes.NewBaseAccountWithAddress(accAddr)
							})

						if tc.authQueryServiceKeeperError {
							testutil.MockAuthQueryServiceWithError(authQueryServiceKeeper, status.Error(codes.PermissionDenied, "not allowed"))
						} else {
							testutil.MockAuthQueryServiceWithAddresses(authQueryServiceKeeper, tc.addresses)
						}

						Convey("and a vm with the account predicate registered", func() {
							interpreter := testutil.NewLightInterpreterMust(ctx)
							interpreter.Register1(engine.NewAtom("account"), account)

							err := interpreter.Compile(ctx, tc.program)
							So(err, ShouldBeNil)

							Convey("When the predicate is called", func() {
								answer, err := util.QueryInterpreter(ctx, interpreter, tc.query, math.NewUint(5))

								Convey("Then the error should be nil", func() {
									So(err, ShouldBeNil)
									So(answer, ShouldNotBeNil)

									Convey("and the answer should be as expected", func() {
										So(answer, ShouldResemble, tc.wantAnswer)
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
