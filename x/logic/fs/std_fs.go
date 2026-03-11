//nolint:revive
package fs

import (
	goctx "context"
	"fmt"
	stdfs "io/fs"

	sdk "github.com/cosmos/cosmos-sdk/types"

	logicbank "github.com/axone-protocol/axoned/v14/x/logic/fs/bank"
	logiccodec "github.com/axone-protocol/axoned/v14/x/logic/fs/codec"
	logicembeddedfs "github.com/axone-protocol/axoned/v14/x/logic/fs/embedded"
	logicshare "github.com/axone-protocol/axoned/v14/x/logic/fs/share"
	logicsyscomet "github.com/axone-protocol/axoned/v14/x/logic/fs/sys/comet"
	logicsysheader "github.com/axone-protocol/axoned/v14/x/logic/fs/sys/header"
	logicvfs "github.com/axone-protocol/axoned/v14/x/logic/fs/vfs"
	logicwasm "github.com/axone-protocol/axoned/v14/x/logic/fs/wasm"
	logiclib "github.com/axone-protocol/axoned/v14/x/logic/lib"
	logictypes "github.com/axone-protocol/axoned/v14/x/logic/types"
)

const (
	v1Root = "/v1"

	// Canonical host namespace paths.
	libPath         = v1Root + "/lib"
	runHeaderPath   = v1Root + "/run/header"
	runCometPath    = v1Root + "/run/comet"
	varLibBankPath  = v1Root + "/var/lib/bank"
	varLibLogicPath = v1Root + "/var/lib/logic/users"
	devCodecPath    = v1Root + "/dev/codec"
	devWasmPath     = v1Root + "/dev/wasm"
)

// Mount defines a filesystem mounted at an absolute path in the logic VFS.
type Mount struct {
	Path string
	FS   stdfs.FS
}

// ProgramKeeper defines the subset of keeper capabilities required by the user library FS.
type ProgramKeeper interface {
	GetStoredProgram(ctx sdk.Context, programID []byte) (logictypes.StoredProgram, bool, error)
	GetProgramPublication(ctx sdk.Context, publisher, programID []byte) (logictypes.ProgramPublication, bool, error)
}

// StandardMounts returns the canonical capability mounts shared across app and tests.
//
// The virtual filesystem is the boundary between logical evaluation and the host
// environment.
//
// The host exposes capabilities as addressable resources.
// Prolog libraries expose the logical relations that operate over them.
//
// The capability hierarchy is exposed through a versioned canonical namespace
// rooted at /v1.
//
// Within this namespace, the path segment "@" denotes the canonical,
// host-defined synthetic view of the addressed capability, intended as its
// primary complete representation.
//
// The hierarchy follows a stable Unix-inspired discipline:
//   - /v1/lib contains immutable host-provided libraries and reference programs.
//   - /v1/run contains invocation-scoped runtime resources supplied by the host.
//   - /v1/var/lib contains persistent host-managed resources exposed as files.
//   - /v1/dev contains interactive device-like capabilities.
//
// Concrete mounts depend on the host environment. Different hosts may expose
// different resources while preserving the same organizing principles.
func StandardMounts(ctx goctx.Context, wasmKeeper logictypes.WasmKeeper, programKeeper ProgramKeeper) []Mount {
	libFS := logicembeddedfs.NewFS(logiclib.Files)
	headerFS := logicsysheader.NewFS(ctx)
	cometFS := logicsyscomet.NewFS(ctx)
	bankFS := logicbank.NewFS(ctx)
	codecFS := logiccodec.NewFS(ctx)
	wasmFS := logicwasm.NewFS(ctx, wasmKeeper)
	shareFS := logicshare.NewFS(ctx, programKeeper)

	mounts := make([]Mount, 0, 7)
	mounts = append(mounts,
		Mount{Path: libPath, FS: libFS},
		Mount{Path: runHeaderPath, FS: headerFS},
		Mount{Path: runCometPath, FS: cometFS},
		Mount{Path: varLibBankPath, FS: bankFS},
		Mount{Path: devCodecPath, FS: codecFS},
		Mount{Path: devWasmPath, FS: wasmFS},
		Mount{Path: varLibLogicPath, FS: shareFS},
	)

	return mounts
}

// NewVFS builds a logic virtual filesystem with standard mounts and optional extras.
func NewVFS(
	ctx goctx.Context, wasmKeeper logictypes.WasmKeeper, programKeeper ProgramKeeper, extraMounts ...Mount,
) (stdfs.FS, error) {
	mounts := append(StandardMounts(ctx, wasmKeeper, programKeeper), extraMounts...)

	pathFS := logicvfs.New()
	for _, m := range mounts {
		if err := pathFS.Mount(m.Path, m.FS); err != nil {
			return nil, fmt.Errorf("failed to mount %s: %w", m.Path, err)
		}
	}

	return pathFS, nil
}
