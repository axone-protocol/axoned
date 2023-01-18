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
//   - Coins represents the balances of the account as a list of pairs of coin denomination and amount.
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
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}
		bankKeeper := sdkContext.Value(types.BankKeeperContextKey).(types.BankKeeper)

		bech32Addr := sdk.AccAddress(nil)
		switch acc := env.Resolve(account).(type) {
		case engine.Variable:
		case engine.Atom:
			bech32Addr, err = sdk.AccAddressFromBech32(acc.String())
			if err != nil {
				return engine.Error(fmt.Errorf("bank_balances/2: %w", err))
			}
		default:
			return engine.Error(fmt.Errorf("bank_balances/2: cannot unify account address with %T", acc))
		}

		if bech32Addr != nil {
			fetchedBalances := AllBalancesSorted(sdkContext, bankKeeper, bech32Addr)

			return engine.Unify(vm, CoinsToTerm(fetchedBalances), balances, cont, env)
		}

		allBalances := bankKeeper.GetAccountsBalances(sdkContext)
		promises := make([]func(ctx context.Context) *engine.Promise, 0, len(allBalances))
		for _, balance := range allBalances {
			address, coins := balance.Address, balance.Coins
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

func BankSpendableCoins(vm *engine.VM, account, balances engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}
		bankKeeper := sdkContext.Value(types.BankKeeperContextKey).(types.BankKeeper)

		bech32Addr := sdk.AccAddress(nil)
		switch acc := env.Resolve(account).(type) {
		case engine.Atom:
			bech32Addr, err = sdk.AccAddressFromBech32(acc.String())
			if err != nil {
				return engine.Error(fmt.Errorf("bank_spendable_coins/2: %w", err))
			}

			fetchedBalances := SpendableCoinsSorted(sdkContext, bankKeeper, bech32Addr)

			return engine.Unify(vm, CoinsToTerm(fetchedBalances), balances, cont, env)
		default:
			return engine.Error(fmt.Errorf("bank_spendable_coins/2: cannot unify account address with %T", acc))
		}

	})
}
