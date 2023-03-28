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
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"
	u "github.com/rjNemo/underscore"
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

	whitelist := util.NonZeroOrDefault(interpreterParams.PredicatesWhitelist, interpreter.RegistryNames)
	blacklist := interpreterParams.GetPredicatesBlacklist()
	interpreted, err := interpreter.New(
		ctx,
		filterPredicates(interpreter.RegistryNames, whitelist, blacklist),
		util.NonZeroOrDefault(interpreterParams.GetBootstrap(), bootstrap.Bootstrap()),
		sdkctx.GasMeter(),
		k.fsProvider(ctx),
	)

	return interpreted, err
}

// filterPredicates constructs the predicate list from the given registry.
// The given whitelist and blacklist are applied to the registry to determine
// the final predicate list.
func filterPredicates(registry []string, whitelist []string, blacklist []string) []string {
	match := func(a string) func(b string) bool {
		return func(b string) bool {
			if strings.Contains(b, "/") {
				return a == b
			}
			return strings.Split(a, "/")[0] == b
		}
	}

	return util.NonZeroOrDefault(
		u.NewPipe(registry).
			Filter(func(predicate string) bool {
				return u.Any(whitelist, match(predicate))
			}).
			Filter(func(predicate string) bool {
				return !u.Any(blacklist, match(predicate))
			}).
			Value,
		[]string{})
}
