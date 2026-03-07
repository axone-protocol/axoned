package codec

import (
	"bytes"
	"context"
	"encoding/hex"
	"io/fs"
	"os"
	"unicode/utf8"

	"github.com/axone-protocol/prolog/v3/engine"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkbech32 "github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/devfile"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
)

// These are protocol response terms serialized by the device.
var (
	atomOK = engine.NewAtom("ok")

	errInvalidRequest = prolog.AtomError.Apply(engine.NewAtom("invalid_request"))
	errInvalidBech32  = prolog.AtomError.Apply(engine.NewAtom("invalid_bech32"))
	errInvalidBytes   = prolog.AtomError.Apply(engine.NewAtom("invalid_bytes"))
	errInvalidHrp     = prolog.AtomError.Apply(engine.NewAtom("invalid_hrp"))
)

const (
	devicePath = "bech32"

	maxHRPBytes      = 255
	maxDataBytes     = 65535
	maxChecksumChars = 6

	requestCommandDecode = "decode"
	requestCommandEncode = "encode"

	// The decode request is delimited by EOF/LF/CRLF, so the request buffer must
	// account for the largest Bech32 string derivable from the largest supported
	// raw byte payload.
	maxDecodeDataChars = (maxDataBytes*8 + 4) / 5

	// Request size limits are for the device input buffer only. They account for
	// the canonical command form with single ASCII spaces and an optional CRLF.
	maxEncodeRequestBytes = len(requestCommandEncode) + 1 + maxHRPBytes + 1 + maxDataBytes*2 + 2
	maxDecodeRequestBytes = len(requestCommandDecode) + 1 + maxHRPBytes + 1 + maxDecodeDataChars + maxChecksumChars + 2
	maxRequestBytes       = maxEncodeRequestBytes

	// The response is a serialized Prolog term streamed by devfile.
	// It is intentionally left uncapped here: this is not a protocol-level
	// error condition for the codec itself.
	maxResponseBytes = 0
)

type vfs struct {
	ctx context.Context
}

var (
	_ fs.FS            = (*vfs)(nil)
	_ iface.OpenFileFS = (*vfs)(nil)
)

// NewFS creates the bech32 codec transactional device filesystem.
//
// The device follows a half-duplex request/response protocol:
//  1. writes build a request buffer;
//  2. the first read commits the request;
//  3. subsequent reads stream a serialized Prolog term response until EOF.
//
// Request validation failures are reported in-band as serialized Prolog terms.
// Stream/runtime failures remain regular VFS errors surfaced by devfile.
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
	if subpath != devicePath {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return devfile.New(
		devfile.WithPath(name),
		devfile.WithModTime(prolog.ResolveHeaderInfo(sdkCtx).Time),
		devfile.WithMaxRequestBytes(maxRequestBytes),
		devfile.WithMaxResponseBytes(maxResponseBytes),
		devfile.WithAllowEmptyRequest(true),
		devfile.WithCommit(commitRequest),
	)
}

// commitRequest is the protocol boundary for the bech32 codec transactional device.
//
// Text request framing:
//   - decode <bech32>
//   - encode <hrp> <hex>
//
// Whitespace rules:
//   - leading and trailing ASCII spaces are ignored;
//   - token separators are one or more ASCII spaces;
//   - requests may terminate with EOF, LF, or CRLF;
//   - TAB and any other control whitespace are invalid_request.
//
// Commands and separators are ASCII. Tokens are UTF-8 text without embedded
// ASCII spaces. The encode payload hex is lowercase/uppercase agnostic.
//
// Serialized Prolog response terms:
//   - decode success: ok(HRP-Bytes).
//   - encode success: ok(Bech32).
//   - protocol error: error(Code).
//
// Protocol errors are always encoded in-band as response terms, for example
// error(invalid_request). Regular Go errors are reserved for VFS/runtime
// failures such as an internal rendering problem.
func commitRequest(request []byte) ([]byte, error) {
	return prologterm.Render(handleRequest(request), true)
}

// handleRequest normalizes the text line and dispatches on the parsed command.
//
// Any malformed request line is mapped to the explicit protocol response term
// error(invalid_request), rather than escaping as a Go error.
func handleRequest(request []byte) engine.Term {
	line, ok := normalizeRequestLine(request)
	if !ok {
		return errInvalidRequest
	}

	tokens := splitRequestLine(line)
	if len(tokens) == 0 {
		return errInvalidRequest
	}

	switch string(tokens[0]) {
	case requestCommandDecode:
		if len(tokens) != 2 {
			return errInvalidRequest
		}

		return handleDecode(tokens[1])
	case requestCommandEncode:
		if len(tokens) != 3 {
			return errInvalidRequest
		}

		return handleEncode(tokens[1], tokens[2])
	default:
		return errInvalidRequest
	}
}

func normalizeRequestLine(request []byte) ([]byte, bool) {
	if len(request) == 0 {
		return nil, false
	}

	switch {
	case bytes.HasSuffix(request, []byte("\r\n")):
		request = request[:len(request)-2]
	case bytes.HasSuffix(request, []byte("\n")):
		request = request[:len(request)-1]
	}

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
	tokens := make([][]byte, 0, 3)

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

func handleDecode(bech32Text []byte) engine.Term {
	hrp, data, err := sdkbech32.DecodeAndConvert(string(bech32Text))
	if err != nil {
		// The request line is well-formed, but the Bech32 payload is invalid.
		return errInvalidBech32
	}

	return atomOK.Apply(prolog.AtomPair.Apply(
		engine.NewAtom(hrp),
		prolog.BytesToByteListTerm(data),
	))
}

func handleEncode(hrpText, hexText []byte) engine.Term {
	if len(hrpText) > maxHRPBytes {
		return errInvalidHrp
	}

	data, err := decodeHex(hexText)
	if err != nil {
		return errInvalidBytes
	}
	if len(data) > maxDataBytes {
		return errInvalidBytes
	}

	bech32Address, err := sdkbech32.ConvertAndEncode(string(hrpText), data)
	if err != nil {
		// The request line is well-formed, but the HRP is invalid for Bech32.
		return errInvalidHrp
	}

	return atomOK.Apply(engine.NewAtom(bech32Address))
}

func decodeHex(hexText []byte) ([]byte, error) {
	if len(hexText)%2 != 0 {
		return nil, hex.ErrLength
	}

	data := make([]byte, hex.DecodedLen(len(hexText)))
	if _, err := hex.Decode(data, hexText); err != nil {
		return nil, err
	}

	return data, nil
}
