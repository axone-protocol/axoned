package embedded

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmbeddedFSReadFile(t *testing.T) {
	Convey("Given an embedded filesystem wrapper", t, func() {
		wrapped := NewFS(fstest.MapFS{
			"foo.pl":        &fstest.MapFile{Data: []byte("foo.")},
			"nested/bar.pl": &fstest.MapFile{Data: []byte("bar.")},
		})

		Convey("when reading an existing file", func() {
			content, err := wrapped.ReadFile("foo.pl")

			So(err, ShouldBeNil)
			So(content, ShouldResemble, []byte("foo."))
		})

		Convey("when reading with a leading slash", func() {
			content, err := wrapped.ReadFile("/nested/bar.pl")

			So(err, ShouldBeNil)
			So(content, ShouldResemble, []byte("bar."))
		})

		Convey("when reading a missing file", func() {
			_, err := wrapped.ReadFile("missing.pl")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading a path escaping root", func() {
			_, err := wrapped.ReadFile("../secret.pl")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func TestEmbeddedFSOpen(t *testing.T) {
	Convey("Given an embedded filesystem wrapper", t, func() {
		wrapped := NewFS(fstest.MapFS{
			"foo.pl": &fstest.MapFile{Data: []byte("foo.")},
		})

		Convey("when opening an existing file", func() {
			file, err := wrapped.Open("foo.pl")

			So(err, ShouldBeNil)
			defer file.Close()

			stat, err := file.Stat()
			So(err, ShouldBeNil)
			So(stat.Name(), ShouldEqual, "foo.pl")
		})

		Convey("when opening a missing file", func() {
			_, err := wrapped.Open("missing.pl")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})
}
