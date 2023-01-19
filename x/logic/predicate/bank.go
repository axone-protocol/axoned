package predicate

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"
)

// BankBalances is a predicate which unifies the given terms with the list of balances (coins) of the given account.
//
//	bank_balances(?Account, ?Balances)
//
// where:
//   - Account represents the account address (in Bech32 format).
//   - Balances represents the balances of the account as a list of pairs of coin denomination and amount.
//
// Example:
//
//	# Query the balances of the account.
//	- bank_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).
//
// # Query the balances of all accounts. The result is a list of pairs of account address and balances.
// - bank_balances(X, Y).
//
// # Query the first balance of the given account by unifying the denomination and amount with the given terms.
// - bank_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A), _]).
func BankBalances(vm *engine.VM, account, balances engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return fetchBalances("bank_balances/2", account, balances, vm, env, cont, func(ctx sdk.Context, bankKeeper types.BankKeeper, address sdk.AccAddress) sdk.Coins {
		return AllBalancesSorted(ctx, bankKeeper, address)
	})
}

// BankSpendableCoins is a predicate which unifies the given terms with the list of spendable coins of the given account.
//
//	bank_spendable_coins(?Account, ?Balances)
//
// where:
//   - Account represents the account address (in Bech32 format).
//   - Balances represents the spendable coins of the account as a list of pairs of coin denomination and amount.
//
// Example:
//
//	# Query the spendable coins of the account.
//	- bank_spendable_coins('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).
//
// # Query the spendable coins of all accounts. The result is a list of pairs of account address and balances.
// - bank_spendable_coins(X, Y).
//
// # Query the first spendable coin of the given account by unifying the denomination and amount with the given terms.
// - bank_spendable_coins('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A), _]).
func BankSpendableCoins(vm *engine.VM, account, balances engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return fetchBalances("bank_spendable_coins/2", account, balances, vm, env, cont, func(ctx sdk.Context, bankKeeper types.BankKeeper, address sdk.AccAddress) sdk.Coins {
		return SpendableCoinsSorted(ctx, bankKeeper, address)
	})
}

func BankLockedCoins(vm *engine.VM, account, balances engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return fetchBalances("bank_locked_coins/2", account, balances, vm, env, cont, func(ctx sdk.Context, bankKeeper types.BankKeeper, address sdk.AccAddress) sdk.Coins {
		return LockedCoinsSorted(ctx, bankKeeper, address)
	})
}

func getBech32(env *engine.Env, account engine.Term) (sdk.AccAddress, error) {
	switch acc := env.Resolve(account).(type) {
	case engine.Variable:
	case engine.Atom:
		return sdk.AccAddressFromBech32(acc.String())
	default:
		return nil, fmt.Errorf("cannot unify account address with %T", acc)
	}
	return sdk.AccAddress(nil), nil
}

func fetchBalances(
	predicate string,
	account, balances engine.Term,
	vm *engine.VM,
	env *engine.Env,
	cont engine.Cont,
	coinsFn func(ctx sdk.Context, bankKeeper types.BankKeeper, address sdk.AccAddress) sdk.Coins) *engine.Promise {

	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}
		bankKeeper := sdkContext.Value(types.BankKeeperContextKey).(types.BankKeeper)

		bech32Addr, err := getBech32(env, account)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: %w", predicate, err))
		}

		if bech32Addr != nil {
			fetchedBalances := coinsFn(sdkContext, bankKeeper, bech32Addr)
			return engine.Unify(vm, CoinsToTerm(fetchedBalances), balances, cont, env)
		}

		allBalances := bankKeeper.GetAccountsBalances(sdkContext)
		promises := make([]func(ctx context.Context) *engine.Promise, 0, len(allBalances))
		for _, balance := range allBalances {
			address := balance.Address
			bech32Addr, err = sdk.AccAddressFromBech32(address)
			if err != nil {
				return engine.Error(fmt.Errorf("%s: %w", predicate, err))
			}
			coins := coinsFn(sdkContext, bankKeeper, bech32Addr)

			promises = append(
				promises,
				func(ctx context.Context) *engine.Promise {
					return engine.Unify(
						vm,
						Tuple(engine.NewAtom(address), CoinsToTerm(coins)),
						Tuple(account, balances),
						cont,
						env,
					)
				})
		}
		return engine.Delay(promises...)
	})
}
