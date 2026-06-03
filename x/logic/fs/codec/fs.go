package codec

import (
	"context"
	"io"
	"io/fs"
	"math"
	"os"
	"unicode/utf8"

	"github.com/axone-protocol/prolog/v3/engine"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/devfile"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
	"github.com/axone-protocol/axoned/v15/x/logic/types"
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

// NewFS creates a transactional device filesystem for codec utilities.
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
	ioCoeff := ioCoeffFromContext(f.ctx)
	return devfile.New(
		devfile.WithPath(name),
		devfile.WithModTime(prolog.ResolveHeaderInfo(sdkCtx).Time),
		devfile.WithMaxRequestBytes(maxRequestBytes),
		devfile.WithMaxResponseBytes(maxResponseBytes),
		devfile.WithCommit(makeCommitFunc(codec)),
		devfile.WithTransferHook(func(_ devfile.TransferDirection, n int) {
			consumeTransferredIOGas(sdkCtx.GasMeter(), n, ioCoeff)
		}),
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

// handleRequest is the protocol boundary for the codec transactional device.
//
// Text request framing:
//   - decode <payload>
//   - encode <payload>
//
// Whitespace rules:
//   - leading ASCII spaces before the command are ignored;
//   - the command must be followed by a separator: SPACE, LF, or CRLF;
//   - the payload is forwarded unchanged after the command separator;
//   - TAB and any other control whitespace in the command are invalid_request.
//
// Commands and separators are ASCII. Payload syntax is codec-specific.
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
// handleRequest dispatches on the parsed command and leaves payload parsing to
// the selected codec.
//
// Any malformed request is mapped to the explicit protocol response term
// error(invalid_request), rather than escaping as a Go error.
func handleRequest(codec Codec, request []byte) engine.Term {
	command, payload, ok := splitRequestCommand(request)
	if !ok {
		return errInvalidRequest
	}

	switch string(command) {
	case requestCommandDecode:
		return codec.Decode(payload)
	case requestCommandEncode:
		return codec.Encode(payload)
	default:
		return errInvalidRequest
	}
}

func splitRequestCommand(request []byte) ([]byte, []byte, bool) {
	if len(request) == 0 || !utf8.Valid(request) {
		return nil, nil, false
	}

	for len(request) > 0 && request[0] == ' ' {
		request = request[1:]
	}
	if len(request) == 0 {
		return nil, nil, false
	}

	commandEnd := 0
	for commandEnd < len(request) {
		switch request[commandEnd] {
		case ' ', '\n', '\r':
			break
		case '\t':
			return nil, nil, false
		default:
			if request[commandEnd] < 0x20 || request[commandEnd] == 0x7f {
				return nil, nil, false
			}
			commandEnd++
			continue
		}
		break
	}
	if commandEnd == 0 || commandEnd == len(request) {
		return nil, nil, false
	}

	switch request[commandEnd] {
	case ' ':
		return request[:commandEnd], request[commandEnd+1:], true
	case '\n':
		return request[:commandEnd], request[commandEnd+1:], true
	case '\r':
		if commandEnd+1 < len(request) && request[commandEnd+1] == '\n' {
			return request[:commandEnd], request[commandEnd+2:], true
		}
	}

	return nil, nil, false
}

func ioCoeffFromContext(ctx context.Context) uint64 {
	coeff, _ := ctx.Value(types.IOCoeffContextKey).(uint64)
	return coeff
}

func consumeTransferredIOGas(gasMeter storetypes.GasMeter, transferred int, coeff uint64) {
	if transferred <= 0 {
		return
	}

	consumeIOGas(gasMeter, uint64(transferred), coeff)
}

func consumeIOGas(gasMeter storetypes.GasMeter, units, coeff uint64) {
	if units == 0 {
		return
	}
	if coeff == 0 {
		coeff = 1
	}

	consumed, overflow := multiplyUint64Overflow(units, coeff)
	if overflow {
		gasMeter.ConsumeGas(math.MaxUint64, "IO")
		return
	}

	gasMeter.ConsumeGas(consumed, "IO")
}

func multiplyUint64Overflow(a, b uint64) (uint64, bool) {
	if a == 0 || b == 0 {
		return 0, false
	}

	c := a * b
	if c/a != b || c/b != a {
		return 0, true
	}

	return c, false
}
