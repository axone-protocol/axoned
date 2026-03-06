package vfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"
	"time"

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
	openFileErr      error
}

func (r *recordingOpenFileFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	r.lastOpenFilePath = name
	r.lastFlag = flag
	r.lastPerm = perm
	if r.openFileErr != nil {
		return nil, r.openFileErr
	}
	return fstest.MapFS{
		"ok": &fstest.MapFile{Data: []byte("ok")},
	}.Open("ok")
}

type runtimeErrorOpenFileFS struct{}

func (f *runtimeErrorOpenFileFS) Open(name string) (fs.File, error) {
	return &runtimeErrorFile{subpath: name}, nil
}

func (f *runtimeErrorOpenFileFS) OpenFile(name string, _ int, _ fs.FileMode) (fs.File, error) {
	return &runtimeErrorFile{subpath: name}, nil
}

type runtimeErrorFile struct {
	subpath string
}

func (f *runtimeErrorFile) Stat() (fs.FileInfo, error) {
	return runtimeErrorFileInfo{}, nil
}

func (f *runtimeErrorFile) Read(_ []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: f.subpath, Err: fs.ErrPermission}
}

func (f *runtimeErrorFile) Write(_ []byte) (int, error) {
	return 0, &fs.PathError{Op: "write", Path: f.subpath, Err: fs.ErrPermission}
}

func (f *runtimeErrorFile) Close() error {
	return nil
}

type runtimeErrorFileInfo struct{}

func (runtimeErrorFileInfo) Name() string       { return "runtime" }
func (runtimeErrorFileInfo) Size() int64        { return 0 }
func (runtimeErrorFileInfo) Mode() fs.FileMode  { return 0o644 }
func (runtimeErrorFileInfo) ModTime() time.Time { return time.Unix(0, 0) }
func (runtimeErrorFileInfo) IsDir() bool        { return false }
func (runtimeErrorFileInfo) Sys() any           { return nil }

type readFileCapableFS struct {
	data         []byte
	readFileErr  error
	lastReadFile string
}

func (r *readFileCapableFS) Open(_ string) (fs.File, error) {
	return fstest.MapFS{
		"ok": &fstest.MapFile{Data: []byte("ok")},
	}.Open("ok")
}

func (r *readFileCapableFS) ReadFile(name string) ([]byte, error) {
	r.lastReadFile = name
	if r.readFileErr != nil {
		return nil, r.readFileErr
	}

	return append([]byte(nil), r.data...), nil
}

type openOnlyFS struct {
	openErr  error
	openFile fs.File
	lastOpen string
}

func (o *openOnlyFS) Open(name string) (fs.File, error) {
	o.lastOpen = name
	if o.openErr != nil {
		return nil, o.openErr
	}

	return o.openFile, nil
}

type flakyFile struct {
	statErr  error
	readErr  error
	closeErr error
	content  []byte
	pos      int
}

func (f *flakyFile) Stat() (fs.FileInfo, error) {
	return runtimeErrorFileInfo{}, f.statErr
}

func (f *flakyFile) Read(p []byte) (int, error) {
	if f.readErr != nil {
		return 0, f.readErr
	}

	if f.pos >= len(f.content) {
		return 0, io.EOF
	}

	n := copy(p, f.content[f.pos:])
	f.pos += n

	return n, nil
}

func (f *flakyFile) Close() error {
	return f.closeErr
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

		Convey("when routing exact mount path /v1/sys", func() {
			mounted, subpath, err := v.Route("/v1/sys")

			So(err, ShouldBeNil)
			So(mounted, ShouldEqual, sys)
			So(subpath, ShouldEqual, ".")
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

	Convey("Given mounts with same prefix length", t, func() {
		v := New()
		aa := &recordingFS{}
		ab := &recordingFS{}
		So(v.Mount("/aa", aa), ShouldBeNil)
		So(v.Mount("/ab", ab), ShouldBeNil)

		Convey("when routing /ab/x", func() {
			mounted, subpath, err := v.Route("/ab/x")

			So(err, ShouldBeNil)
			So(mounted, ShouldEqual, ab)
			So(subpath, ShouldEqual, "x")
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

		Convey("when mounted Open fails", func() {
			bad := &openOnlyFS{openErr: errors.New("open failed")}
			So(v.Mount("/v1/bad", bad), ShouldBeNil)

			_, err := v.Open("/v1/bad/block")
			So(err, ShouldNotBeNil)
			So(errors.Is(err, errors.New("open failed")), ShouldBeFalse)

			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, "/v1/bad/block")
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

		Convey("when mounted OpenFile fails", func() {
			sys.openFileErr = &fs.PathError{Op: "open", Path: "block", Err: fs.ErrPermission}
			_, err := v.OpenFile("/v1/sys/block", 42, 0o640)

			So(err, ShouldNotBeNil)
			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, "/v1/sys/block")
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})

	Convey("Given a VFS with mount not implementing OpenFile", t, func() {
		v := New()
		So(v.Mount("/v1/sys", &recordingFS{}), ShouldBeNil)

		Convey("when opening /v1/sys/block in OpenFile mode", func() {
			_, err := v.OpenFile("/v1/sys/block", 42, 0o640)

			So(err, ShouldNotBeNil)
			So(errors.Is(err, errors.ErrUnsupported), ShouldBeTrue)
		})
	})

	Convey("Given a VFS with no matching mount for OpenFile", t, func() {
		v := New()
		_, err := v.OpenFile("/v1/sys/block", 42, 0o640)

		So(err, ShouldNotBeNil)
		So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
	})

	Convey("Given a VFS with a mount returning runtime PathError values", t, func() {
		v := New()
		So(v.Mount("/v1/dev", &runtimeErrorOpenFileFS{}), ShouldBeNil)

		Convey("when reading from an opened file", func() {
			f, err := v.OpenFile("/v1/dev/block", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			_, err = f.Read(make([]byte, 1))
			So(err, ShouldNotBeNil)

			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, "/v1/dev/block")
		})

		Convey("when writing to an opened file", func() {
			f, err := v.OpenFile("/v1/dev/block", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			writer, ok := f.(interface{ Write([]byte) (int, error) })
			So(ok, ShouldBeTrue)

			_, err = writer.Write([]byte("x"))
			So(err, ShouldNotBeNil)

			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, "/v1/dev/block")
		})
	})
}

func TestFileSystemReadFileDispatch(t *testing.T) {
	Convey("Given a VFS with ReadFile-capable mount", t, func() {
		v := New()
		sys := &readFileCapableFS{data: []byte("payload")}
		So(v.Mount("/v1/sys", sys), ShouldBeNil)

		Convey("when reading /v1/sys/block", func() {
			content, err := v.ReadFile("/v1/sys/block")

			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "payload")
			So(sys.lastReadFile, ShouldEqual, "block")
		})

		Convey("when mounted ReadFile fails", func() {
			sys.readFileErr = &fs.PathError{Op: "open", Path: "block", Err: fs.ErrPermission}
			_, err := v.ReadFile("/v1/sys/block")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)

			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, "/v1/sys/block")
		})
	})

	Convey("Given a VFS with Open-only mount", t, func() {
		v := New()
		So(v.Mount("/v1/sys", &openOnlyFS{openFile: &flakyFile{content: []byte("ok")}}), ShouldBeNil)

		Convey("when reading through Open + io.ReadAll", func() {
			content, err := v.ReadFile("/v1/sys/block")

			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "ok")
		})
	})

	Convey("Given a VFS where Open fails for ReadFile fallback", t, func() {
		v := New()
		So(v.Mount("/v1/sys", &openOnlyFS{openErr: fs.ErrNotExist}), ShouldBeNil)

		Convey("when reading /v1/sys/block", func() {
			_, err := v.ReadFile("/v1/sys/block")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})

	Convey("Given a VFS where fallback reader fails", t, func() {
		v := New()
		So(v.Mount("/v1/sys", &openOnlyFS{openFile: &flakyFile{readErr: errors.New("boom")}}), ShouldBeNil)

		Convey("when reading /v1/sys/block", func() {
			_, err := v.ReadFile("/v1/sys/block")

			So(err, ShouldNotBeNil)
			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, "/v1/sys/block")
		})
	})

	Convey("Given a VFS with no matching mount", t, func() {
		v := New()
		_, err := v.ReadFile("/v1/sys/block")

		So(err, ShouldNotBeNil)
		So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
	})
}

func TestWrappedFileStatAndClose(t *testing.T) {
	Convey("Given a wrapped file from vfs.Open", t, func() {
		v := New()
		So(v.Mount("/v1/sys", &openOnlyFS{openFile: &flakyFile{content: []byte("ok")}}), ShouldBeNil)

		f, err := v.Open("/v1/sys/block")
		So(err, ShouldBeNil)

		Convey("when calling Stat and Close with nil underlying errors", func() {
			_, err := f.Stat()
			So(err, ShouldBeNil)

			So(f.Close(), ShouldBeNil)
		})
	})

	Convey("Given a wrapped file with stat/close errors", t, func() {
		v := New()
		underlying := &flakyFile{
			statErr:  &fs.PathError{Op: "stat", Path: "block", Err: fs.ErrPermission},
			closeErr: &fs.PathError{Op: "close", Path: "block", Err: fs.ErrPermission},
		}
		So(v.Mount("/v1/sys", &openOnlyFS{openFile: underlying}), ShouldBeNil)

		f, err := v.Open("/v1/sys/block")
		So(err, ShouldBeNil)

		Convey("when calling Stat", func() {
			_, err := f.Stat()
			So(err, ShouldNotBeNil)
			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, "/v1/sys/block")
		})

		Convey("when calling Close", func() {
			err := f.Close()
			So(err, ShouldNotBeNil)
			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, "/v1/sys/block")
		})
	})
}

func TestWrapPathErrorBranches(t *testing.T) {
	Convey("wrapPathError nil branch", t, func() {
		So(wrapPathError("open", "/v1/sys/block", nil), ShouldBeNil)
	})

	Convey("wrapPathError EOF branch", t, func() {
		err := wrapPathError("read", "/v1/sys/block", io.EOF)
		So(errors.Is(err, io.EOF), ShouldBeTrue)
	})
}
