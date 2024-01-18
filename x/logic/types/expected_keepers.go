package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias).
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetAccountsBalances(ctx context.Context) []bank.Balance
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	LockedCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

// WasmKeeper defines the expected interface needed to request smart contracts.
type WasmKeeper interface {
	QuerySmart(ctx context.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
}
