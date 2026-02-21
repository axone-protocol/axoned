package virtualfile

import (
	"errors"
	"io"
	"io/fs"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVirtualFile(t *testing.T) {
	Convey("Given a virtual file", t, func() {
		modTime := time.Date(2026, 2, 21, 12, 34, 56, 0, time.UTC)
		f := New("foo.pl", []byte("foo."), modTime)

		Convey("when reading metadata", func() {
			info, err := f.Stat()

			So(err, ShouldBeNil)
			So(info.Name(), ShouldEqual, "foo.pl")
			So(info.Size(), ShouldEqual, 4)
			So(info.Mode(), ShouldEqual, fs.FileMode(0o444))
			So(info.Mode().IsRegular(), ShouldBeTrue)
			So(info.ModTime(), ShouldEqual, modTime)
			So(info.IsDir(), ShouldBeFalse)
			So(info.Sys(), ShouldBeNil)
		})

		Convey("when reading content", func() {
			content, err := io.ReadAll(f)

			So(err, ShouldBeNil)
			So(content, ShouldResemble, []byte("foo."))

			_, err = f.Read(make([]byte, 1))
			So(errors.Is(err, io.EOF), ShouldBeTrue)
		})

		Convey("when closing", func() {
			So(f.Close(), ShouldBeNil)
		})
	})
}
