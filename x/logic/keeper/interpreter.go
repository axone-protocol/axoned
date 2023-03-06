package keeper

import (
	goctx "context"
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog"
	"github.com/okp4/okp4d/x/logic/interpreter"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"
)

func (k Keeper) limits(ctx goctx.Context) types.Limits {
	params := k.GetParams(sdk.UnwrapSDKContext(ctx))
	return params.GetLimits()
}

func (k Keeper) enhanceContext(ctx goctx.Context) goctx.Context {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithValue(types.AuthKeeperContextKey, k.authKeeper)
	sdkCtx = sdkCtx.WithValue(types.BankKeeperContextKey, k.bankKeeper)

	sdkCtx = sdkCtx.WithValue(types.WasmKeeperContextKey, k.WasmKeeper)

	return sdkCtx
}

func (k Keeper) execute(ctx goctx.Context, program, query string) (*types.QueryServiceAskResponse, error) {
	ctx = k.enhanceContext(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	i, err := k.newInterpreter(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.Internal, "error creating interpreter: %v", err.Error())
	}

	if err := i.ExecContext(ctx, program); err != nil {
		return nil, sdkerrors.Wrapf(types.InvalidArgument, "error compiling query: %v", err.Error())
	}

	sols, err := i.QueryContext(ctx, query)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.InvalidArgument, "error executing query: %v", err.Error())
	}
	defer func() {
		_ = sols.Close()
	}()

	success := false
	limits := k.limits(ctx)
	var variables []string
	results := make([]types.Result, 0)
	for nb := sdkmath.ZeroUint(); nb.LT(*limits.MaxResultCount) && sols.Next(); nb = nb.Incr() {
		success = true

		m := types.TermResults{}
		if err := sols.Scan(m); err != nil {
			return nil, sdkerrors.Wrapf(types.Internal, "error scanning solution: %v", err.Error())
		}

		if nb.IsZero() {
			variables = m.ToVariables()
		}

		results = append(results, types.Result{Substitutions: m.ToSubstitutions()})
	}

	return &types.QueryServiceAskResponse{
		Height:  uint64(sdkCtx.BlockHeight()),
		GasUsed: sdkCtx.GasMeter().GasConsumed(),
		Answer: &types.Answer{
			Success:   success,
			HasMore:   sols.Next(),
			Variables: variables,
			Results:   results,
		},
	}, nil
}

func checkLimits(request *types.QueryServiceAskRequest, limits types.Limits) error {
	size := sdkmath.NewUint(uint64(len(request.GetQuery())))
	limit := util.DerefOrDefault(limits.MaxSize, sdkmath.NewUint(math.MaxInt64))
	if size.GT(limit) {
		return sdkerrors.Wrapf(types.LimitExceeded, "query: %d > MaxSize: %d", size, limit)
	}

	return nil
}

// newInterpreter creates a new interpreter properly configured.
func (k Keeper) newInterpreter(ctx goctx.Context) (*prolog.Interpreter, error) {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkctx)

	interpreterParams := params.GetInterpreter()

	interpreted, err := interpreter.New(
		ctx,
		util.NonZeroOrDefault(interpreterParams.GetRegisteredPredicates(), interpreter.RegistryNames),
		util.NonZeroOrDefault(interpreterParams.GetBootstrap(), interpreter.Bootstrap()),
		sdkctx.GasMeter(),
		k.WasmKeeper,
	)

	return interpreted, err
}
