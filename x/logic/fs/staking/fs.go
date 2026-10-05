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
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/virtualfile"
	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
	logictypes "github.com/axone-protocol/axoned/v15/x/logic/types"
)

const (
	opOpen = "open"

	delegationsPath                 = "delegations"
	unbondingDelegationsPath        = "unbonding_delegations"
	redelegationsPath               = "redelegations"
	validatorsPath                  = "validators"
	paramsPath                      = "params"
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
	atomOperator             = engine.NewAtom("operator")
	atomStatus               = engine.NewAtom("status")
	atomStakingParams        = engine.NewAtom("staking_params")
	atomBondDenom            = engine.NewAtom("bond_denom")

	errVFSUnavailable = errors.New("vfs_unavailable")
)

type vfs struct {
	ctx context.Context
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates a read-only filesystem exposing staking positions, validators, and bond denomination.
func NewFS(ctx context.Context) fs.ReadFileFS {
	return &vfs{ctx: ctx}
}

func (f *vfs) Open(name string) (fs.File, error) {
	sdkCtx := sdk.UnwrapSDKContext(f.ctx)

	subject, collection, err := validatePath(name)
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
		return newDelegationsFile(f.ctx, name, modTime, service, subject), nil
	case unbondingDelegationsPath:
		return newUnbondingDelegationsFile(f.ctx, name, modTime, service, subject)
	case redelegationsPath:
		return newRedelegationsFile(f.ctx, name, modTime, service, subject)
	case validatorsPath:
		return newValidatorsFile(f.ctx, name, modTime, service, subject), nil
	case paramsPath:
		denom, err := queryBondDenom(f.ctx, service)
		if err != nil {
			return nil, &fs.PathError{Op: opOpen, Path: name, Err: err}
		}
		// The tag and key are fixed atoms, and the value is always an atom. Rendering this
		// finite term to an in-memory stream cannot fail.
		term, _ := engine.NewDict([]engine.Term{atomStakingParams, atomBondDenom, engine.NewAtom(denom)})
		content, _ := prologterm.Render(term, true)
		return virtualfile.New(name, content, modTime), nil
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

func validatePath(name string) (string, string, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return "", "", err
	}

	segments := strings.Split(subpath, "/")
	if len(segments) == 2 && segments[0] == paramsPath && segments[1] == atPath {
		return "", paramsPath, nil
	}
	if len(segments) != 3 || segments[2] != atPath {
		return "", "", fs.ErrNotExist
	}
	if segments[0] == validatorsPath {
		switch segments[1] {
		case "bonded", "unbonding", "unbonded":
			return segments[1], validatorsPath, nil
		default:
			return "", "", fs.ErrNotExist
		}
	}

	address, err := sdk.AccAddressFromBech32(segments[0])
	if err != nil {
		return "", "", fs.ErrNotExist
	}
	switch segments[1] {
	case delegationsPath, unbondingDelegationsPath, redelegationsPath:
		return address.String(), segments[1], nil
	default:
		return "", "", fs.ErrNotExist
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

func newValidatorsFile(
	ctx context.Context, name string, modTime time.Time, service logictypes.StakingQueryService, status string,
) fs.File {
	status = "BOND_STATUS_" + strings.ToUpper(status)
	return streamingfile.New(name, modTime, newPagedCursor(func(key []byte) ([]stakingtypes.Validator, []byte, error) {
		response, err := service.Validators(ctx, &stakingtypes.QueryValidatorsRequest{
			Status:     status,
			Pagination: &query.PageRequest{Key: key, Limit: pageSize},
		})
		if err != nil {
			return nil, nil, err
		}
		return response.Validators, response.Pagination.GetNextKey(), nil
	}), renderValidator)
}

func renderValidator(validator stakingtypes.Validator) ([]byte, error) {
	status := strings.ToLower(strings.TrimPrefix(validator.Status.String(), "BOND_STATUS_"))
	term, err := engine.NewDict([]engine.Term{
		atomValidator,
		atomOperator, engine.NewAtom(validator.OperatorAddress),
		atomStatus, engine.NewAtom(status),
	})
	if err != nil {
		return nil, err
	}
	return prologterm.Render(term, true)
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
