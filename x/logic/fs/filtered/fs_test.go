package filtered

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/samber/lo"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	fsiface "github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/wasm"
	"github.com/axone-protocol/axoned/v14/x/logic/testutil"
	"github.com/axone-protocol/axoned/v14/x/logic/util"
)

type mockOpenFileFS struct {
	open     func(name string) (fs.File, error)
	readFile func(name string) ([]byte, error)
	openFile func(name string, flag int, perm fs.FileMode) (fs.File, error)
}

func (m mockOpenFileFS) Open(name string) (fs.File, error) {
	return m.open(name)
}

func (m mockOpenFileFS) ReadFile(name string) ([]byte, error) {
	return m.readFile(name)
}

func (m mockOpenFileFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	return m.openFile(name, flag, perm)
}

type mockFS struct {
	open func(name string) (fs.File, error)
}

func (m mockFS) Open(name string) (fs.File, error) {
	return m.open(name)
}

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
					So(err, ShouldEqual, &fs.PathError{Op: "open", Path: "file", Err: fs.ErrInvalid})
				})
			})
		})
	})
}

func TestFilteredVFSOpenFile(t *testing.T) {
	Convey("Given a filtered file system backed by an OpenFile-capable filesystem", t, func() {
		content := []byte("42")
		var gotName string
		var gotFlag int
		var gotPerm fs.FileMode

		filteredFS := NewFS(
			mockOpenFileFS{
				open: func(name string) (fs.File, error) {
					return wasm.NewVirtualFile(name, content, time.Unix(1681389446, 0)), nil
				},
				readFile: func(_ string) ([]byte, error) {
					return content, nil
				},
				openFile: func(name string, flag int, perm fs.FileMode) (fs.File, error) {
					gotName = name
					gotFlag = flag
					gotPerm = perm
					return wasm.NewVirtualFile(name, content, time.Unix(1681389446, 0)), nil
				},
			},
			lo.Map([]string{"file1"}, util.Indexed(util.ParseURLMust)),
			nil,
		)

		Convey("when OpenFile is called on an allowed path", func() {
			ofs, ok := filteredFS.(fsiface.OpenFileFS)
			So(ok, ShouldBeTrue)

			file, err := ofs.OpenFile("file1", 123, 0o640)

			Convey("then it should delegate OpenFile to the underlying filesystem", func() {
				So(err, ShouldBeNil)
				So(file, ShouldNotBeNil)
				So(gotName, ShouldEqual, "file1")
				So(gotFlag, ShouldEqual, 123)
				So(gotPerm, ShouldEqual, 0o640)
			})
		})
	})

	Convey("Given a filtered file system backed by a filesystem without OpenFile support", t, func() {
		filteredFS := NewFS(
			mockFS{
				open: func(name string) (fs.File, error) {
					return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
				},
			},
			nil,
			nil,
		)

		Convey("when OpenFile is called", func() {
			ofs, ok := filteredFS.(fsiface.OpenFileFS)
			So(ok, ShouldBeTrue)

			_, err := ofs.OpenFile("file1", 123, 0o640)

			Convey("then it should return an unsupported operation error", func() {
				So(err, ShouldNotBeNil)
				So(errors.Is(err, errors.ErrUnsupported), ShouldBeTrue)
			})
		})
	})

	Convey("Given a filtered file system with restrictive whitelist", t, func() {
		content := []byte("42")
		filteredFS := NewFS(
			mockOpenFileFS{
				open: func(name string) (fs.File, error) {
					return wasm.NewVirtualFile(name, content, time.Unix(1681389446, 0)), nil
				},
				readFile: func(_ string) ([]byte, error) {
					return content, nil
				},
				openFile: func(name string, _ int, _ fs.FileMode) (fs.File, error) {
					return wasm.NewVirtualFile(name, content, time.Unix(1681389446, 0)), nil
				},
			},
			lo.Map([]string{"file2"}, util.Indexed(util.ParseURLMust)),
			nil,
		)

		Convey("when OpenFile is called on a denied path", func() {
			ofs, ok := filteredFS.(fsiface.OpenFileFS)
			So(ok, ShouldBeTrue)

			_, err := ofs.OpenFile("file1", 123, 0o640)

			Convey("then it should fail with a permission error", func() {
				So(err, ShouldResemble, &fs.PathError{Op: "open", Path: "file1", Err: fs.ErrPermission})
			})
		})
	})
}
