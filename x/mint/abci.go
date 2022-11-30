package mint

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/mint/keeper"
	"github.com/okp4/okp4d/x/mint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	totalSupply := k.TokenSupply(ctx, params.MintDenom)

	if uint64(ctx.BlockHeight()) == 1 {
		minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalSupply)
		minter.TargetSupply = totalSupply.Add(minter.AnnualProvisions.TruncateInt())
		k.SetMinter(ctx, minter)
	}

	// If we have reached the end of the year by reaching the targeted supply for the year
	// We need to re-calculate the next inflation for the next year.
	if totalSupply.GTE(minter.TargetSupply) {
		minter.Inflation = minter.NextInflation(params)
		minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalSupply)
		minter.TargetSupply = totalSupply.Add(minter.AnnualProvisions.TruncateInt())
		k.SetMinter(ctx, minter)
	}

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params, totalSupply)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
