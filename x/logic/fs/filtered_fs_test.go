package fs

import (
	"fmt"
	"io/fs"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/samber/lo"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/axone/axoned/v7/x/logic/testutil"
	"github.com/axone/axoned/v7/x/logic/util"
)

func TestSourceFile(t *testing.T) {
	Convey("Given test cases", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			whitelist []string
			blacklist []string
			file      string
			wantError error
		}{
			{
				whitelist: []string{"file1", "file2"},
				blacklist: []string{"file3"},
				file:      "file1",
				wantError: nil,
			},
			{
				whitelist: []string{"tel:"},
				blacklist: []string{},
				file:      "tel:123456789",
				wantError: nil,
			},
			{
				whitelist: []string{"file2"},
				blacklist: []string{"file3"},
				file:      "file1",
				wantError: &fs.PathError{
					Op:   "open",
					Path: "file1",
					Err:  fs.ErrPermission,
				},
			},
			{
				whitelist: []string{},
				blacklist: []string{"file2"},
				file:      "file2",
				wantError: &fs.PathError{
					Op:   "open",
					Path: "file2",
					Err:  fs.ErrPermission,
				},
			},
			{
				whitelist: []string{},
				blacklist: []string{},
				file:      "https://foo{bar/",
				wantError: &fs.PathError{
					Op:   "open",
					Path: "https://foo{bar/",
					Err:  &url.Error{Op: "parse", URL: "https://foo{bar/", Err: url.InvalidHostError("{")},
				},
			},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the test case #%d - file %s", nc, tc.file), func() {
				Convey("and a mocked file system", func() {
					mockedFS := testutil.NewMockFS(ctrl)
					mockedFS.EXPECT().Open(tc.file).Times(lo.If(util.IsNil(tc.wantError), 1).Else(0)).
						DoAndReturn(func(file string) (VirtualFile, error) {
							return NewVirtualFile(
								[]byte("42"),
								util.ParseURLMust(file),
								time.Unix(1681389446, 0)), nil
						})
					Convey("and a filtered file system under test", func() {
						filteredFS := NewFilteredFS(
							lo.Map(tc.whitelist, util.Indexed(util.ParseURLMust)),
							lo.Map(tc.blacklist, util.Indexed(util.ParseURLMust)),
							mockedFS)

						Convey(fmt.Sprintf(`When the open("%s") is called`, tc.file), func() {
							result, err := filteredFS.Open(tc.file)

							Convey("Then the result should be as expected", func() {
								if util.IsNil(tc.wantError) {
									So(err, ShouldBeNil)

									stat, _ := result.Stat()
									So(stat.Name(), ShouldEqual, tc.file)
									So(stat.Size(), ShouldEqual, 2)
									So(stat.ModTime(), ShouldEqual, time.Unix(1681389446, 0))
									So(stat.IsDir(), ShouldBeFalse)
								} else {
									So(err, ShouldNotBeNil)
									So(err, ShouldResemble, tc.wantError)
								}
							})
						})
					})
				})
			})
		}
	})
}
