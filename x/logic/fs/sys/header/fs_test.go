package header

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSysHeaderVFSReadFile(t *testing.T) {
	Convey("Given a sys/header VFS", t, func() {
		headerInfo := coreheader.Info{
			Height:  42,
			Hash:    []byte{1, 2, 3},
			Time:    time.Date(2024, 4, 10, 10, 44, 27, 0, time.UTC),
			ChainID: "axone-testchain-1",
			AppHash: []byte{4, 5, 6},
		}
		vfs := NewFS(newTestContext(headerInfo))

		cases := []struct {
			path string
			want []byte
		}{
			{path: "@", want: []byte("header{app_hash:[4,5,6],chain_id:'axone-testchain-1',hash:[1,2,3],height:42,time:1712745867}.\n")},
			{path: "height", want: []byte("42.\n")},
			{path: "hash", want: []byte("[1,2,3].\n")},
			{path: "time", want: []byte("1712745867.\n")},
			{path: "chain_id", want: []byte("'axone-testchain-1'.\n")},
			{path: "app_hash", want: []byte("[4,5,6].\n")},
		}

		for i, tc := range cases {
			Convey(fmt.Sprintf("when reading case #%d path %s", i, tc.path), func() {
				got, err := vfs.ReadFile(tc.path)

				So(err, ShouldBeNil)
				So(got, ShouldResemble, tc.want)
			})
		}
	})
}

func TestSysHeaderVFSOpen(t *testing.T) {
	Convey("Given a sys/header VFS", t, func() {
		headerTime := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
		headerInfo := coreheader.Info{Time: headerTime}
		vfs := NewFS(newTestContext(headerInfo))

		Convey("when opening time", func() {
			f, err := vfs.Open("time")

			So(err, ShouldBeNil)
			defer f.Close()

			info, err := f.Stat()
			So(err, ShouldBeNil)
			So(info.Name(), ShouldEqual, "time")
			So(info.ModTime(), ShouldEqual, headerTime)

			content, err := io.ReadAll(f)
			So(err, ShouldBeNil)
			So(content, ShouldResemble, []byte(fmt.Sprintf("%d.\n", headerTime.Unix())))
		})
	})
}

func TestSysHeaderVFSErrors(t *testing.T) {
	Convey("Given a sys/header VFS", t, func() {
		vfs := NewFS(newTestContext(coreheader.Info{}))

		Convey("when reading an unknown path", func() {
			_, err := vfs.ReadFile("unknown")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading a path escaping root", func() {
			_, err := vfs.ReadFile("../height")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func newTestContext(headerInfo coreheader.Info) sdk.Context {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	return sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger()).WithHeaderInfo(headerInfo)
}
