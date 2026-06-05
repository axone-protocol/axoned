package pathutil

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testFooPath       = "foo.pl"
	testNestedBarPath = "nested/bar.pl"
)

func TestNormalizeSubpath(t *testing.T) {
	Convey("Given normalization cases", t, func() {
		cases := []struct {
			in      string
			want    string
			wantErr error
		}{
			{in: "", want: "."},
			{in: ".", want: "."},
			{in: "/", want: "."},
			{in: testFooPath, want: testFooPath},
			{in: "/foo.pl", want: testFooPath},
			{in: testNestedBarPath, want: testNestedBarPath},
			{in: "nested/./bar.pl", want: testNestedBarPath},
			{in: "../secret", wantErr: fs.ErrPermission},
			{in: "nested/../secret", wantErr: fs.ErrPermission},
		}

		for i, tc := range cases {
			Convey(fmt.Sprintf("when normalizing case #%d (%q)", i, tc.in), func() {
				got, err := NormalizeSubpath(tc.in)
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

func TestUnwrapPathError(t *testing.T) {
	Convey("Given a path error", t, func() {
		err := &fs.PathError{Op: "open", Path: "foo", Err: fs.ErrNotExist}

		Convey("it should unwrap the underlying error", func() {
			So(UnwrapPathError(err), ShouldEqual, fs.ErrNotExist)
		})
	})
}
