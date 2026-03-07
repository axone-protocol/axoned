package codec

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"unicode/utf8"

	"github.com/axone-protocol/prolog/v3/engine"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/devfile"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
)

// These are protocol response terms serialized by the device.
var (
	errInvalidRequest = prolog.AtomError.Apply(engine.NewAtom("invalid_request"))
)

const (
	requestCommandDecode = "decode"
	requestCommandEncode = "encode"

	// Generic buffer size for the device input, sized to accommodate
	// all supported codecs (bech32, base64, etc.).
	// Each codec is responsible for its own validation limits.
	maxRequestBytes = 256 * 1024 // 256 KB

	// Maximum size for the serialized Prolog term response.
	// This must account for the worst-case serialization of any codec response.
	maxResponseBytes = 1024 * 1024 // 1 MB
)

type vfs struct {
	ctx context.Context
}

var (
	_ fs.FS            = (*vfs)(nil)
	_ iface.OpenFileFS = (*vfs)(nil)
)

// NewFS creates the codec transactional device filesystem.
//
// The device follows a half-duplex request/response protocol:
//  1. writes build a request buffer;
//  2. the first read commits the request;
//  3. subsequent reads stream a serialized Prolog term response until EOF.
//
// Request validation failures are reported in-band as serialized Prolog terms.
// Stream/runtime failures remain regular VFS errors surfaced by devfile.
//
// The VFS supports multiple codecs via relative paths: {codec_name}
// The codec name is resolved from the filesystem path relative to the mount point.
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

	codec := Get(subpath)
	if codec == nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return devfile.New(
		devfile.WithPath(name),
		devfile.WithModTime(prolog.ResolveHeaderInfo(sdkCtx).Time),
		devfile.WithMaxRequestBytes(maxRequestBytes),
		devfile.WithMaxResponseBytes(maxResponseBytes),
		devfile.WithCommit(makeCommitFunc(codec)),
	)
}

// makeCommitFunc creates a closure that captures the codec for request processing.
func makeCommitFunc(codec Codec) func(io.Reader, io.Writer) error {
	return func(reader io.Reader, writer io.Writer) error {
		request, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		responseTerm := handleRequest(codec, request)
		responseBytes, err := prologterm.Render(responseTerm, true)
		if err != nil {
			return err
		}

		_, err = writer.Write(responseBytes)
		return err
	}
}

// commitRequest is the protocol boundary for the codec transactional device.
//
// Text request framing:
//   - decode <input>
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
//   - decode success: codec-specific (e.g., ok(HRP-Bytes) for bech32)
//   - encode success: codec-specific (e.g., ok(Bech32) for bech32)
//   - protocol error: error(Code)
//
// Protocol errors are always encoded in-band as response terms, for example
// error(invalid_request). Regular Go errors are reserved for VFS/runtime
// failures such as an internal rendering problem.
//
// handleRequest normalizes the text line and dispatches on the parsed command.
//
// Any malformed request line is mapped to the explicit protocol response term
// error(invalid_request), rather than escaping as a Go error.
func handleRequest(codec Codec, request []byte) engine.Term {
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

		return codec.Decode(tokens[1:])
	case requestCommandEncode:
		if len(tokens) != 3 {
			return errInvalidRequest
		}

		return codec.Encode(tokens[1:])
	default:
		return errInvalidRequest
	}
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
