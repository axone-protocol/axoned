package header

import (
	"context"
	"io/fs"

	"github.com/axone-protocol/prolog/v3/engine"

	coreheader "cosmossdk.io/core/header"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/virtualfile"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
)

const (
	atPath      = "@"
	heightPath  = "height"
	hashPath    = "hash"
	timePath    = "time"
	chainIDPath = "chain_id"
	appHashPath = "app_hash"
)

var (
	atomHeader  = engine.NewAtom("header")
	atomHeight  = engine.NewAtom("height")
	atomHash    = engine.NewAtom("hash")
	atomTime    = engine.NewAtom("time")
	atomChainID = engine.NewAtom("chain_id")
	atomAppHash = engine.NewAtom("app_hash")
)

type vfs struct {
	ctx context.Context
}

type headerTerms struct {
	height  engine.Term
	hash    engine.Term
	time    engine.Term
	chainID engine.Term
	appHash engine.Term
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates the /v1/sys/header snapshot filesystem.
func NewFS(ctx context.Context) fs.ReadFileFS {
	return &vfs{ctx: ctx}
}

func (f *vfs) Open(name string) (fs.File, error) {
	data, err := f.readFile("open", name)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return virtualfile.New(name, data, prolog.ResolveHeaderInfo(sdkCtx).Time), nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	return f.readFile("readfile", name)
}

func (f *vfs) readFile(op, name string) ([]byte, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: op, Path: name, Err: err}
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	headerInfo := prolog.ResolveHeaderInfo(sdkCtx)

	content, err := renderFile(headerInfo, subpath)
	if err != nil {
		return nil, &fs.PathError{Op: op, Path: name, Err: err}
	}

	return content, nil
}

func renderFile(headerInfo coreheader.Info, subpath string) ([]byte, error) {
	terms := newHeaderTerms(headerInfo)

	switch subpath {
	case atPath:
		dictTerm, err := terms.dict()
		if err != nil {
			return nil, err
		}
		return prologterm.Render(dictTerm, true)
	case heightPath:
		return prologterm.Render(terms.height, true)
	case hashPath:
		return prologterm.Render(terms.hash, true)
	case timePath:
		return prologterm.Render(terms.time, true)
	case chainIDPath:
		return prologterm.Render(terms.chainID, true)
	case appHashPath:
		return prologterm.Render(terms.appHash, true)
	default:
		return nil, fs.ErrNotExist
	}
}

func newHeaderTerms(headerInfo coreheader.Info) headerTerms {
	return headerTerms{
		height:  engine.Integer(headerInfo.Height),
		hash:    prolog.BytesToByteListTerm(headerInfo.Hash),
		time:    engine.Integer(headerInfo.Time.Unix()),
		chainID: engine.NewAtom(headerInfo.ChainID),
		appHash: prolog.BytesToByteListTerm(headerInfo.AppHash),
	}
}

func (t headerTerms) dict() (engine.Term, error) {
	return engine.NewDict([]engine.Term{
		atomHeader,
		atomHeight, t.height,
		atomHash, t.hash,
		atomTime, t.time,
		atomChainID, t.chainID,
		atomAppHash, t.appHash,
	})
}
