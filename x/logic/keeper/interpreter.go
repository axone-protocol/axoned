package keeper

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/axone-protocol/prolog/v2"
	"github.com/axone-protocol/prolog/v2/engine"
	"github.com/samber/lo"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v12/x/logic/fs/filtered"
	"github.com/axone-protocol/axoned/v12/x/logic/interpreter"
	"github.com/axone-protocol/axoned/v12/x/logic/interpreter/bootstrap"
	"github.com/axone-protocol/axoned/v12/x/logic/meter"
	prolog2 "github.com/axone-protocol/axoned/v12/x/logic/prolog"
	"github.com/axone-protocol/axoned/v12/x/logic/types"
	"github.com/axone-protocol/axoned/v12/x/logic/util"
)

const (
	defaultPredicateCost = uint64(1)
	defaultWeightFactor  = uint64(1)
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
) (*types.QueryServiceAskResponse, error) {
	ctx = k.enhanceContext(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	i, userOutput, err := k.newInterpreter(ctx, params)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "error creating interpreter: %v", err.Error())
	}
	if err := i.ExecContext(ctx, program); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "error compiling query: %v", err.Error())
	}

	answer, err := k.queryInterpreter(ctx, i, query, calculateSolutionLimit(solutionsLimit, params.GetLimits().MaxResultCount))
	if err != nil {
		return nil, err
	}

	return &types.QueryServiceAskResponse{
		Height:     uint64(sdkCtx.BlockHeight()), //nolint:gosec // disable G115
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
	sdkctx := sdk.UnwrapSDKContext(ctx)

	interpreterParams := params.GetInterpreter()
	whitelistPredicates := util.NonZeroOrDefault(interpreterParams.PredicatesFilter.Whitelist, interpreter.RegistryNames)
	blacklistPredicates := interpreterParams.PredicatesFilter.Blacklist

	whitelistUrls := lo.Map(
		util.NonZeroOrDefault(interpreterParams.VirtualFilesFilter.Whitelist, []string{}),
		util.Indexed(util.ParseURLMust))
	blacklistUrls := lo.Map(
		util.NonZeroOrDefault(interpreterParams.VirtualFilesFilter.Whitelist, []string{}),
		util.Indexed(util.ParseURLMust))

	var userOutputBuffer writerStringer
	limits := params.GetLimits()
	if limits.MaxUserOutputSize > 0 {
		userOutputBuffer = util.NewBoundedBufferMust(int(limits.MaxUserOutputSize)) //nolint:gosec // disable G115
	} else {
		userOutputBuffer = new(strings.Builder)
	}

	options := []interpreter.Option{
		interpreter.WithHooks(
			whitelistBlacklistHookFn(whitelistPredicates, blacklistPredicates),
			gasMeterHookFn(sdkctx, params.GetGasPolicy()),
			telemetryPredicateCallCounterHookFn(),
			telemetryPredicateDurationHookFn(),
		),
		interpreter.WithPredicates(ctx, interpreter.RegistryNames),
		interpreter.WithBootstrap(ctx, util.NonZeroOrDefault(interpreterParams.GetBootstrap(), bootstrap.Bootstrap())),
		interpreter.WithFS(filtered.NewFS(k.fsProvider(ctx), whitelistUrls, blacklistUrls)),
		interpreter.WithUserOutputWriter(userOutputBuffer),
		interpreter.WithMaxVariables(limits.MaxVariables),
	}

	i, err := interpreter.New(options...)

	return i, userOutputBuffer, err
}

// whitelistBlacklistHookFn returns a hook function that checks if the given predicate is allowed to be executed.
// The predicate is allowed if it is in the whitelist or not in the blacklist.
func whitelistBlacklistHookFn(whitelist, blacklist []string) engine.HookFunc {
	allowed := lo.Reduce(
		lo.Filter(interpreter.RegistryNames,
			util.Indexed(util.WhitelistBlacklistMatches(whitelist, blacklist, prolog2.PredicateMatches))),
		func(agg *orderedmap.OrderedMap[string, struct{}], item string, _ int) *orderedmap.OrderedMap[string, struct{}] {
			agg.Set(item, struct{}{})
			return agg
		},
		orderedmap.New[string, struct{}](orderedmap.WithCapacity[string, struct{}](len(interpreter.RegistryNames))))

	return func(opcode engine.Opcode, operand engine.Term, env *engine.Env) error {
		if opcode != engine.OpCall {
			return nil
		}

		predicate, ok := stringifyOperand(operand)
		if !ok {
			return engine.SyntaxError(operand, env)
		}

		if !interpreter.IsRegistered(predicate) {
			return nil
		}

		if _, found := allowed.Get(predicate); !found {
			return engine.PermissionError(
				prolog2.AtomOperationExecute,
				prolog2.AtomPermissionForbiddenPredicate,
				engine.NewAtom(predicate),
				env,
			)
		}

		return nil
	}
}

// gasMeterHookFn returns a hook function that consumes gas based on the cost of the executed predicate.
func gasMeterHookFn(ctx context.Context, gasPolicy types.GasPolicy) engine.HookFunc {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	gasMeter := meter.WithWeightedMeter(sdkctx.GasMeter(), lo.CoalesceOrEmpty(gasPolicy.WeightingFactor, defaultWeightFactor))

	return func(opcode engine.Opcode, operand engine.Term, env *engine.Env) (err error) {
		if opcode != engine.OpCall {
			return nil
		}

		predicate, ok := stringifyOperand(operand)
		if !ok {
			return engine.SyntaxError(operand, env)
		}

		cost := lookupCost(predicate,
			lo.CoalesceOrEmpty(gasPolicy.DefaultPredicateCost, defaultPredicateCost),
			gasPolicy.PredicateCosts)

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case storetypes.ErrorOutOfGas:
					err = errorsmod.Wrapf(
						types.ErrLimitExceeded, "out of gas: %s <%s> (%d/%d)",
						types.ModuleName, rType.Descriptor, sdkctx.GasMeter().GasConsumed(), sdkctx.GasMeter().Limit())
				default:
					panic(r)
				}
			}
		}()
		gasMeter.ConsumeGas(cost, predicate)

		return nil
	}
}

func lookupCost(predicate string, defaultCost uint64, costs []types.PredicateCost) uint64 {
	if !interpreter.IsRegistered(predicate) {
		return defaultCost
	}

	for _, c := range costs {
		if prolog2.PredicateMatches(predicate)(c.Predicate) {
			return lo.CoalesceOrEmpty(c.Cost, defaultCost)
		}
	}

	return defaultCost
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
