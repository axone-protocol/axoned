package bank

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"

	"github.com/axone-protocol/prolog/v3/engine"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/streamingfile"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

const (
	balancesPath = "balances"
	atPath       = "@"
)

type vfs struct {
	ctx context.Context
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates the /v1/bank filesystem.
func NewFS(ctx context.Context) fs.ReadFileFS {
	return &vfs{ctx: ctx}
}

func (f *vfs) Open(name string) (fs.File, error) {
	sdkCtx := sdk.UnwrapSDKContext(f.ctx)

	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	addr, err := f.validatePath(subpath)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	bankKeeper, err := prolog.ContextValue[types.BankKeeper](f.ctx, types.BankKeeperContextKey, nil)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fmt.Errorf("bank keeper not found in context: %w", err)}
	}

	return newStreamingFile(f.ctx, name, sdkCtx.HeaderInfo().Time, bankKeeper, addr), nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	file, err := f.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func (f *vfs) validatePath(subpath string) (sdk.AccAddress, error) {
	segments := strings.Split(subpath, "/")
	if len(segments) != 3 || segments[1] != balancesPath || segments[2] != atPath {
		return nil, fs.ErrNotExist
	}

	address := segments[0]
	if address == "" {
		return nil, fs.ErrNotExist
	}

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	return addr, nil
}

func newStreamingFile(
	ctx context.Context,
	name string,
	modTime time.Time,
	keeper types.BankKeeper,
	addr sdk.AccAddress,
) fs.File {
	open := func() (streamingfile.Next[sdk.Coin], streamingfile.Stop, error) {
		coins := keeper.GetAllBalances(ctx, addr)

		idx := 0
		next := func() (sdk.Coin, bool, error) {
			if idx >= len(coins) {
				return sdk.Coin{}, false, nil
			}
			coin := coins[idx]
			idx++
			return coin, true, nil
		}

		stop := func() error {
			return nil
		}

		return next, stop, nil
	}

	render := func(coin sdk.Coin) ([]byte, error) {
		term := prolog.AtomPair.Apply(
			engine.NewAtom(coin.Denom),
			coinAmountTerm(coin.Amount),
		)
		return prologterm.Render(term, true)
	}

	return streamingfile.New(name, modTime, open, render)
}

func coinAmountTerm(amount math.Int) engine.Term {
	if amount.IsInt64() {
		return engine.Integer(amount.Int64())
	}

	return engine.NewAtom(amount.String())
}
