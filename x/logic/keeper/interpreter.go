package keeper

import (
	"context"
	"math"
	"strings"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/logic/fs"
	"github.com/okp4/okp4d/x/logic/interpreter"
	"github.com/okp4/okp4d/x/logic/interpreter/bootstrap"
	"github.com/okp4/okp4d/x/logic/meter"
	prolog2 "github.com/okp4/okp4d/x/logic/prolog"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"
)

const (
	defaultPredicateCost = uint64(1)
	defaultWeightFactor  = uint64(1)
)

func (k Keeper) limits(ctx context.Context) types.Limits {
	params := k.GetParams(sdk.UnwrapSDKContext(ctx))
	return params.GetLimits()
}

func (k Keeper) enhanceContext(ctx context.Context) context.Context {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithValue(types.AuthKeeperContextKey, k.authKeeper)
	sdkCtx = sdkCtx.WithValue(types.BankKeeperContextKey, k.bankKeeper)
	return sdkCtx
}

func (k Keeper) execute(ctx context.Context, program, query string) (*types.QueryServiceAskResponse, error) {
	ctx = k.enhanceContext(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	limits := k.limits(sdkCtx)

	i, userOutputBuffer, err := k.newInterpreter(ctx)
	if err != nil {
		return nil, errorsmod.Wrapf(types.Internal, "error creating interpreter: %v", err.Error())
	}
	if err := i.ExecContext(ctx, program); err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error compiling query: %v", err.Error())
	}

	answer, err := k.queryInterpreter(ctx, i, query, *limits.MaxResultCount)
	if err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error executing query: %v", err.Error())
	}

	var userOutput string
	if userOutputBuffer != nil {
		userOutput = userOutputBuffer.String()
	}

	return &types.QueryServiceAskResponse{
		Height:     uint64(sdkCtx.BlockHeight()),
		GasUsed:    sdkCtx.GasMeter().GasConsumed(),
		Answer:     answer,
		UserOutput: userOutput,
	}, nil
}

// queryInterpreter executes the given query on the given interpreter and returns the answer.
func (k Keeper) queryInterpreter(ctx context.Context, i *prolog.Interpreter, query string, limit sdkmath.Uint) (*types.Answer, error) {
	p := engine.NewParser(&i.VM, strings.NewReader(query))
	t, err := p.Term()
	if err != nil {
		return nil, err
	}

	var env *engine.Env
	count := sdkmath.ZeroUint()
	envs := make([]*engine.Env, 0, limit.Uint64())
	_, callErr := engine.Call(&i.VM, t, func(env *engine.Env) *engine.Promise {
		if count.LT(limit) {
			envs = append(envs, env)
		}
		count = count.Incr()
		return engine.Bool(count.GT(limit))
	}, env).Force(ctx)

	answerErr := lo.IfF(callErr != nil, func() string {
		return callErr.Error()
	}).Else("")
	success := len(envs) > 0
	hasMore := count.GT(limit)
	vars := parsedVarsToVars(p.Vars)
	results, err := envsToResults(envs, p.Vars, i)
	if err != nil {
		return nil, err
	}

	return &types.Answer{
		Success:   success,
		Error:     answerErr,
		HasMore:   hasMore,
		Variables: vars,
		Results:   results,
	}, nil
}

// newInterpreter creates a new interpreter properly configured.
func (k Keeper) newInterpreter(ctx context.Context) (*prolog.Interpreter, *util.BoundedBuffer, error) {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkctx)

	interpreterParams := params.GetInterpreter()
	gasPolicy := params.GetGasPolicy()
	limits := params.GetLimits()
	gasMeter := meter.WithWeightedMeter(sdkctx.GasMeter(), nonNilNorZeroOrDefaultUint64(gasPolicy.WeightingFactor, defaultWeightFactor))

	whitelistPredicates := util.NonZeroOrDefault(interpreterParams.PredicatesFilter.Whitelist, interpreter.RegistryNames)
	blacklistPredicates := interpreterParams.PredicatesFilter.Blacklist
	predicates := lo.Reduce(
		lo.Map(
			lo.Filter(
				interpreter.RegistryNames,
				util.Indexed(util.WhitelistBlacklistMatches(whitelistPredicates, blacklistPredicates, prolog2.PredicateMatches))),
			toPredicate(
				nonNilNorZeroOrDefaultUint64(gasPolicy.DefaultPredicateCost, defaultPredicateCost),
				gasPolicy.GetPredicateCosts())),
		func(agg interpreter.Predicates, item lo.Tuple2[string, uint64], _ int) interpreter.Predicates {
			agg[item.A] = item.B
			return agg
		},
		interpreter.Predicates{})

	whitelistUrls := lo.Map(
		util.NonZeroOrDefault(interpreterParams.VirtualFilesFilter.Whitelist, []string{}),
		util.Indexed(util.ParseURLMust))
	blacklistUrls := lo.Map(
		util.NonZeroOrDefault(interpreterParams.VirtualFilesFilter.Whitelist, []string{}),
		util.Indexed(util.ParseURLMust))

	options := []interpreter.Option{
		interpreter.WithPredicates(ctx, predicates, gasMeter),
		interpreter.WithBootstrap(ctx, util.NonZeroOrDefault(interpreterParams.GetBootstrap(), bootstrap.Bootstrap())),
		interpreter.WithFS(fs.NewFilteredFS(whitelistUrls, blacklistUrls, k.fsProvider(ctx))),
	}

	var userOutputBuffer *util.BoundedBuffer
	if limits.MaxUserOutputSize != nil && limits.MaxUserOutputSize.GT(sdkmath.ZeroUint()) {
		userOutputBuffer = util.NewBoundedBufferMust(int(limits.MaxUserOutputSize.Uint64()))
		options = append(options, interpreter.WithUserOutputWriter(userOutputBuffer))
	}

	i, err := interpreter.New(options...)

	return i, userOutputBuffer, err
}

func checkLimits(request *types.QueryServiceAskRequest, limits types.Limits) error {
	size := sdkmath.NewUint(uint64(len(request.GetQuery())))
	limit := util.DerefOrDefault(limits.MaxSize, sdkmath.NewUint(math.MaxInt64))
	if size.GT(limit) {
		return errorsmod.Wrapf(types.LimitExceeded, "query: %d > MaxSize: %d", size, limit)
	}

	return nil
}

// toPredicate converts the given predicate costs to a function that returns the cost for the given predicate as
// a pair of predicate name and cost.
func toPredicate(defaultCost uint64, predicateCosts []types.PredicateCost) func(string, int) lo.Tuple2[string, uint64] {
	return func(predicate string, _ int) lo.Tuple2[string, uint64] {
		for _, c := range predicateCosts {
			if prolog2.PredicateMatches(predicate)(c.Predicate) {
				return lo.T2(predicate, nonNilNorZeroOrDefaultUint64(c.Cost, defaultCost))
			}
		}

		return lo.T2(predicate, defaultCost)
	}
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
			var expression prolog.TermString
			err := expression.Scan(&i.VM, v.Variable, rEnv)
			if err != nil {
				return nil, err
			}
			substitution := types.Substitution{
				Variable:   v.Name.String(),
				Expression: string(expression),
			}
			substitutions = append(substitutions, substitution)
		}
		results = append(results, types.Result{Substitutions: substitutions})
	}
	return results, nil
}
