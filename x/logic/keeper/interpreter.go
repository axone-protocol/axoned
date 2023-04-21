package keeper

import (
	goctx "context"
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog"
	"github.com/okp4/okp4d/x/logic/fs"
	"github.com/okp4/okp4d/x/logic/interpreter"
	"github.com/okp4/okp4d/x/logic/interpreter/bootstrap"
	"github.com/okp4/okp4d/x/logic/meter"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"
	"github.com/samber/lo"
)

const (
	defaultPredicateCost = uint64(1)
	defaultWeightFactor  = uint64(1)
)

func (k Keeper) limits(ctx goctx.Context) types.Limits {
	params := k.GetParams(sdk.UnwrapSDKContext(ctx))
	return params.GetLimits()
}

func (k Keeper) enhanceContext(ctx goctx.Context) goctx.Context {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithValue(types.AuthKeeperContextKey, k.authKeeper)
	sdkCtx = sdkCtx.WithValue(types.BankKeeperContextKey, k.bankKeeper)
	return sdkCtx
}

func (k Keeper) execute(ctx goctx.Context, program, query string) (*types.QueryServiceAskResponse, error) {
	ctx = k.enhanceContext(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	i, userOutputBuffer, err := k.newInterpreter(ctx)
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

	if sols.Err() != nil && sdkCtx.GasMeter().IsOutOfGas() {
		panic(sdk.ErrorOutOfGas{Descriptor: "Prolog interpreter execution"})
	}

	var userOutput string
	if userOutputBuffer != nil {
		userOutput = userOutputBuffer.String()
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
		UserOutput: userOutput,
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
func (k Keeper) newInterpreter(ctx goctx.Context) (*prolog.Interpreter, *util.BoundedBuffer, error) {
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
				util.Indexed(util.WhitelistBlacklistMatches(whitelistPredicates, blacklistPredicates, util.PredicateMatches))),
			toPredicate(
				nonNilNorZeroOrDefaultUint64(gasPolicy.DefaultPredicateCost, defaultPredicateCost),
				gasPolicy.GetPredicateCosts())),
		func(agg interpreter.Predicates, item lo.Tuple2[string, uint64], index int) interpreter.Predicates {
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

// toPredicate converts the given predicate costs to a function that returns the cost for the given predicate as
// a pair of predicate name and cost.
func toPredicate(defaultCost uint64, predicateCosts []types.PredicateCost) func(string, int) lo.Tuple2[string, uint64] {
	return func(predicate string, _ int) lo.Tuple2[string, uint64] {
		for _, c := range predicateCosts {
			if util.PredicateMatches(predicate)(c.Predicate) {
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
