package predicate

import (
	"context"
	"sort"

	"github.com/axone-protocol/prolog/v2/engine"
	"github.com/samber/lo"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/axone-protocol/axoned/v12/x/logic/prolog"
	"github.com/axone-protocol/axoned/v12/x/logic/types"
)

// SortBalances by coin denomination.
func SortBalances(balances sdk.Coins) {
	sort.SliceStable(balances, func(i, j int) bool {
		return balances[i].Denom < balances[j].Denom
	})
}

// IterMap transforms the output of an iterator by applying a given mapping function to each element.
func IterMap[U any, V any](next func() (U, bool), f func(U) V) func() (V, bool) {
	return func() (V, bool) {
		u, ok := next()
		if !ok {
			var zeroV V
			return zeroV, false
		}

		return f(u), true
	}
}

// Accounts returns an iterator that iterates over all accounts.
//
//nolint:lll
func Accounts(ctx context.Context, authQueryService types.AuthQueryService, unpacker cdctypes.AnyUnpacker) func() (lo.Tuple2[sdk.AccountI, error], bool) {
	var (
		finished bool
		key      []byte
	)

	return func() (lo.Tuple2[sdk.AccountI, error], bool) {
		if finished {
			return lo.Tuple2[sdk.AccountI, error]{A: nil, B: nil}, false
		}
		res, err := authQueryService.Accounts(ctx,
			&auth.QueryAccountsRequest{
				Pagination: &query.PageRequest{
					Key:   key,
					Limit: 1,
				},
			})
		if err != nil {
			finished = true
			return lo.Tuple2[sdk.AccountI, error]{A: nil, B: err}, true
		}

		if len(res.Accounts) == 0 {
			finished = true
			return lo.Tuple2[sdk.AccountI, error]{A: nil, B: nil}, false
		}

		key = res.Pagination.NextKey
		finished = len(key) == 0

		var account sdk.AccountI
		if err := unpacker.UnpackAny(res.Accounts[0], &account); err != nil {
			return lo.Tuple2[sdk.AccountI, error]{A: nil, B: err}, true
		}

		return lo.Tuple2[sdk.AccountI, error]{A: account, B: nil}, true
	}
}

// AllBalancesSorted returns the list of balances for the given address, sorted by coin denomination.
func AllBalancesSorted(ctx context.Context, bankKeeper types.BankKeeper, bech32Addr sdk.AccAddress) sdk.Coins {
	fetchedBalances := bankKeeper.GetAllBalances(ctx, bech32Addr)
	SortBalances(fetchedBalances)
	return fetchedBalances
}

// SpendableCoinsSorted returns the list of spendable coins for the given address, sorted by coin denomination.
func SpendableCoinsSorted(ctx context.Context, bankKeeper types.BankKeeper, bech32Addr sdk.AccAddress) sdk.Coins {
	fetchedBalances := bankKeeper.SpendableCoins(ctx, bech32Addr)
	SortBalances(fetchedBalances)
	return fetchedBalances
}

// LockedCoinsSorted returns the list of spendable coins for the given address, sorted by coin denomination.
func LockedCoinsSorted(ctx context.Context, bankKeeper types.BankKeeper, bech32Addr sdk.AccAddress) sdk.Coins {
	fetchedBalances := bankKeeper.LockedCoins(ctx, bech32Addr)
	SortBalances(fetchedBalances)
	return fetchedBalances
}

// CoinsToTerm converts the given coins to a term of the form:
//
//	[-(Denom, Amount), -(Denom, Amount), ...]
func CoinsToTerm(coins sdk.Coins) engine.Term {
	terms := make([]engine.Term, 0, len(coins))
	for _, coin := range coins {
		terms = append(terms, prolog.AtomPair.Apply(engine.NewAtom(coin.Denom), engine.Integer(coin.Amount.Int64())))
	}

	return engine.List(terms...)
}
