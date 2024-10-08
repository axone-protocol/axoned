package keeper

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/axone-protocol/prolog"
	"github.com/axone-protocol/prolog/engine"
	"github.com/samber/lo"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v10/x/logic/fs/filtered"
	"github.com/axone-protocol/axoned/v10/x/logic/interpreter"
	"github.com/axone-protocol/axoned/v10/x/logic/interpreter/bootstrap"
	"github.com/axone-protocol/axoned/v10/x/logic/meter"
	prolog2 "github.com/axone-protocol/axoned/v10/x/logic/prolog"
	"github.com/axone-protocol/axoned/v10/x/logic/types"
	"github.com/axone-protocol/axoned/v10/x/logic/util"
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
	ctx context.Context, params types.Params, program, query string, solutionsLimit sdkmath.Uint,
) (*types.QueryServiceAskResponse, error) {
	ctx = k.enhanceContext(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	i, userOutput, err := k.newInterpreter(ctx, params)
	if err != nil {
		return nil, errorsmod.Wrapf(types.Internal, "error creating interpreter: %v", err.Error())
	}
	if err := i.ExecContext(ctx, program); err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error compiling query: %v", err.Error())
	}

	answer, err := k.queryInterpreter(ctx, i, query, solutionsLimit)
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
	ctx context.Context, i *prolog.Interpreter, query string, solutionsLimit sdkmath.Uint,
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
	if limits.MaxUserOutputSize != nil && limits.MaxUserOutputSize.GT(sdkmath.ZeroUint()) {
		userOutputBuffer = util.NewBoundedBufferMust(int(limits.MaxUserOutputSize.Uint64())) //nolint:gosec // disable G115
	} else {
		userOutputBuffer = new(strings.Builder)
	}

	options := []interpreter.Option{
		interpreter.WithHooks(
			whitelistBlacklistHookFn(whitelistPredicates, blacklistPredicates),
			gasMeterHookFn(sdkctx, params.GetGasPolicy()),
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

		predicateStringer, ok := operand.(fmt.Stringer)
		if !ok {
			return engine.SyntaxError(operand, env)
		}

		predicate := predicateStringer.String()

		if interpreter.IsRegistered(predicate) {
			if _, found := allowed.Get(predicate); !found {
				return engine.PermissionError(
					prolog2.AtomOperationExecute,
					prolog2.AtomPermissionForbiddenPredicate,
					engine.NewAtom(predicate),
					env,
				)
			}
		}
		return nil
	}
}

// gasMeterHookFn returns a hook function that consumes gas based on the cost of the executed predicate.
func gasMeterHookFn(ctx context.Context, gasPolicy types.GasPolicy) engine.HookFunc {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	gasMeter := meter.WithWeightedMeter(sdkctx.GasMeter(), nonNilNorZeroOrDefaultUint64(gasPolicy.WeightingFactor, defaultWeightFactor))

	return func(opcode engine.Opcode, operand engine.Term, env *engine.Env) (err error) {
		if opcode != engine.OpCall {
			return nil
		}

		operandStringer, ok := operand.(fmt.Stringer)
		if !ok {
			return engine.SyntaxError(operand, env)
		}

		predicate := operandStringer.String()

		cost := lookupCost(predicate,
			nonNilNorZeroOrDefaultUint64(gasPolicy.DefaultPredicateCost, defaultPredicateCost),
			gasPolicy.PredicateCosts)

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case storetypes.ErrorOutOfGas:
					err = errorsmod.Wrapf(
						types.LimitExceeded, "out of gas: %s <%s> (%d/%d)",
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
			return nonNilNorZeroOrDefaultUint64(c.Cost, defaultCost)
		}
	}

	return defaultCost
}

// nonNilNorZeroOrDefaultUint64 returns the value of the given sdkmath.Uint if it is not nil and not zero, otherwise it returns the
// given default value.
func nonNilNorZeroOrDefaultUint64(v *sdkmath.Uint, defaultValue uint64) uint64 {
	if v == nil || v.IsZero() {
		return defaultValue
	}

	return v.Uint64()
}
