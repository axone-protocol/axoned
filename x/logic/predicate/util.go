package predicate

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/types"
)

// BalancesSorted returns the list of balances for the given address, sorted by coin denomination.
func BalancesSorted(sdkContext sdk.Context, bankKeeper types.BankKeeper, bech32Addr sdk.AccAddress) sdk.Coins {
	fetchedBalances := bankKeeper.GetAllBalances(sdkContext, bech32Addr)
	sort.SliceStable(fetchedBalances, func(i, j int) bool {
		return fetchedBalances[i].Denom < fetchedBalances[j].Denom
	})

	return fetchedBalances
}

// CoinsToTerm converts the given coins to a term of the form:
//
//	[-(Denom, Amount), -(Denom, Amount), ...]
func CoinsToTerm(coins sdk.Coins) engine.Term {
	terms := make([]engine.Term, 0, len(coins))
	for _, coin := range coins {
		terms = append(terms, AtomPair.Apply(engine.NewAtom(coin.Denom), engine.Integer(coin.Amount.Int64())))
	}

	return engine.List(terms...)
}

// Tuple is a predicate which unifies the given term with a tuple of the given arity.
func Tuple(args ...engine.Term) engine.Term {
	return engine.Atom(0).Apply(args...)
}
