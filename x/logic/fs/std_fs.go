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

// StandardMounts returns the common logic runtime mounts shared across app and tests.
func StandardMounts(ctx goctx.Context, wasmKeeper logictypes.WasmKeeper, programKeeper ProgramKeeper) []Mount {
	mounts := make([]Mount, 0, 7)
	mounts = append(mounts,
		Mount{Path: "/v1/lib", FS: logicembeddedfs.NewFS(logiclib.Files)},
		Mount{Path: "/v1/sys/header", FS: logicsysheader.NewFS(ctx)},
		Mount{Path: "/v1/sys/comet", FS: logicsyscomet.NewFS(ctx)},
		Mount{Path: "/v1/state/bank", FS: logicbank.NewFS(ctx)},
		Mount{Path: "/v1/dev/codec", FS: logiccodec.NewFS(ctx)},
		Mount{Path: "/v1/dev/wasm", FS: logicwasm.NewFS(ctx, wasmKeeper)},
		Mount{
			Path: "/v1/usr/share/logic",
			FS:   logicshare.NewFS(ctx, programKeeper),
		},
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
