package predicate

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/types"
)

// SortBalances by coin denomination.
func SortBalances(balances sdk.Coins) {
	sort.SliceStable(balances, func(i, j int) bool {
		return balances[i].Denom < balances[j].Denom
	})
}

// AllBalancesSorted returns the list of balances for the given address, sorted by coin denomination.
func AllBalancesSorted(sdkContext sdk.Context, bankKeeper types.BankKeeper, bech32Addr sdk.AccAddress) sdk.Coins {
	fetchedBalances := bankKeeper.GetAllBalances(sdkContext, bech32Addr)
	SortBalances(fetchedBalances)
	return fetchedBalances
}

// SpendableCoinsSorted returns the list of spendable coins for the given address, sorted by coin denomination.
func SpendableCoinsSorted(sdkContext sdk.Context, bankKeeper types.BankKeeper, bech32Addr sdk.AccAddress) sdk.Coins {
	fetchedBalances := bankKeeper.SpendableCoins(sdkContext, bech32Addr)
	SortBalances(fetchedBalances)
	return fetchedBalances
}

// LockedCoinsSorted returns the list of spendable coins for the given address, sorted by coin denomination.
func LockedCoinsSorted(sdkContext sdk.Context, bankKeeper types.BankKeeper, bech32Addr sdk.AccAddress) sdk.Coins {
	fetchedBalances := bankKeeper.LockedCoins(sdkContext, bech32Addr)
	SortBalances(fetchedBalances)
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

func BytesToList(bt []byte) engine.Term {
	terms := make([]engine.Term, 0, len(bt))
	for _, b := range bt {
		terms = append(terms, engine.Integer(b))
	}
	return engine.List(terms...)
}

func ListToBytes(terms engine.ListIterator, env *engine.Env) ([]byte, error) {
	bt := make([]byte, 0)
	for terms.Next() {
		term := env.Resolve(terms.Current())
		switch t := term.(type) {
		case engine.Integer:
			bt = append(bt, byte(t))
		default:
			return nil, fmt.Errorf("invalid term type in list %T, only integer allowed", term)
		}
	}
	return bt, nil
}

func AtomBool(b bool) engine.Term {
	var r engine.Atom
	if b {
		r = engine.NewAtom("true")
	} else {
		r = engine.NewAtom("false")
	}
	return engine.NewAtom("@").Apply(r)
}

var AtomNull = engine.NewAtom("@").Apply(engine.NewAtom("null"))
