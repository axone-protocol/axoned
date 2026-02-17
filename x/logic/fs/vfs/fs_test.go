package vfs

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"

	. "github.com/smartystreets/goconvey/convey"
)

type recordingFS struct {
	lastOpen string
}

func (r *recordingFS) Open(name string) (fs.File, error) {
	r.lastOpen = name
	return fstest.MapFS{
		"ok": &fstest.MapFile{Data: []byte("ok")},
	}.Open("ok")
}

type recordingOpenFileFS struct {
	recordingFS
	lastOpenFilePath string
	lastFlag         int
	lastPerm         fs.FileMode
}

func (r *recordingOpenFileFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	r.lastOpenFilePath = name
	r.lastFlag = flag
	r.lastPerm = perm
	return fstest.MapFS{
		"ok": &fstest.MapFile{Data: []byte("ok")},
	}.Open("ok")
}

func TestNormalizePath(t *testing.T) {
	Convey("Given normalization cases", t, func() {
		cases := []struct {
			name    string
			in      string
			want    string
			wantErr error
		}{
			{
				name: "relative path",
				in:   "v1/sys/block",
				want: "/v1/sys/block",
			},
			{
				name: "dotdot collapse",
				in:   "/v1/sys/../params/max_tx",
				want: "/v1/params/max_tx",
			},
			{
				name: "repeated slash",
				in:   "/v1//sys///block",
				want: "/v1/sys/block",
			},
			{
				name:    "escape root absolute",
				in:      "/../../etc",
				wantErr: fs.ErrPermission,
			},
			{
				name:    "escape root relative",
				in:      "../../etc",
				wantErr: fs.ErrPermission,
			},
			{
				name: "root",
				in:   "/",
				want: "/",
			},
			{
				name: "dot",
				in:   ".",
				want: "/",
			},
			{
				name: "empty",
				in:   "",
				want: "/",
			},
		}

		for i, tc := range cases {
			Convey(fmt.Sprintf("when case #%d (%s) is normalized", i, tc.name), func() {
				got, err := normalizePath(tc.in)
				if tc.wantErr != nil {
					So(err, ShouldNotBeNil)
					So(errors.Is(err, tc.wantErr), ShouldBeTrue)
					return
				}

				So(err, ShouldBeNil)
				So(got, ShouldEqual, tc.want)
			})
		}
	})
}

func TestFileSystemRoute(t *testing.T) {
	Convey("Given a mounted VFS", t, func() {
		v := New()
		root := &recordingFS{}
		v1 := &recordingFS{}
		sys := &recordingFS{}

		So(v.Mount("/", root), ShouldBeNil)
		So(v.Mount("/v1", v1), ShouldBeNil)
		So(v.Mount("/v1/sys", sys), ShouldBeNil)

		Convey("when routing /v1/sys/block", func() {
			mounted, subpath, err := v.Route("/v1/sys/block")

			So(err, ShouldBeNil)
			So(mounted, ShouldEqual, sys)
			So(subpath, ShouldEqual, "block")
		})

		Convey("when routing /v1/system", func() {
			mounted, subpath, err := v.Route("/v1/system")

			So(err, ShouldBeNil)
			So(mounted, ShouldEqual, v1)
			So(subpath, ShouldEqual, "system")
		})

		Convey("when routing relative path v1/sys/./block", func() {
			mounted, subpath, err := v.Route("v1/sys/./block")

			So(err, ShouldBeNil)
			So(mounted, ShouldEqual, sys)
			So(subpath, ShouldEqual, "block")
		})

		Convey("when routing /", func() {
			mounted, subpath, err := v.Route("/")

			So(err, ShouldBeNil)
			So(mounted, ShouldEqual, root)
			So(subpath, ShouldEqual, ".")
		})

		Convey("when routing an unmatched path with root mount", func() {
			mounted, subpath, err := v.Route("/v2/sys/block")

			So(err, ShouldBeNil)
			So(mounted, ShouldEqual, root)
			So(subpath, ShouldEqual, "v2/sys/block")
		})

		Convey("when routing a path escaping root", func() {
			_, _, err := v.Route("/../../etc")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})

	Convey("Given a mounted VFS without root mount", t, func() {
		v := New()
		So(v.Mount("/v1/sys", &recordingFS{}), ShouldBeNil)

		Convey("when routing a missing mount", func() {
			_, _, err := v.Route("/v2/sys/block")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})
}

func TestFileSystemMountErrors(t *testing.T) {
	Convey("Given a VFS with /v1/sys mounted", t, func() {
		v := New()
		So(v.Mount("/v1/sys", &recordingFS{}), ShouldBeNil)

		Convey("when mounting the same canonical prefix", func() {
			err := v.Mount("v1/sys", &recordingFS{})

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrExist), ShouldBeTrue)
		})

		Convey("when mounting a nil filesystem", func() {
			err := v.Mount("/v2/sys", nil)

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrInvalid), ShouldBeTrue)
		})

		Convey("when mounting an invalid escaping prefix", func() {
			err := v.Mount("/../../etc", &recordingFS{})

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func TestFileSystemOpenDispatch(t *testing.T) {
	Convey("Given a VFS with /v1/sys mounted", t, func() {
		v := New()
		sys := &recordingFS{}
		So(v.Mount("/v1/sys", sys), ShouldBeNil)

		Convey("when opening v1/sys/block", func() {
			_, err := v.Open("v1/sys/block")

			So(err, ShouldBeNil)
			So(sys.lastOpen, ShouldEqual, "block")
		})

		Convey("when opening a missing mount", func() {
			_, err := v.Open("/v2/sys/block")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})
}

func TestFileSystemOpenFileDispatch(t *testing.T) {
	Convey("Given a VFS with OpenFile-capable mount", t, func() {
		v := New()
		sys := &recordingOpenFileFS{}
		So(v.Mount("/v1/sys", sys), ShouldBeNil)

		Convey("when opening /v1/sys/block in OpenFile mode", func() {
			_, err := v.OpenFile("/v1/sys/block", 42, 0o640)

			So(err, ShouldBeNil)
			So(sys.lastOpenFilePath, ShouldEqual, "block")
			So(sys.lastFlag, ShouldEqual, 42)
			So(sys.lastPerm, ShouldEqual, 0o640)
		})
	})

	Convey("Given a VFS with mount not implementing OpenFile", t, func() {
		v := New()
		So(v.Mount("/v1/sys", &recordingFS{}), ShouldBeNil)

		Convey("when opening /v1/sys/block in OpenFile mode", func() {
			_, err := v.OpenFile("/v1/sys/block", 42, 0o640)

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}
