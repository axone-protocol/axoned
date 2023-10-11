package predicate

import (
	"encoding/hex"
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

func OptionsContains(atom engine.Atom, options engine.Term, env *engine.Env) (engine.Compound, error) {
	switch opts := env.Resolve(options).(type) {
	case engine.Compound:
		if opts.Functor() == atom {
			return opts, nil
		} else if opts.Arity() == 2 && opts.Functor().String() == "." {
			iter := engine.ListIterator{List: opts, Env: env}

			for iter.Next() {
				opt := env.Resolve(iter.Current())
				term, err := OptionsContains(atom, opt, env)
				if err != nil {
					return nil, err
				}
				if term != nil {
					return term, nil
				}
			}
		}
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid options term, should be compound, give %T", opts)
	}
}

// TermToBytes try to convert a term to native golang []byte.
// By default, if no encoding options is given the term is considered as hexadecimal value.
// Available encoding option is `text`, `octet` and `hex` (default value)
func TermToBytes(term, options engine.Term, env *engine.Env) ([]byte, error) {
	encoding, err := OptionsContains(AtomEncoding, options, env)
	if err != nil {
		return nil, err
	}

	if encoding == nil {
		encoding = AtomEncoding.Apply(engine.NewAtom("hex")).(engine.Compound)
	}

	if encoding.Arity() != 1 {
		return nil, fmt.Errorf("invalid arity for encoding option, should be 1")
	}

	switch enc := env.Resolve(encoding.Arg(0)).(type) {
	case engine.Atom:
		switch enc.String() {
		case "octet":
			switch b := env.Resolve(term).(type) {
			case engine.Compound:
				if b.Arity() != 2 || b.Functor().String() != "." {
					return nil, fmt.Errorf("term should be a List, give %T", b)
				}
				iter := engine.ListIterator{List: b, Env: env}

				return ListToBytes(iter, env)
			default:
				return nil, fmt.Errorf("invalid term type: %T, should be a List", term)
			}
		case "hex":
			switch b := env.Resolve(term).(type) {
			case engine.Atom:
				src := []byte(b.String())
				result := make([]byte, hex.DecodedLen(len(src)))
				_, err := hex.Decode(result, src)
				return result, err
			default:
				return nil, fmt.Errorf("invalid term type: %T, should be String", term)
			}
		default:
			return nil, fmt.Errorf("invalid encoding option: %s, valid value are 'hex' or 'octet'", enc.String())
		}
	default:
		return nil, fmt.Errorf("invalid given options")
	}
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
