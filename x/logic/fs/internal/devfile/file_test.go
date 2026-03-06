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
			WithAllowEmptyRequest(false),
			WithCommit(func(request []byte) ([]byte, error) {
				commitCalls++
				So(string(request), ShouldEqual, "ping")
				return []byte("pong"), nil
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
	Convey("Given a half-duplex file that rejects empty requests", t, func() {
		commitCalled := false
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),
			WithAllowEmptyRequest(false),
			WithCommit(func(_ []byte) ([]byte, error) {
				commitCalled = true
				return nil, nil
			}),
		)
		So(err, ShouldBeNil)

		_, err = file.Read(make([]byte, 1))
		So(errors.Is(err, ErrInvalidRequest), ShouldBeTrue)

		_, err = file.Read(make([]byte, 1))
		So(errors.Is(err, ErrInvalidRequest), ShouldBeTrue)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err := writer.Write([]byte("x"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		So(commitCalled, ShouldBeFalse)
	})
}

func TestHalfDuplexAllowEmptyRequest(t *testing.T) {
	Convey("Given a half-duplex file that allows empty requests", t, func() {
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),
			WithAllowEmptyRequest(true),
			WithCommit(func(request []byte) ([]byte, error) {
				So(len(request), ShouldEqual, 0)
				return []byte("ok"), nil
			}),
		)
		So(err, ShouldBeNil)

		buf := make([]byte, 2)
		n, err := file.Read(buf)
		So(err, ShouldBeNil)
		So(string(buf[:n]), ShouldEqual, "ok")
	})
}

func TestHalfDuplexMaxRequestBytes(t *testing.T) {
	Convey("Given a half-duplex file with max request bytes", t, func() {
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),
			WithMaxRequestBytes(4),
			WithAllowEmptyRequest(false),
			WithCommit(func(request []byte) ([]byte, error) {
				return request, nil
			}),
		)
		So(err, ShouldBeNil)

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		So(ok, ShouldBeTrue)

		n, err := writer.Write([]byte("12345"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
	})
}

func TestHalfDuplexClosed(t *testing.T) {
	Convey("Given a closed half-duplex file", t, func() {
		file, err := New(
			WithPath("/v1/dev/test"),
			WithModTime(time.Unix(0, 0).UTC()),
			WithAllowEmptyRequest(false),
			WithCommit(func(request []byte) ([]byte, error) {
				return request, nil
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
			WithAllowEmptyRequest(false),
			WithCommit(func(_ []byte) ([]byte, error) {
				commitCalls++
				return nil, expected
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
			WithAllowEmptyRequest(false),
			WithMaxResponseBytes(2),
			WithCommit(func(_ []byte) ([]byte, error) {
				commitCalls++
				return []byte("too-big"), nil
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
		So(errors.Is(err, ErrResponseTooLarge), ShouldBeTrue)

		n, err = writer.Write([]byte("y"))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)

		n, err = file.Read(make([]byte, 16))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, ErrResponseTooLarge), ShouldBeTrue)

		So(commitCalls, ShouldEqual, 1)
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

func TestHalfDuplexEmptyRequestCustomError(t *testing.T) {
	Convey("Given a half-duplex file with custom empty request error", t, func() {
		customErr := errors.New("custom-empty-request")
		commitCalled := false

		file, err := New(
			WithPath("/v1/dev/test"),
			WithAllowEmptyRequest(false),
			WithEmptyRequestError(customErr),
			WithCommit(func(_ []byte) ([]byte, error) {
				commitCalled = true
				return nil, nil
			}),
		)
		So(err, ShouldBeNil)

		n, err := file.Read(make([]byte, 1))
		So(n, ShouldEqual, 0)
		So(errors.Is(err, customErr), ShouldBeTrue)
		So(commitCalled, ShouldBeFalse)
	})
}

func TestHalfDuplexStat(t *testing.T) {
	Convey("Given a half-duplex file", t, func() {
		modTime := time.Unix(123, 0).UTC()
		file, err := New(
			WithPath("/v1/dev/echo"),
			WithModTime(modTime),
			WithAllowEmptyRequest(true),
			WithCommit(func(_ []byte) ([]byte, error) {
				return nil, nil
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
			tc := tc
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
