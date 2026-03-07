//nolint:revive // Package name 'fs' conflicts with stdlib io/fs but provides domain-specific filesystem abstractions
package fs

import (
	goctx "context"
	"fmt"
	stdfs "io/fs"

	logiccodec "github.com/axone-protocol/axoned/v14/x/logic/fs/codec"
	logicembeddedfs "github.com/axone-protocol/axoned/v14/x/logic/fs/embedded"
	logicsyscomet "github.com/axone-protocol/axoned/v14/x/logic/fs/sys/comet"
	logicsysheader "github.com/axone-protocol/axoned/v14/x/logic/fs/sys/header"
	logicvfs "github.com/axone-protocol/axoned/v14/x/logic/fs/vfs"
	logiclib "github.com/axone-protocol/axoned/v14/x/logic/lib"
)

// Mount defines a filesystem mounted at an absolute path in the logic VFS.
type Mount struct {
	Path string
	FS   stdfs.FS
}

// StandardMounts returns the common logic runtime mounts shared across app and tests.
func StandardMounts(ctx goctx.Context, bankFS, wasmFS stdfs.FS) []Mount {
	return []Mount{
		{Path: "/v1/lib", FS: logicembeddedfs.NewFS(logiclib.Files)},
		{Path: "/v1/sys/header", FS: logicsysheader.NewFS(ctx)},
		{Path: "/v1/sys/comet", FS: logicsyscomet.NewFS(ctx)},
		{Path: "/v1/state/bank", FS: bankFS},
		{Path: "/v1/dev/codec", FS: logiccodec.NewFS(ctx)},
		{Path: "/v1/dev/wasm", FS: wasmFS},
	}
}

// NewVFS builds a logic virtual filesystem with standard mounts and optional extras.
func NewVFS(ctx goctx.Context, bankFS, wasmFS stdfs.FS, extraMounts ...Mount) (stdfs.FS, error) {
	mounts := append(StandardMounts(ctx, bankFS, wasmFS), extraMounts...)

	pathFS := logicvfs.New()
	for _, m := range mounts {
		if err := pathFS.Mount(m.Path, m.FS); err != nil {
			return nil, fmt.Errorf("failed to mount %s: %w", m.Path, err)
		}
	}

	return pathFS, nil
}
