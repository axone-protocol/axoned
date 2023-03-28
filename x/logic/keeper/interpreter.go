package keeper

import (
	goctx "context"
	"math"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog"
	"github.com/okp4/okp4d/x/logic/interpreter"
	"github.com/okp4/okp4d/x/logic/interpreter/bootstrap"
	"github.com/okp4/okp4d/x/logic/meter"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"
	"github.com/samber/lo"
)

const defaultCost = uint64(1)

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
	gasPolicy := params.GetGasPolicy()

	whitelist := util.NonZeroOrDefault(interpreterParams.PredicatesWhitelist, interpreter.RegistryNames)
	blacklist := interpreterParams.GetPredicatesBlacklist()
	gasMeter := meter.WithWeightedMeter(sdkctx.GasMeter(), gasPolicy.WeightingFactor.Uint64())

	predicates := lo.Reduce(
		lo.Map(
			lo.Filter(interpreter.RegistryNames, filterPredicates(whitelist, blacklist)),
			toPredicate(gasPolicy.GetPredicateCosts())),
		func(agg interpreter.Predicates, item lo.Tuple2[string, uint64], index int) interpreter.Predicates {
			agg[item.A] = item.B
			return agg
		}, interpreter.Predicates{})

	interpreted, err := interpreter.New(
		ctx,
		predicates,
		util.NonZeroOrDefault(interpreterParams.GetBootstrap(), bootstrap.Bootstrap()),
		gasMeter,
		k.fsProvider(ctx),
	)

	return interpreted, err
}

// filterPredicates filters the given predicate (with arity) according to the given whitelist and blacklist.
// The whitelist and blacklist are applied to the registry to determine the final predicate list.
// The whitelist and blacklist can contain predicates with or without arity, e.g. "foo/0", "foo", "bar/1".
func filterPredicates(whitelist []string, blacklist []string) func(string, int) bool {
	return func(predicate string, _ int) bool {
		return lo.ContainsBy(whitelist, matchPredicate(predicate)) && !lo.ContainsBy(blacklist, matchPredicate(predicate))
	}
}

// matchPredicate returns a function that matches the given predicate against the given other predicate.
// If the other predicate contains a slash, it is matched as is. Otherwise, the other predicate is matched against the
// first part of the given predicate.
// For example:
//   - matchPredicate("foo/0")("foo/0") -> true
//   - matchPredicate("foo/0")("foo/1") -> false
//   - matchPredicate("foo/0")("foo") -> true
//   - matchPredicate("foo/0")("bar") -> false
func matchPredicate(predicate string) func(b string) bool {
	return func(other string) bool {
		if strings.Contains(other, "/") {
			return predicate == other
		}
		return strings.Split(predicate, "/")[0] == other
	}
}

// toPredicate converts the given predicate costs to a function that returns the cost for the given predicate as
// a pair of predicate name and cost.
func toPredicate(predicateCosts []types.PredicateCost) func(string, int) lo.Tuple2[string, uint64] {
	return func(predicate string, _ int) lo.Tuple2[string, uint64] {
		for _, c := range predicateCosts {
			if matchPredicate(predicate)(c.Predicate) {
				return lo.T2(predicate, c.Cost.Uint64())
			}
		}

		return lo.T2(predicate, defaultCost)
	}
}
