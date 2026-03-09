package devfile

import (
	"errors"
	"io"
	"io/fs"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHalfDuplexLifecycle(t *testing.T) {
	Convey("Given a half-duplex file with a valid request and response", t, func() {
		commitCalls := 0
		file, err := New(
			WithPath("/v1/dev/echo"),
			WithModTime(time.Unix(0, 0).UTC()),

			WithCommit(func(r io.Reader, w io.Writer) error {
				commitCalls++
				request, _ := io.ReadAll(r)
				So(string(request), ShouldEqual, "ping")
				_, err := w.Write([]byte("pong"))
				return err
			}),
		)
		So(err, ShouldBeNil)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err := writer.Write([]byte("ping"))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, 4)

		buf := make([]byte, 8)
		n, err = file.Read(buf)
		So(err, ShouldBeNil)
		So(string(buf[:n]), ShouldEqual, "pong")

		n, err = file.Read(buf)
		So(n, ShouldEqual, 0)
		So(errors.Is(err, io.EOF), ShouldBeTrue)

		n, err = writer.Write([]byte("x"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)

		So(commitCalls, ShouldEqual, 1)
	})
}

func TestHalfDuplexEmptyRequestValidation(t *testing.T) {
	Convey("Given a half-duplex file with empty request", t, func() {
		commitCalled := false
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),
			WithCommit(func(r io.Reader, _ io.Writer) error {
				commitCalled = true
				request, _ := io.ReadAll(r)
				So(len(request), ShouldEqual, 0)
				return nil
			}),
		)
		So(err, ShouldBeNil)

		n, err := file.Read(make([]byte, 1))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, io.EOF), ShouldBeTrue)
		So(commitCalled, ShouldBeTrue)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err = writer.Write([]byte("x"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
	})
}

func TestHalfDuplexMaxRequestBytes(t *testing.T) {
	Convey("Given a half-duplex file with max request bytes", t, func() {
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),
			WithMaxRequestBytes(4),

			WithCommit(func(r io.Reader, w io.Writer) error {
				request, err := io.ReadAll(r)
				if err != nil {
					return err
				}
				_, err = w.Write(request)
				return err
			}),
		)
		So(err, ShouldBeNil)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err := writer.Write([]byte("12345"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, ErrWriteLimit), ShouldBeTrue)
	})
}

func TestHalfDuplexClosed(t *testing.T) {
	Convey("Given a closed half-duplex file", t, func() {
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),

			WithCommit(func(r io.Reader, w io.Writer) error {
				request, err := io.ReadAll(r)
				if err != nil {
					return err
				}
				_, err = w.Write(request)
				return err
			}),
		)
		So(err, ShouldBeNil)
		So(file.Close(), ShouldBeNil)

		n, err := file.Read(make([]byte, 1))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrClosed), ShouldBeTrue)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err = writer.Write([]byte("x"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrClosed), ShouldBeTrue)
	})
}

func TestHalfDuplexCommitErrorPropagation(t *testing.T) {
	Convey("Given a half-duplex file whose commit fails", t, func() {
		expected := errors.New("boom")
		commitCalls := 0
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),

			WithCommit(func(_ io.Reader, _ io.Writer) error {
				commitCalls++
				return expected
			}),
		)
		So(err, ShouldBeNil)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err := writer.Write([]byte("x"))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, 1)

		n, err = file.Read(make([]byte, 1))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, expected), ShouldBeTrue)

		n, err = writer.Write([]byte("y"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)

		n, err = file.Read(make([]byte, 1))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, expected), ShouldBeTrue)

		So(commitCalls, ShouldEqual, 1)
	})
}

func TestHalfDuplexMaxResponseBytes(t *testing.T) {
	Convey("Given a half-duplex file with max response bytes", t, func() {
		commitCalls := 0
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),

			WithMaxResponseBytes(2),
			WithCommit(func(_ io.Reader, w io.Writer) error {
				commitCalls++
				_, err := w.Write([]byte("too-big"))
				return err
			}),
		)
		So(err, ShouldBeNil)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err := writer.Write([]byte("x"))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, 1)

		n, err = file.Read(make([]byte, 16))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, ErrWriteLimit), ShouldBeTrue)

		n, err = writer.Write([]byte("y"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)

		n, err = file.Read(make([]byte, 16))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, ErrWriteLimit), ShouldBeTrue)

		So(commitCalls, ShouldEqual, 1)
	})
}

func TestHalfDuplexTransferHook(t *testing.T) {
	Convey("Given a half-duplex file with transfer instrumentation", t, func() {
		var transfers []struct {
			dir TransferDirection
			n   int
		}

		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),
			WithTransferHook(func(dir TransferDirection, n int) {
				transfers = append(transfers, struct {
					dir TransferDirection
					n   int
				}{dir: dir, n: n})
			}),
			WithCommit(func(r io.Reader, w io.Writer) error {
				request, err := io.ReadAll(r)
				if err != nil {
					return err
				}
				_, err = w.Write(append([]byte("ok:"), request...))
				return err
			}),
		)
		So(err, ShouldBeNil)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err := writer.Write([]byte("ping"))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, 4)

		response, err := io.ReadAll(file)
		So(err, ShouldBeNil)
		So(string(response), ShouldEqual, "ok:ping")
		So(transfers, ShouldResemble, []struct {
			dir TransferDirection
			n   int
		}{
			{dir: TransferRequest, n: 4},
			{dir: TransferResponse, n: 7},
		})
	})
}

func TestHalfDuplexMissingCommit(t *testing.T) {
	Convey("Given a half-duplex file without commit function", t, func() {
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),
		)

		So(file, ShouldBeNil)
		So(errors.Is(err, ErrMissingCommit), ShouldBeTrue)
	})
}

func TestHalfDuplexStat(t *testing.T) {
	Convey("Given a half-duplex file", t, func() {
		modTime := time.Unix(123, 0).UTC()
		file, err := New(
			WithPath("/v1/dev/echo"),
			WithModTime(modTime),
			WithCommit(func(_ io.Reader, _ io.Writer) error {
				return nil
			}),
		)
		So(err, ShouldBeNil)

		info, err := file.Stat()
		So(err, ShouldBeNil)
		So(info.Name(), ShouldEqual, "echo")
		So(info.Size(), ShouldEqual, int64(0))
		So(info.Mode(), ShouldEqual, fs.FileMode(0o666))
		So(info.ModTime(), ShouldEqual, modTime)
		So(info.IsDir(), ShouldBeFalse)
		So(info.Sys(), ShouldBeNil)
	})
}

func TestBaseName(t *testing.T) {
	Convey("Given path variants", t, func() {
		testCases := []struct {
			path string
			want string
		}{
			{path: "/v1/dev/echo", want: "echo"},
			{path: "echo", want: "echo"},
			{path: "/", want: ""},
		}

		for _, tc := range testCases {
			Convey(tc.path, func() {
				So(baseName(tc.path), ShouldEqual, tc.want)
			})
		}
	})
}

func TestPathError(t *testing.T) {
	Convey("Given a halfDuplex file", t, func() {
		f := &halfDuplexFile{cfg: halfDuplexConfig{path: "/v1/dev/test"}}

		Convey("when the input error is nil", func() {
			So(f.pathError("read", nil), ShouldBeNil)
		})

		Convey("when the input error is EOF", func() {
			So(f.pathError("read", io.EOF), ShouldEqual, io.EOF)
		})

		Convey("when the input error is already an fs.PathError", func() {
			pathErr := &fs.PathError{Op: "open", Path: "/other", Err: fs.ErrNotExist}
			So(f.pathError("read", pathErr), ShouldEqual, pathErr)
		})

		Convey("when the input error is not an fs.PathError", func() {
			innerErr := errors.New("boom")
			err := f.pathError("write", innerErr)

			var wrapped *fs.PathError
			So(errors.As(err, &wrapped), ShouldBeTrue)
			So(wrapped.Op, ShouldEqual, "write")
			So(wrapped.Path, ShouldEqual, "/v1/dev/test")
			So(errors.Is(wrapped.Err, innerErr), ShouldBeTrue)
		})
	})
}
