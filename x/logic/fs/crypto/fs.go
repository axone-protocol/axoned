package cryptofs

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/axone-protocol/prolog/v3/engine"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/devfile"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
	"github.com/axone-protocol/axoned/v15/x/logic/util"
)

const (
	maxRequestBytes  = 256 * 1024
	maxResponseBytes = 1024

	requestCommandVerify = "verify"
)

type vfs struct {
	ctx context.Context
}

type device struct {
	hashAlg util.HashAlg
	keyAlg  util.KeyAlg
	kind    deviceKind
}

type deviceKind int

const (
	deviceKindHash deviceKind = iota + 1
	deviceKindSignature
)

var (
	_ fs.FS            = (*vfs)(nil)
	_ iface.OpenFileFS = (*vfs)(nil)

	atomOK    = engine.NewAtom("ok")
	atomTrue  = engine.NewAtom("true")
	atomFalse = engine.NewAtom("false")

	errInvalidRequest       = prolog.AtomError.Apply(engine.NewAtom("invalid_request"))
	errUnsupportedOperation = prolog.AtomError.Apply(engine.NewAtom("unsupported_operation"))
	errInvalidKey           = prolog.AtomError.Apply(engine.NewAtom("invalid_key"))
)

// NewFS creates a transactional device filesystem for crypto utilities.
func NewFS(ctx context.Context) fs.FS {
	return &vfs{ctx: ctx}
}

func (f *vfs) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
}

func (f *vfs) OpenFile(name string, flag int, _ fs.FileMode) (fs.File, error) {
	if flag != os.O_RDWR {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
	}

	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	dev, err := resolveDevice(subpath)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return devfile.New(
		devfile.WithPath(name),
		devfile.WithModTime(prolog.ResolveHeaderInfo(sdkCtx).Time),
		devfile.WithMaxRequestBytes(maxRequestBytes),
		devfile.WithMaxResponseBytes(maxResponseBytes),
		devfile.WithCommit(dev.makeCommitFunc()),
	)
}

func resolveDevice(subpath string) (device, error) {
	if strings.Contains(subpath, "/") {
		return device{}, fs.ErrNotExist
	}

	if algorithm, err := util.ParseHashAlg(subpath); err == nil {
		return device{hashAlg: algorithm, kind: deviceKindHash}, nil
	}

	if algorithm, err := util.ParseKeyAlg(subpath); err == nil {
		return device{keyAlg: algorithm, kind: deviceKindSignature}, nil
	}

	return device{}, fs.ErrNotExist
}

func (d device) makeCommitFunc() func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		request, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		var response []byte
		switch d.kind {
		case deviceKindHash:
			response, err = util.Hash(d.hashAlg, request)
		case deviceKindSignature:
			response, err = handleSignatureRequest(d.keyAlg, request)
		default:
			err = fs.ErrInvalid
		}
		if err != nil {
			return err
		}

		_, err = w.Write(response)
		return err
	}
}

func handleSignatureRequest(algorithm util.KeyAlg, request []byte) ([]byte, error) {
	line, ok := normalizeRequestLine(request)
	if !ok {
		return renderResponse(errInvalidRequest)
	}

	tokens := splitRequestLine(line)
	if len(tokens) == 0 {
		return renderResponse(errInvalidRequest)
	}
	if string(tokens[0]) != requestCommandVerify {
		return renderResponse(errUnsupportedOperation)
	}
	if len(tokens) != 4 {
		return renderResponse(errInvalidRequest)
	}

	pubKey, err := decodeHexToken(tokens[1])
	if err != nil {
		return renderResponse(errInvalidRequest)
	}
	msg, err := decodeHexToken(tokens[2])
	if err != nil {
		return renderResponse(errInvalidRequest)
	}
	sig, err := decodeHexToken(tokens[3])
	if err != nil {
		return renderResponse(errInvalidRequest)
	}

	verified, err := util.VerifySignature(algorithm, pubKey, msg, sig)
	if err != nil {
		return renderResponse(errInvalidKey)
	}
	if verified {
		return renderResponse(atomOK.Apply(atomTrue))
	}

	return renderResponse(atomOK.Apply(atomFalse))
}

func normalizeRequestLine(request []byte) ([]byte, bool) {
	if len(request) == 0 {
		return nil, false
	}

	request = bytes.TrimRight(request, "\r\n")
	request = bytes.Trim(request, " ")
	if len(request) == 0 || !utf8.Valid(request) {
		return nil, false
	}

	for _, b := range request {
		switch {
		case b == ' ':
			continue
		case b == '\t', b == '\n', b == '\r':
			return nil, false
		case b < 0x20 || b == 0x7f:
			return nil, false
		}
	}

	return request, true
}

func splitRequestLine(line []byte) [][]byte {
	tokens := make([][]byte, 0, 4)

	for start := 0; start < len(line); {
		for start < len(line) && line[start] == ' ' {
			start++
		}
		if start >= len(line) {
			break
		}

		end := start
		for end < len(line) && line[end] != ' ' {
			end++
		}
		tokens = append(tokens, line[start:end])
		start = end
	}

	return tokens
}

func decodeHexToken(token []byte) ([]byte, error) {
	if len(token)%2 != 0 {
		return nil, hex.ErrLength
	}

	dst := make([]byte, hex.DecodedLen(len(token)))
	_, err := hex.Decode(dst, token)
	return dst, err
}

func renderResponse(term engine.Term) ([]byte, error) {
	return prologterm.Render(term, true)
}
