package keeper

import (
	goctx "context"
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog"
	"github.com/okp4/okp4d/x/logic/context"
	"github.com/okp4/okp4d/x/logic/interpreter"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"
)

// withLimitContext returns a context with the limits configured for the module.
func (k Keeper) withLimitContext(ctx goctx.Context) (goctx.Context, context.IncrementCountByFunc) {
	limits := k.getLimits(ctx)

	maxGas := util.DerefOrDefault(limits.MaxGas, sdkmath.NewUint(math.MaxInt64))

	return context.WithLimit(ctx, maxGas.Uint64())
}

func (k Keeper) getLimits(ctx goctx.Context) types.Limits {
	params := k.GetParams(sdk.UnwrapSDKContext(ctx))
	return params.GetLimits()
}

func (k Keeper) execute(goctx goctx.Context, program, query string) (*types.QueryServiceAskResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(goctx)

	i, limitContext, err := k.newInterpreter(goctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.Internal, "error creating interpreter: %v", err.Error())
	}

	if err := i.ExecContext(limitContext, program); err != nil {
		return nil, sdkerrors.Wrapf(types.InvalidArgument, "error compiling query: %v", err.Error())
	}

	sols, err := i.QueryContext(limitContext, query)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.InvalidArgument, "error executing query: %v", err.Error())
	}
	defer func() {
		_ = sols.Close()
	}()

	success := false
	limits := k.getLimits(sdkCtx)
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

// newInterpreter creates a new interpreter with the limits configured for the module, and initialized with the
// interpreter's module settings.
func (k Keeper) newInterpreter(ctx goctx.Context) (*prolog.Interpreter, goctx.Context, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(ctx))
	interpreterParams := params.GetInterpreter()
	limitContext, inc := k.withLimitContext(ctx)

	interpreted, err := interpreter.NewInstrumentedInterpreter(
		limitContext,
		util.NonZeroOrDefault(interpreterParams.GetRegisteredPredicates(), interpreter.RegistryNames),
		util.NonZeroOrDefault(interpreterParams.GetBootstrap(), interpreter.Bootstrap()),
		inc)

	return interpreted, limitContext, err
}
