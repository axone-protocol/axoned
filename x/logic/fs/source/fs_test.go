package source

import (
	"errors"
	"io/fs"
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

	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

func TestSourceFS(t *testing.T) {
	Convey("Given a source filesystem", t, func() {
		ctx := newSDKContext()
		ctx = ctx.WithValue(types.SourceFilesProviderContextKey, FilesProvider(func() []string {
			return []string{"/v1/lib/foo.pl", "/v1/lib/bar.pl"}
		}))
		vfs := NewFS(ctx)

		Convey("when reading loaded source files", func() {
			content, err := fs.ReadFile(vfs, "files")
			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "['/v1/lib/foo.pl','/v1/lib/bar.pl'].\n")
		})

		Convey("when opening loaded source files", func() {
			file, err := vfs.Open("files")
			So(err, ShouldBeNil)
			defer file.Close()

			info, err := file.Stat()
			So(err, ShouldBeNil)
			So(info.ModTime(), ShouldEqual, time.Date(2026, time.March, 5, 10, 0, 0, 0, time.UTC))
		})

		Convey("when the provider is missing", func() {
			content, err := fs.ReadFile(NewFS(newSDKContext()), "files")
			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "[].\n")
		})

		Convey("when the path is unknown", func() {
			_, err := fs.ReadFile(vfs, "unknown")
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when the path escapes", func() {
			_, err := fs.ReadFile(vfs, "../files")
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func newSDKContext() sdk.Context {
	stateStore := store.NewCommitMultiStore(dbm.NewMemDB(), log.NewNopLogger(), metrics.NewNoOpMetrics())
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return ctx.WithHeaderInfo(coreheader.Info{
		Height: 42,
		Time:   time.Date(2026, time.March, 5, 10, 0, 0, 0, time.UTC),
	})
}
