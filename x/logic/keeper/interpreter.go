package keeper

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/axone-protocol/prolog/v3"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/interpreter"
	"github.com/axone-protocol/axoned/v15/x/logic/interpreter/bootstrap"
	"github.com/axone-protocol/axoned/v15/x/logic/meter"
	prolog2 "github.com/axone-protocol/axoned/v15/x/logic/prolog"
	"github.com/axone-protocol/axoned/v15/x/logic/types"
	"github.com/axone-protocol/axoned/v15/x/logic/util"
)

// writerStringer is an interface that combines io.Writer with capabilities of fmt.Stringer.
type writerStringer interface {
	io.Writer
	fmt.Stringer
}

func (k Keeper) enhanceContext(ctx context.Context) context.Context {
	return sdk.UnwrapSDKContext(ctx).
		WithValue(types.InterfaceRegistryContextKey, k.interfaceRegistry).
		WithValue(types.AuthKeeperContextKey, k.authKeeper).
		WithValue(types.AuthQueryServiceContextKey, k.authQueryService).
		WithValue(types.BankKeeperContextKey, k.bankKeeper)
}

func (k Keeper) execute(
	ctx context.Context, params types.Params, program, query string, solutionsLimit uint64,
) (*types.QueryAskResponse, error) {
	ctx = k.enhanceContext(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	i, userOutput, err := k.newInterpreter(ctx, params)
	if err != nil {
		if limitErr := util.AsLimitExceededError(ctx, err); limitErr != nil {
			return nil, limitErr
		}
		return nil, errorsmod.Wrapf(types.ErrInternal, "error creating interpreter: %v", err.Error())
	}
	if err := i.ExecContext(ctx, program); err != nil {
		if limitErr := util.AsLimitExceededError(ctx, err); limitErr != nil {
			return nil, limitErr
		}
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "error compiling query: %v", err.Error())
	}

	answer, err := k.queryInterpreter(ctx, i, query, calculateSolutionLimit(solutionsLimit, params.GetLimits().MaxResultCount))
	if err != nil {
		return nil, err
	}

	return &types.QueryAskResponse{
		Height:     uint64(prolog2.ResolveHeaderInfo(sdkCtx).Height), //nolint:gosec // disable G115
		GasUsed:    sdkCtx.GasMeter().GasConsumed(),
		Answer:     answer,
		UserOutput: userOutput.String(),
	}, nil
}

// queryInterpreter executes the given query on the given interpreter and returns the answer.
func (k Keeper) queryInterpreter(
	ctx context.Context, i *prolog.Interpreter, query string, solutionsLimit uint64,
) (*types.Answer, error) {
	return util.QueryInterpreter(ctx, i, query, solutionsLimit)
}

// newInterpreter creates a new interpreter properly configured.
func (k Keeper) newInterpreter(ctx context.Context, params types.Params) (*prolog.Interpreter, fmt.Stringer, error) {
	sdkctx := sdk.UnwrapSDKContext(ctx).WithValue(types.IOCoeffContextKey, params.GetGasPolicy().IoCoeff)
	ctx = sdkctx

	var userOutputBuffer writerStringer
	limits := params.GetLimits()
	if limits.MaxUserOutputSize > 0 {
		userOutputBuffer = util.NewBoundedBufferMust(int(limits.MaxUserOutputSize)) //nolint:gosec // disable G115
	} else {
		userOutputBuffer = new(strings.Builder)
	}

	fsProvider, err := k.fsProvider(ctx)
	if err != nil {
		return nil, nil, errorsmod.Wrapf(types.ErrInternal, "error getting filesystem: %v", err.Error())
	}

	options := []interpreter.Option{
		interpreter.WithHooks(telemetryPredicateCallCounterHookFn(), telemetryPredicateDurationHookFn()),
		interpreter.WithPredicates(ctx, interpreter.RegistryNames),
		// Bootstrap is part of the kernel space and must not affect user-space gas accounting.
		interpreter.WithBootstrap(ctx, bootstrap.Bootstrap()),
		interpreter.WithMeter(
			meter.NewVMMeter(
				sdkctx.GasMeter(),
				params.GetGasPolicy().ComputeCoeff,
				params.GetGasPolicy().MemoryCoeff,
				params.GetGasPolicy().UnifyCoeff,
			),
		),
		interpreter.WithFS(fsProvider),
		interpreter.WithUserOutputWriter(userOutputBuffer),
		interpreter.WithMaxVariables(limits.MaxVariables),
	}

	i, err := interpreter.New(options...)

	return i, userOutputBuffer, err
}

// calculateSolutionLimit returns the final number of solutions to be returned based on the given number of solutions and the
// maximum result count.
func calculateSolutionLimit(nbSolutions uint64, maxResultCount uint64) uint64 {
	nbSolutions = max(nbSolutions, 1)
	if maxResultCount == 0 {
		return nbSolutions
	}

	return min(nbSolutions, maxResultCount)
}
