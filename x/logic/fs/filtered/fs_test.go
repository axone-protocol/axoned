package filtered

import (
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/samber/lo"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/axone-protocol/axoned/v13/x/logic/fs/wasm"
	"github.com/axone-protocol/axoned/v13/x/logic/testutil"
	"github.com/axone-protocol/axoned/v13/x/logic/util"
)

func TestFilteredVFS(t *testing.T) {
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
					Err:  fmt.Errorf("invalid argument"),
				},
			},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the test case #%d - file %s", nc, tc.file), func() {
				Convey("and a mocked file system", func() {
					content := []byte("42")
					mockedFS := testutil.NewMockReadFileFS(ctrl)
					mockedFS.EXPECT().Open(tc.file).AnyTimes().
						DoAndReturn(func(file string) (fs.File, error) {
							return wasm.NewVirtualFile(
								file,
								content,
								time.Unix(1681389446, 0)), nil
						})
					mockedFS.EXPECT().ReadFile(tc.file).AnyTimes().
						DoAndReturn(func(_ string) ([]byte, error) {
							return content, nil
						})
					Convey("and a filtered file system under test", func() {
						filteredFS := NewFS(
							mockedFS,
							lo.Map(tc.whitelist, util.Indexed(util.ParseURLMust)),
							lo.Map(tc.blacklist, util.Indexed(util.ParseURLMust)),
						)

						Convey(fmt.Sprintf(`when the open("%s") is called`, tc.file), func() {
							result, err := filteredFS.Open(tc.file)

							Convey("then the result should be as expected", func() {
								if lo.IsNil(tc.wantError) {
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

						Convey(fmt.Sprintf(`when the readFile("%s") is called`, tc.file), func() {
							result, err := filteredFS.ReadFile(tc.file)

							Convey("Then the result should be as expected", func() {
								if lo.IsNil(tc.wantError) {
									So(err, ShouldBeNil)
									So(result, ShouldResemble, content)
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

	Convey("Given a mocked fs that does not implement ReadFileFS", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedFS := testutil.NewMockFS(ctrl)
		Convey("and a filtered file system under test", func() {
			filteredFS := NewFS(mockedFS, nil, nil)

			Convey("when readFile is called", func() {
				_, err := filteredFS.ReadFile("file")

				Convey("then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err, ShouldEqual, &fs.PathError{Op: "readfile", Path: "file", Err: fs.ErrInvalid})
				})
			})
		})
	})
}
