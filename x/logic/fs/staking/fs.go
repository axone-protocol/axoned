package staking

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"strings"
	"time"

	"github.com/axone-protocol/prolog/v3/engine"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/streamingfile"
	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
	logictypes "github.com/axone-protocol/axoned/v15/x/logic/types"
)

const (
	opOpen = "open"

	delegationsPath                 = "delegations"
	unbondingDelegationsPath        = "unbonding_delegations"
	redelegationsPath               = "redelegations"
	atPath                          = "@"
	pageSize                 uint64 = 100
)

var (
	atomDelegation           = engine.NewAtom("delegation")
	atomUnbondingDelegation  = engine.NewAtom("unbonding_delegation")
	atomRedelegation         = engine.NewAtom("redelegation")
	atomUnbondingEntry       = engine.NewAtom("unbonding_entry")
	atomRedelegationEntry    = engine.NewAtom("redelegation_entry")
	atomCoin                 = engine.NewAtom("coin")
	atomValidator            = engine.NewAtom("validator")
	atomSourceValidator      = engine.NewAtom("source_validator")
	atomDestinationValidator = engine.NewAtom("destination_validator")
	atomBalance              = engine.NewAtom("balance")
	atomEntries              = engine.NewAtom("entries")
	atomCompletionTime       = engine.NewAtom("completion_time")

	errVFSUnavailable = errors.New("vfs_unavailable")
)

type vfs struct {
	ctx context.Context
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates a read-only filesystem exposing delegator staking positions.
func NewFS(ctx context.Context) fs.ReadFileFS {
	return &vfs{ctx: ctx}
}

func (f *vfs) Open(name string) (fs.File, error) {
	sdkCtx := sdk.UnwrapSDKContext(f.ctx)

	address, collection, err := validatePath(name)
	if err != nil {
		return nil, &fs.PathError{Op: opOpen, Path: name, Err: err}
	}

	service, err := prolog.ContextValue[logictypes.StakingQueryService](f.ctx, logictypes.StakingQueryServiceContextKey, nil)
	if err != nil {
		return nil, &fs.PathError{Op: opOpen, Path: name, Err: errVFSUnavailable}
	}

	modTime := prolog.ResolveHeaderInfo(sdkCtx).Time
	switch collection {
	case delegationsPath:
		return newDelegationsFile(f.ctx, name, modTime, service, address.String()), nil
	case unbondingDelegationsPath:
		return newUnbondingDelegationsFile(f.ctx, name, modTime, service, address.String())
	case redelegationsPath:
		return newRedelegationsFile(f.ctx, name, modTime, service, address.String())
	default:
		return nil, &fs.PathError{Op: opOpen, Path: name, Err: fs.ErrNotExist}
	}
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	file, err := f.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func validatePath(name string) (sdk.AccAddress, string, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, "", err
	}

	segments := strings.Split(subpath, "/")
	if len(segments) != 3 || segments[2] != atPath {
		return nil, "", fs.ErrNotExist
	}

	address, err := sdk.AccAddressFromBech32(segments[0])
	if err != nil {
		return nil, "", fs.ErrNotExist
	}

	switch segments[1] {
	case delegationsPath, unbondingDelegationsPath, redelegationsPath:
		return address, segments[1], nil
	default:
		return nil, "", fs.ErrNotExist
	}
}

func newDelegationsFile(
	ctx context.Context, name string, modTime time.Time, service logictypes.StakingQueryService, address string,
) fs.File {
	return streamingfile.New(name, modTime, newPagedCursor(func(key []byte) ([]stakingtypes.DelegationResponse, []byte, error) {
		response, err := service.DelegatorDelegations(ctx, &stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: address,
			Pagination:    &query.PageRequest{Key: key, Limit: pageSize},
		})
		if err != nil {
			return nil, nil, err
		}
		return response.DelegationResponses, response.Pagination.GetNextKey(), nil
	}), renderDelegation)
}

func newUnbondingDelegationsFile(
	ctx context.Context, name string, modTime time.Time, service logictypes.StakingQueryService, address string,
) (fs.File, error) {
	bondDenom, err := queryBondDenom(ctx, service)
	if err != nil {
		return nil, err
	}

	return streamingfile.New(name, modTime, newPagedCursor(func(key []byte) ([]stakingtypes.UnbondingDelegation, []byte, error) {
		response, err := service.DelegatorUnbondingDelegations(ctx, &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
			DelegatorAddr: address,
			Pagination:    &query.PageRequest{Key: key, Limit: pageSize},
		})
		if err != nil {
			return nil, nil, err
		}
		return response.UnbondingResponses, response.Pagination.GetNextKey(), nil
	}), func(delegation stakingtypes.UnbondingDelegation) ([]byte, error) {
		return renderUnbondingDelegation(delegation, bondDenom)
	}), nil
}

func newRedelegationsFile(
	ctx context.Context, name string, modTime time.Time, service logictypes.StakingQueryService, address string,
) (fs.File, error) {
	bondDenom, err := queryBondDenom(ctx, service)
	if err != nil {
		return nil, err
	}

	return streamingfile.New(name, modTime, newPagedCursor(func(key []byte) ([]stakingtypes.RedelegationResponse, []byte, error) {
		response, err := service.Redelegations(ctx, &stakingtypes.QueryRedelegationsRequest{
			DelegatorAddr: address,
			Pagination:    &query.PageRequest{Key: key, Limit: pageSize},
		})
		if err != nil {
			return nil, nil, err
		}
		return response.RedelegationResponses, response.Pagination.GetNextKey(), nil
	}), func(redelegation stakingtypes.RedelegationResponse) ([]byte, error) {
		return renderRedelegation(redelegation, bondDenom)
	}), nil
}

func queryBondDenom(ctx context.Context, service logictypes.StakingQueryService) (string, error) {
	response, err := service.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return "", err
	}
	return response.Params.BondDenom, nil
}

func newPagedCursor[T any](fetch func([]byte) ([]T, []byte, error)) streamingfile.OpenCursor[T] {
	return func() (streamingfile.Next[T], streamingfile.Stop, error) {
		var (
			items   []T
			index   int
			nextKey []byte
			done    bool
		)

		next := func() (T, bool, error) {
			var zero T
			for index == len(items) && !done {
				page, key, err := fetch(nextKey)
				if err != nil {
					return zero, false, err
				}
				items = page
				index = 0
				nextKey = bytes.Clone(key)
				done = len(nextKey) == 0
				if len(items) == 0 && done {
					return zero, false, nil
				}
			}
			if index == len(items) {
				return zero, false, nil
			}

			item := items[index]
			index++
			return item, true, nil
		}

		return next, func() error { return nil }, nil
	}
}

func renderDelegation(response stakingtypes.DelegationResponse) ([]byte, error) {
	term, err := engine.NewDict([]engine.Term{
		atomDelegation,
		atomValidator, engine.NewAtom(response.Delegation.ValidatorAddress),
		atomBalance, coinTerm(response.Balance.Denom, response.Balance.Amount),
	})
	if err != nil {
		return nil, err
	}
	return prologterm.Render(term, true)
}

func renderUnbondingDelegation(delegation stakingtypes.UnbondingDelegation, bondDenom string) ([]byte, error) {
	entries := make([]engine.Term, 0, len(delegation.Entries))
	for _, entry := range delegation.Entries {
		term, err := engine.NewDict([]engine.Term{
			atomUnbondingEntry,
			atomBalance, coinTerm(bondDenom, entry.Balance),
			atomCompletionTime, engine.Integer(entry.CompletionTime.Unix()),
		})
		if err != nil {
			return nil, err
		}
		entries = append(entries, term)
	}

	term, err := engine.NewDict([]engine.Term{
		atomUnbondingDelegation,
		atomValidator, engine.NewAtom(delegation.ValidatorAddress),
		atomEntries, engine.List(entries...),
	})
	if err != nil {
		return nil, err
	}
	return prologterm.Render(term, true)
}

func renderRedelegation(response stakingtypes.RedelegationResponse, bondDenom string) ([]byte, error) {
	entries := make([]engine.Term, 0, len(response.Entries))
	for _, entry := range response.Entries {
		term, err := engine.NewDict([]engine.Term{
			atomRedelegationEntry,
			atomBalance, coinTerm(bondDenom, entry.Balance),
			atomCompletionTime, engine.Integer(entry.RedelegationEntry.CompletionTime.Unix()),
		})
		if err != nil {
			return nil, err
		}
		entries = append(entries, term)
	}

	term, err := engine.NewDict([]engine.Term{
		atomRedelegation,
		atomSourceValidator, engine.NewAtom(response.Redelegation.ValidatorSrcAddress),
		atomDestinationValidator, engine.NewAtom(response.Redelegation.ValidatorDstAddress),
		atomEntries, engine.List(entries...),
	})
	if err != nil {
		return nil, err
	}
	return prologterm.Render(term, true)
}

func coinTerm(denom string, amount math.Int) engine.Term {
	return atomCoin.Apply(engine.NewAtom(denom), amountTerm(amount))
}

func amountTerm(amount math.Int) engine.Term {
	if amount.IsInt64() {
		return engine.Integer(amount.Int64())
	}
	return engine.NewAtom(amount.String())
}
