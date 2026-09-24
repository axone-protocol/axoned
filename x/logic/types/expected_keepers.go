package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias).
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

type AuthQueryService interface {
	Accounts(ctx context.Context, req *auth.QueryAccountsRequest) (*auth.QueryAccountsResponse, error)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	LockedCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

// StakingQueryService defines the staking queries exposed to logic programs.
type StakingQueryService interface {
	DelegatorDelegations(
		ctx context.Context, req *staking.QueryDelegatorDelegationsRequest,
	) (*staking.QueryDelegatorDelegationsResponse, error)
	DelegatorUnbondingDelegations(
		ctx context.Context, req *staking.QueryDelegatorUnbondingDelegationsRequest,
	) (*staking.QueryDelegatorUnbondingDelegationsResponse, error)
	Redelegations(ctx context.Context, req *staking.QueryRedelegationsRequest) (*staking.QueryRedelegationsResponse, error)
	Params(ctx context.Context, req *staking.QueryParamsRequest) (*staking.QueryParamsResponse, error)
}

// WasmKeeper defines the expected interface needed to request smart contracts.
type WasmKeeper interface {
	QuerySmart(ctx context.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
}
