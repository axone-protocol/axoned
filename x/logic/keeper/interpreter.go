package keeper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v8/x/logic/fs/filtered"
	"github.com/axone-protocol/axoned/v8/x/logic/interpreter"
	"github.com/axone-protocol/axoned/v8/x/logic/interpreter/bootstrap"
	"github.com/axone-protocol/axoned/v8/x/logic/meter"
	prolog2 "github.com/axone-protocol/axoned/v8/x/logic/prolog"
	"github.com/axone-protocol/axoned/v8/x/logic/types"
	"github.com/axone-protocol/axoned/v8/x/logic/util"
)

const (
	defaultEnvCap        = uint64(50)
	defaultPredicateCost = uint64(1)
	defaultWeightFactor  = uint64(1)
)

// writerStringer is an interface that combines io.Writer with capabilities of fmt.Stringer.
type writerStringer interface {
	io.Writer
	fmt.Stringer
}

func (k Keeper) enhanceContext(ctx context.Context) context.Context {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithValue(types.AuthKeeperContextKey, k.authKeeper)
	sdkCtx = sdkCtx.WithValue(types.BankKeeperContextKey, k.bankKeeper)
	return sdkCtx
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
		Height:     uint64(sdkCtx.BlockHeight()),
		GasUsed:    sdkCtx.GasMeter().GasConsumed(),
		Answer:     answer,
		UserOutput: userOutput.String(),
	}, nil
}

// queryInterpreter executes the given query on the given interpreter and returns the answer.
func (k Keeper) queryInterpreter(
	ctx context.Context, i *prolog.Interpreter, query string, solutionsLimit sdkmath.Uint,
) (*types.Answer, error) {
	p := engine.NewParser(&i.VM, strings.NewReader(query))
	t, err := p.Term()
	if err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error executing query: %v", err.Error())
	}

	var env *engine.Env
	count := sdkmath.ZeroUint()
	envs := make([]*engine.Env, 0, sdkmath.MinUint(solutionsLimit, sdkmath.NewUint(defaultEnvCap)).Uint64())
	_, callErr := engine.Call(&i.VM, t, func(env *engine.Env) *engine.Promise {
		if count.LT(solutionsLimit) {
			envs = append(envs, env)
		}
		count = count.Incr()
		return engine.Bool(count.GT(solutionsLimit))
	}, env).Force(ctx)

	vars := parsedVarsToVars(p.Vars)

	results, err := envsToResults(envs, p.Vars, i)
	if err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error executing query: %v", err.Error())
	}

	if callErr != nil {
		if sdkmath.NewUint(uint64(len(results))).LT(solutionsLimit) {
			// error is not part of the look-ahead and should be included in the solutions
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			var panicErr engine.PanicError
			switch {
			case errors.Is(callErr, types.LimitExceeded):
				return nil, callErr
			case errors.As(callErr, &panicErr) && errors.Is(panicErr.OriginErr, engine.ErrMaxVariables):
				return nil, errorsmod.Wrapf(types.LimitExceeded, panicErr.OriginErr.Error())
			case sdkCtx.GasMeter().IsOutOfGas():
				return nil, errorsmod.Wrapf(
					types.LimitExceeded, "out of gas: %s <%s> (%d/%d)",
					types.ModuleName, callErr.Error(), sdkCtx.GasMeter().GasConsumed(), sdkCtx.GasMeter().Limit())
			}
			results = append(results, types.Result{Error: callErr.Error()})
		} else {
			// error is part of the look-ahead, so let's consider that there's one more solution
			count = count.Incr()
		}
	}

	return &types.Answer{
		HasMore:   count.GT(solutionsLimit),
		Variables: vars,
		Results:   results,
	}, nil
}

// newInterpreter creates a new interpreter properly configured.
func (k Keeper) newInterpreter(ctx context.Context, params types.Params) (*prolog.Interpreter, fmt.Stringer, error) {
	sdkctx := sdk.UnwrapSDKContext(ctx)

	interpreterParams := params.GetInterpreter()
	gasPolicy := params.GetGasPolicy()
	limits := params.GetLimits()
	gasMeter := meter.WithWeightedMeter(sdkctx.GasMeter(), nonNilNorZeroOrDefaultUint64(gasPolicy.WeightingFactor, defaultWeightFactor))

	whitelistPredicates := util.NonZeroOrDefault(interpreterParams.PredicatesFilter.Whitelist, interpreter.RegistryNames)
	blacklistPredicates := interpreterParams.PredicatesFilter.Blacklist

	hook := func(predicate string) func(env *engine.Env) (err error) {
		return func(env *engine.Env) (err error) {
			if !util.WhitelistBlacklistMatches(whitelistPredicates, blacklistPredicates, prolog2.PredicateMatches)(predicate) {
				return engine.PermissionError(
					prolog2.AtomOperationExecute, prolog2.AtomPermissionForbiddenPredicate, engine.NewAtom(predicate), env)
			}
			cost := lookupCost(predicate, defaultPredicateCost, gasPolicy.PredicateCosts)

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
			return err
		}
	}

	whitelistUrls := lo.Map(
		util.NonZeroOrDefault(interpreterParams.VirtualFilesFilter.Whitelist, []string{}),
		util.Indexed(util.ParseURLMust))
	blacklistUrls := lo.Map(
		util.NonZeroOrDefault(interpreterParams.VirtualFilesFilter.Whitelist, []string{}),
		util.Indexed(util.ParseURLMust))

	var userOutputBuffer writerStringer
	if limits.MaxUserOutputSize != nil && limits.MaxUserOutputSize.GT(sdkmath.ZeroUint()) {
		userOutputBuffer = util.NewBoundedBufferMust(int(limits.MaxUserOutputSize.Uint64()))
	} else {
		userOutputBuffer = new(strings.Builder)
	}

	options := []interpreter.Option{
		interpreter.WithPredicates(ctx, interpreter.RegistryNames, hook),
		interpreter.WithBootstrap(ctx, util.NonZeroOrDefault(interpreterParams.GetBootstrap(), bootstrap.Bootstrap())),
		interpreter.WithFS(filtered.NewFS(k.fsProvider(ctx), whitelistUrls, blacklistUrls)),
		interpreter.WithUserOutputWriter(userOutputBuffer),
		interpreter.WithMaxVariables(limits.MaxVariables),
	}

	i, err := interpreter.New(options...)

	return i, userOutputBuffer, err
}

func lookupCost(predicate string, defaultCost uint64, costs []types.PredicateCost) uint64 {
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

func parsedVarsToVars(vars []engine.ParsedVariable) []string {
	return lo.Map(vars, func(v engine.ParsedVariable, _ int) string {
		return v.Name.String()
	})
}

func envsToResults(envs []*engine.Env, vars []engine.ParsedVariable, i *prolog.Interpreter) ([]types.Result, error) {
	results := make([]types.Result, 0, len(envs))
	for _, rEnv := range envs {
		substitutions := make([]types.Substitution, 0, len(vars))
		for _, v := range vars {
			if !isBound(v, rEnv) {
				// skip parsed variables that are not bound (singletons variables or other)
				continue
			}

			substitution, err := scanExpression(i, v, rEnv)
			if err != nil {
				return results, err
			}
			substitutions = append(substitutions, substitution)
		}
		results = append(results, types.Result{Substitutions: substitutions})
	}
	return results, nil
}

func scanExpression(i *prolog.Interpreter, v engine.ParsedVariable, rEnv *engine.Env) (types.Substitution, error) {
	var expression prolog.TermString
	err := expression.Scan(&i.VM, v.Variable, rEnv)
	if err != nil {
		return types.Substitution{}, err
	}
	substitution := types.Substitution{
		Variable:   v.Name.String(),
		Expression: string(expression),
	}

	return substitution, nil
}

// isBound returns true if the given parsed variable is bound in the given environment.
func isBound(v engine.ParsedVariable, env *engine.Env) bool {
	_, ok := env.Resolve(v.Variable).(engine.Variable)

	return !ok
}
