package predicate

import (
	"fmt"
	"sort"

	"github.com/ichiban/prolog/engine"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

// ExtractJSONTerm is an utility function that would extract all attribute of a JSON object
// that is represented in prolog with the `json` atom.
//
// This function will ensure the json atom follow our json object representation in prolog.
//
// A JSON object is represented like this :
//
// ```
// json([foo-bar])
// ```
//
// That give a JSON object: `{"foo": "bar"}`
// Returns the map of all attributes with its term value.
func ExtractJSONTerm(term engine.Compound, env *engine.Env) (map[string]engine.Term, error) {
	if term.Functor() != AtomJSON {
		return nil, fmt.Errorf("invalid functor %s. Expected %s", term.Functor().String(), AtomJSON.String())
	} else if term.Arity() != 1 {
		return nil, fmt.Errorf("invalid compound arity : %d but expected %d", term.Arity(), 1)
	}

	list := term.Arg(0)
	switch l := env.Resolve(list).(type) {
	case engine.Compound:
		iter := engine.ListIterator{
			List: l,
			Env:  env,
		}
		terms := make(map[string]engine.Term, 0)
		for iter.Next() {
			pair, ok := env.Resolve(iter.Current()).(engine.Compound)
			if !ok || pair.Functor() != AtomPair || pair.Arity() != 2 {
				return nil, fmt.Errorf("json attributes should be a pair")
			}

			key, ok := env.Resolve(pair.Arg(0)).(engine.Atom)
			if !ok {
				return nil, fmt.Errorf("first pair arg should be an atom")
			}
			terms[key.String()] = pair.Arg(1)
		}
		return terms, nil
	default:
		return nil, fmt.Errorf("json compound should contains one list, give %T", l)
	}
}
