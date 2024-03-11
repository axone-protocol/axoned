package keeper

import (
	"context"
	"math"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/v7/x/logic/fs"
	"github.com/okp4/okp4d/v7/x/logic/interpreter"
	"github.com/okp4/okp4d/v7/x/logic/interpreter/bootstrap"
	"github.com/okp4/okp4d/v7/x/logic/meter"
	prolog2 "github.com/okp4/okp4d/v7/x/logic/prolog"
	"github.com/okp4/okp4d/v7/x/logic/types"
	"github.com/okp4/okp4d/v7/x/logic/util"
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

func (k Keeper) execute(ctx context.Context, program, query string, limit sdkmath.Uint) (*types.QueryServiceAskResponse, error) {
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

	answer, err := k.queryInterpreter(ctx, i, query, sdkmath.MinUint(limit, *limits.MaxResultCount))
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
	return util.QueryInterpreter(ctx, i, query, limit)
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

	options := []interpreter.Option{
		interpreter.WithPredicates(ctx, interpreter.RegistryNames, hook),
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
