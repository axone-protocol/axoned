package codec

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"sync"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fsiface "github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
)

func TestCodecDeviceFSOpenFileValidation(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("codec fs should implement OpenFileFS")
	}

	Convey("Given a codec device filesystem", t, func() {
		Convey("when opening with unsupported mode", func() {
			_, err := ofs.OpenFile("bech32", os.O_RDONLY, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})

		Convey("when opening an unknown path", func() {
			_, err := ofs.OpenFile("unknown", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})
}

func TestCodecDeviceFSFunctional(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs := vfs.(fsiface.OpenFileFS)

	readAll := func(codecName string, request []byte) ([]byte, error) {
		file, err := ofs.OpenFile(codecName, os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}

		if len(request) > 0 {
			if _, err := file.(interface{ Write([]byte) (int, error) }).Write(request); err != nil {
				return nil, err
			}
		}

		defer file.Close()
		return io.ReadAll(file)
	}

	testCases := []struct {
		name           string
		codecName      string
		request        []byte
		expectedOutput string
	}{
		// Protocol-level validation (codec-agnostic)
		{
			name:           "empty request",
			codecName:      "bech32",
			request:        nil,
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "unknown command",
			codecName:      "bech32",
			request:        []byte("unknown"),
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "malformed UTF-8",
			codecName:      "bech32",
			request:        []byte{0xff},
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "tab separator not allowed",
			codecName:      "bech32",
			request:        []byte("encode\taxone\t00"),
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "decode with insufficient arguments",
			codecName:      "bech32",
			request:        []byte("decode"),
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "encode with insufficient arguments",
			codecName:      "bech32",
			request:        []byte("encode axone"),
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "encode with too many arguments",
			codecName:      "bech32",
			request:        []byte("encode hrp hex extra"),
			expectedOutput: "error(invalid_request).\n",
		},

		// Bech32 codec - decode tests
		{
			name:           "bech32 decode valid address",
			codecName:      "bech32",
			request:        []byte("decode axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4"),
			expectedOutput: "ok(-(axone,[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30])).\n",
		},
		{
			name:           "bech32 decode with whitespace normalization",
			codecName:      "bech32",
			request:        []byte("  decode   axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4  \r\n"),
			expectedOutput: "ok(-(axone,[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30])).\n",
		},
		{
			name:           "bech32 decode invalid payload",
			codecName:      "bech32",
			request:        []byte("decode bad"),
			expectedOutput: "error(invalid_bech32).\n",
		},

		// Bech32 codec - encode tests
		{
			name:           "bech32 encode valid bytes",
			codecName:      "bech32",
			request:        []byte("encode axone a3a717f4a2af31a2aa0fb58d44868da81238f71e"),
			expectedOutput: "ok(axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4).\n",
		},
		{
			name:           "bech32 encode with uppercase hex",
			codecName:      "bech32",
			request:        []byte("encode axone A3A717F4A2AF31A2AA0FB58D44868DA81238F71E"),
			expectedOutput: "ok(axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4).\n",
		},
		{
			name:           "bech32 encode invalid hex characters",
			codecName:      "bech32",
			request:        []byte("encode axone 0011zz"),
			expectedOutput: "error(invalid_bytes).\n",
		},
		{
			name:           "bech32 encode odd-length hex",
			codecName:      "bech32",
			request:        []byte("encode axone 123"),
			expectedOutput: "error(invalid_bytes).\n",
		},
	}

	Convey("Given a codec device filesystem", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				response, err := readAll(tc.codecName, tc.request)
				So(err, ShouldBeNil)
				So(string(response), ShouldEqual, tc.expectedOutput)
			})
		}
	})
}

func newSDKContext() sdk.Context {
	bech32ConfigOnce.Do(func() {
		sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")
	})
	stateStore := store.NewCommitMultiStore(dbm.NewMemDB(), log.NewNopLogger(), metrics.NewNoOpMetrics())
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return ctx.WithHeaderInfo(coreheader.Info{
		Height: 42,
		Time:   time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC),
	})
}

var bech32ConfigOnce sync.Once
