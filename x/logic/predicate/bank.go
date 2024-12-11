package predicate

import (
	"context"

	"github.com/axone-protocol/prolog/v2/engine"
	"github.com/samber/lo"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v11/x/logic/prolog"
	"github.com/axone-protocol/axoned/v11/x/logic/types"
)

// BankBalances is a predicate which unifies the given terms with the list of balances (coins) of the given account.
//
// The signature is as follows:
//
//	bank_balances(?Address, ?Balances)
//
// where:
//   - Address represents the account address (in Bech32 format).
//   - Balances represents the balances of the account as a list of pairs of coin denomination and amount.
//
// # Examples:
//
//	# Query the balances of the account.
//	- bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).
//
//	# Query the balances of all accounts. The result is a list of pairs of account address and balances.
//	- bank_balances(X, Y).
//
//	# Query the first balance of the given account by unifying the denomination and amount with the given terms.
//	- bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [-(D, A), _]).
func BankBalances(vm *engine.VM, address, balances engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return fetchBalances(vm, address, balances, AllBalancesSorted, cont, env)
}

// BankSpendableBalances is a predicate which unifies the given terms with the list of spendable coins of the given account.
//
// The signature is as follows:
//
//	bank_spendable_balances(?Address, ?Balances)
//
// where:
//   - Address represents the account address (in Bech32 format).
//   - Balances represents the spendable balances of the account as a list of pairs of coin denomination and amount.
//
// # Examples:
//
//	# Query the spendable balances of the account.
//	- bank_spendable_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).
//
//	# Query the spendable balances of all accounts. The result is a list of pairs of account address and balances.
//	- bank_spendable_balances(X, Y).
//
//	# Query the first spendable balances of the given account by unifying the denomination and amount with the given terms.
//	- bank_spendable_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [-(D, A), _]).
func BankSpendableBalances(vm *engine.VM, address, balances engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return fetchBalances(vm, address, balances, SpendableCoinsSorted, cont, env)
}

// BankLockedBalances is a predicate which unifies the given terms with the list of locked coins of the given account.
//
// The signature is as follows:
//
//	bank_locked_balances(?Address, ?Balances)
//
// where:
//   - Address represents the account address (in Bech32 format).
//   - Balances represents the locked balances of the account as a list of pairs of coin denomination and amount.
//
// # Examples:
//
//	# Query the locked coins of the account.
//	- bank_locked_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).
//
//	# Query the locked balances of all accounts. The result is a list of pairs of account address and balances.
//	- bank_locked_balances(X, Y).
//
//	# Query the first locked balances of the given account by unifying the denomination and amount with the given terms.
//	- bank_locked_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [-(D, A), _]).
func BankLockedBalances(vm *engine.VM, address, balances engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return fetchBalances(vm, address, balances, LockedCoinsSorted, cont, env)
}

// account is a predicate which unifies the given term with the list of account addresses existing in the blockchain.
//
// The signature is as follows:
//
//	account(?Address)
//
// where:
//   - Address represents the account address (in Bech32 format).
//
// # Examples:
//
//	# Query the locked coins of the account.
//	- account('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa').
//
//	# Query the all accounts existing in the blockchain.
//	- account(Address).
func account(vm *engine.VM, address engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		authKeeper, err := prolog.ContextValue[types.AccountKeeper](ctx, types.AuthKeeperContextKey, env)
		if err != nil {
			return engine.Error(err)
		}
		authQueryService, err := prolog.ContextValue[types.AuthQueryService](ctx, types.AuthQueryServiceContextKey, env)
		if err != nil {
			return engine.Error(err)
		}
		interfaceRegistry, err := prolog.ContextValue[cdctypes.InterfaceRegistry](ctx, types.InterfaceRegistryContextKey, env)
		if err != nil {
			return engine.Error(err)
		}

		switch acc := env.Resolve(address).(type) {
		case engine.Atom:
			return engine.Delay(
				func(ctx context.Context) *engine.Promise {
					addr, err := sdk.AccAddressFromBech32(acc.String())
					if err != nil {
						return engine.Error(prolog.WithError(
							engine.DomainError(prolog.ValidEncoding("bech32"), engine.NewAtom(acc.String()), env), err, env))
					}
					if exists := authKeeper.GetAccount(ctx, addr) != nil; !exists {
						return engine.Bool(false)
					}

					return cont(env)
				})
		case engine.Variable:
			return engine.DelaySeq(IterMap(Accounts(ctx, authQueryService, interfaceRegistry),
				func(it lo.Tuple2[sdk.AccountI, error]) engine.PromiseFunc {
					return func(_ context.Context) *engine.Promise {
						addr, err := lo.Unpack2(it)
						if err != nil {
							return engine.Error(prolog.WithError(engine.ResourceError(prolog.ResourceModule("auth"), env), err, env))
						}
						return engine.Unify(vm, address, engine.NewAtom(addr.GetAddress().String()), cont, env)
					}
				}))
		default:
			return engine.Error(engine.TypeError(prolog.AtomTypeAtom, address, env))
		}
	})
}

// fetchBalances is a helper function to fetch the balances of the given account using a given function which returns the coins for
// the given address.
func fetchBalances(vm *engine.VM, address engine.Term, balances engine.Term, coinsFn func(ctx context.Context,
	bankKeeper types.BankKeeper, address sdk.AccAddress) sdk.Coins, cont engine.Cont, env *engine.Env,
) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		bankKeeper, err := prolog.ContextValue[types.BankKeeper](ctx, types.BankKeeperContextKey, env)
		if err != nil {
			return engine.Error(err)
		}

		return account(vm, address, func(env *engine.Env) *engine.Promise {
			switch acc := env.Resolve(address).(type) {
			case engine.Atom:
				coins := coinsFn(ctx, bankKeeper, sdk.MustAccAddressFromBech32(acc.String()))

				return engine.Unify(vm, CoinsToTerm(coins), balances, cont, env)
			default:
				return engine.Error(engine.TypeError(prolog.AtomTypeAtom, address, env))
			}
		}, env)
	})
}
