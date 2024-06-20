package keeper

import (
	"context"
	"fmt"
	"io"
	"math"
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
	defaultPredicateCost = uint64(1)
	defaultWeightFactor  = uint64(1)
)

// writerStringer is an interface that combines io.Writer with capabilities of fmt.Stringer.
type writerStringer interface {
	io.Writer
	fmt.Stringer
}

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

func (k Keeper) execute(
	ctx context.Context, program, query string, solutionsLimit sdkmath.Uint,
) (*types.QueryServiceAskResponse, error) {
	ctx = k.enhanceContext(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	limits := k.limits(sdkCtx)

	i, userOutput, err := k.newInterpreter(ctx)
	if err != nil {
		return nil, errorsmod.Wrapf(types.Internal, "error creating interpreter: %v", err.Error())
	}
	if err := i.ExecContext(ctx, program); err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error compiling query: %v", err.Error())
	}

	answer, err := k.queryInterpreter(ctx, i, query, sdkmath.MinUint(solutionsLimit, *limits.MaxResultCount))
	if err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error executing query: %v", err.Error())
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
	return util.QueryInterpreter(ctx, i, query, solutionsLimit)
}

// newInterpreter creates a new interpreter properly configured.
func (k Keeper) newInterpreter(ctx context.Context) (*prolog.Interpreter, fmt.Stringer, error) {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkctx)

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
					if gasError, ok := r.(storetypes.ErrorOutOfGas); ok {
						err = engine.ResourceError(prolog2.ResourceGas(gasError.Descriptor, gasMeter.GasConsumed(), gasMeter.Limit()), env)
						return
					}

					panic(r)
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
	}

	i, err := interpreter.New(options...)

	return i, userOutputBuffer, err
}

func checkLimits(request *types.QueryServiceAskRequest, limits types.Limits) error {
	size := sdkmath.NewUint(uint64(len(request.GetQuery())))
	maxSize := util.DerefOrDefault(limits.MaxSize, sdkmath.NewUint(math.MaxInt64))
	if size.GT(maxSize) {
		return errorsmod.Wrapf(types.LimitExceeded, "query: %d > MaxSize: %d", size, maxSize)
	}

	return nil
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
