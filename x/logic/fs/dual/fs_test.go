package dual

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"

	. "github.com/smartystreets/goconvey/convey"
)

type recordingFS struct {
	openCalls int
	lastOpen  string
}

func (r *recordingFS) Open(name string) (fs.File, error) {
	r.openCalls++
	r.lastOpen = name
	return fstest.MapFS{
		"ok": &fstest.MapFile{Data: []byte("ok")},
	}.Open("ok")
}

type recordingOpenFileFS struct {
	recordingFS
	openFileCalls int
	lastOpenFile  string
	lastFlag      int
	lastPerm      fs.FileMode
}

func (r *recordingOpenFileFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	r.openFileCalls++
	r.lastOpenFile = name
	r.lastFlag = flag
	r.lastPerm = perm
	return fstest.MapFS{
		"ok": &fstest.MapFile{Data: []byte("ok")},
	}.Open("ok")
}

type recordingReadFileFS struct {
	recordingFS
	readFileCalls int
	lastReadFile  string
}

func (r *recordingReadFileFS) ReadFile(name string) ([]byte, error) {
	r.readFileCalls++
	r.lastReadFile = name
	return []byte("ok"), nil
}

func TestDualFSOpenDispatch(t *testing.T) {
	Convey("Given a dual stack FS", t, func() {
		pathFS := &recordingFS{}
		legacyFS := &recordingFS{}
		dualFS := NewFS(pathFS, legacyFS)

		Convey("when opening a path source-sink", func() {
			_, err := dualFS.Open("/v1/sys/block")

			So(err, ShouldBeNil)
			So(pathFS.openCalls, ShouldEqual, 1)
			So(pathFS.lastOpen, ShouldEqual, "/v1/sys/block")
			So(legacyFS.openCalls, ShouldEqual, 0)
		})

		Convey("when opening a relative path source-sink", func() {
			_, err := dualFS.Open("v1/sys/block")

			So(err, ShouldBeNil)
			So(pathFS.openCalls, ShouldEqual, 1)
			So(pathFS.lastOpen, ShouldEqual, "v1/sys/block")
			So(legacyFS.openCalls, ShouldEqual, 0)
		})

		Convey("when opening a URI source-sink", func() {
			_, err := dualFS.Open("cosmwasm:storage:addr?query=foo")

			So(err, ShouldBeNil)
			So(pathFS.openCalls, ShouldEqual, 0)
			So(legacyFS.openCalls, ShouldEqual, 1)
			So(legacyFS.lastOpen, ShouldEqual, "cosmwasm:storage:addr?query=foo")
		})
	})
}

func TestDualFSOpenFileDispatch(t *testing.T) {
	Convey("Given a dual stack FS where both branches support OpenFile", t, func() {
		pathFS := &recordingOpenFileFS{}
		legacyFS := &recordingOpenFileFS{}
		dualFS := NewFS(pathFS, legacyFS)

		Convey("when opening a path source-sink in OpenFile mode", func() {
			ofs, ok := dualFS.(openFileFS)
			So(ok, ShouldBeTrue)

			_, err := ofs.OpenFile("/v1/codec/bech32", 42, 0o640)

			So(err, ShouldBeNil)
			So(pathFS.openFileCalls, ShouldEqual, 1)
			So(pathFS.lastOpenFile, ShouldEqual, "/v1/codec/bech32")
			So(pathFS.lastFlag, ShouldEqual, 42)
			So(pathFS.lastPerm, ShouldEqual, 0o640)
			So(legacyFS.openFileCalls, ShouldEqual, 0)
		})

		Convey("when opening a URI source-sink in OpenFile mode", func() {
			ofs, ok := dualFS.(openFileFS)
			So(ok, ShouldBeTrue)

			_, err := ofs.OpenFile("cosmwasm:storage:addr?query=foo", 84, 0o600)

			So(err, ShouldBeNil)
			So(pathFS.openFileCalls, ShouldEqual, 0)
			So(legacyFS.openFileCalls, ShouldEqual, 1)
			So(legacyFS.lastOpenFile, ShouldEqual, "cosmwasm:storage:addr?query=foo")
			So(legacyFS.lastFlag, ShouldEqual, 84)
			So(legacyFS.lastPerm, ShouldEqual, 0o600)
		})
	})

	Convey("Given a dual stack FS where path branch does not support OpenFile", t, func() {
		pathFS := &recordingFS{}
		legacyFS := &recordingOpenFileFS{}
		dualFS := NewFS(pathFS, legacyFS)

		Convey("when opening a path source-sink in OpenFile mode", func() {
			ofs, ok := dualFS.(openFileFS)
			So(ok, ShouldBeTrue)

			_, err := ofs.OpenFile("/v1/codec/bech32", 42, 0o640)

			So(err, ShouldNotBeNil)
			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Op, ShouldEqual, "open")
			So(pathErr.Path, ShouldEqual, "/v1/codec/bech32")
			So(errors.Is(err, errors.ErrUnsupported), ShouldBeTrue)
		})
	})

	Convey("Given a dual stack FS where legacy branch does not support OpenFile", t, func() {
		pathFS := &recordingOpenFileFS{}
		legacyFS := &recordingFS{}
		dualFS := NewFS(pathFS, legacyFS)

		Convey("when opening a URI source-sink in OpenFile mode", func() {
			ofs, ok := dualFS.(openFileFS)
			So(ok, ShouldBeTrue)

			_, err := ofs.OpenFile("cosmwasm:storage:addr?query=foo", 84, 0o600)

			So(err, ShouldNotBeNil)
			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Op, ShouldEqual, "open")
			So(pathErr.Path, ShouldEqual, "cosmwasm:storage:addr?query=foo")
			So(errors.Is(err, errors.ErrUnsupported), ShouldBeTrue)
		})
	})
}

func TestDualFSReadFileDispatch(t *testing.T) {
	Convey("Given a dual stack FS where both branches support ReadFile", t, func() {
		pathFS := &recordingReadFileFS{}
		legacyFS := &recordingReadFileFS{}
		dualFS := NewFS(pathFS, legacyFS)

		Convey("when reading a path source-sink", func() {
			rfs, ok := dualFS.(fs.ReadFileFS)
			So(ok, ShouldBeTrue)

			data, err := rfs.ReadFile("/v1/sys/block")

			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte("ok"))
			So(pathFS.readFileCalls, ShouldEqual, 1)
			So(pathFS.lastReadFile, ShouldEqual, "/v1/sys/block")
			So(pathFS.openCalls, ShouldEqual, 0)
			So(legacyFS.readFileCalls, ShouldEqual, 0)
		})

		Convey("when reading a URI source-sink", func() {
			rfs, ok := dualFS.(fs.ReadFileFS)
			So(ok, ShouldBeTrue)

			data, err := rfs.ReadFile("cosmwasm:storage:addr?query=foo")

			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte("ok"))
			So(pathFS.readFileCalls, ShouldEqual, 0)
			So(legacyFS.readFileCalls, ShouldEqual, 1)
			So(legacyFS.lastReadFile, ShouldEqual, "cosmwasm:storage:addr?query=foo")
		})
	})

	Convey("Given a dual stack FS where path branch does not support ReadFile", t, func() {
		pathFS := &recordingFS{}
		legacyFS := &recordingReadFileFS{}
		dualFS := NewFS(pathFS, legacyFS)

		Convey("when reading a path source-sink", func() {
			rfs, ok := dualFS.(fs.ReadFileFS)
			So(ok, ShouldBeTrue)

			data, err := rfs.ReadFile("/v1/sys/block")

			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte("ok"))
			So(pathFS.openCalls, ShouldEqual, 1)
			So(pathFS.lastOpen, ShouldEqual, "/v1/sys/block")
		})
	})
}

func TestIsURI(t *testing.T) {
	Convey("Given URI classification inputs", t, func() {
		cases := []struct {
			in   string
			want bool
		}{
			{in: "cosmwasm:storage:addr?query=foo", want: true},
			{in: "/v1/sys/block", want: false},
			{in: "v1/sys/block", want: false},
			{in: "https://axone.xyz", want: true},
			{in: "% %", want: false},
		}

		for i, tc := range cases {
			Convey(fmt.Sprintf("when case #%d is checked", i), func() {
				So(isURI(tc.in), ShouldEqual, tc.want)
			})
		}
	})
}
