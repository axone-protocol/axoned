package composite

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/wasm"
	"github.com/axone-protocol/axoned/v14/x/logic/testutil"
)

type fileSpec struct {
	name        string
	content     []byte
	modTime     time.Time
	isCorrupted bool
}

var (
	vfs1File1 = fileSpec{
		name:    "vfs1:///file1",
		content: []byte("vfs1-file1"),
		modTime: time.Unix(1681389446, 0),
	}

	vfs1File2 = fileSpec{
		name:    "vfs1:///file2",
		content: []byte("vfs1-file2"),
		modTime: time.Unix(1681389446, 0),
	}

	vfs2File1 = fileSpec{
		name:        "vfs2:///file1",
		content:     []byte("vfs1-file1"),
		modTime:     time.Unix(1681389446, 0),
		isCorrupted: true,
	}

	vfs2File2 = fileSpec{
		name:    "vfs2:///file2",
		content: []byte("vfs1-file2"),
		modTime: time.Unix(1681389446, 0),
	}
)

type fsType int

const (
	fsTypeFile = iota
	fsTypeReadFile
)

//nolint:gocognit
func TestFilteredVFS(t *testing.T) {
	Convey("Given test cases", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			files     map[string][]fileSpec
			file      string
			fsType    fsType
			wantError string
			want      *fileSpec
		}{
			{
				files:     map[string][]fileSpec{"vfs1": {vfs1File1, vfs1File2}},
				file:      "vfs1:///file1",
				wantError: "",
				want:      &vfs1File1,
			},
			{
				files:     map[string][]fileSpec{"vfs1": {vfs1File1, vfs1File2}},
				file:      "vfs1:///file1",
				fsType:    fsTypeReadFile,
				wantError: "",
				want:      &vfs1File1,
			},
			{
				files:     map[string][]fileSpec{"vfs1": {vfs1File1, vfs1File2}, "vfs2": {vfs2File2}},
				file:      "vfs2:///file2",
				wantError: "",
				want:      &vfs2File2,
			},
			{
				files:     map[string][]fileSpec{"vfs1": {vfs1File1, vfs1File2}},
				file:      "vfs3:///file1",
				wantError: "vfs3:///file1: file does not exist",
				want:      nil,
			},
			{
				files:     map[string][]fileSpec{"vfs1": {vfs1File1, vfs1File2, vfs2File2}},
				file:      "vfs3:///file1",
				wantError: "vfs3:///file1: file does not exist",
				want:      nil,
			},
			{
				files:     map[string][]fileSpec{},
				file:      "% %",
				wantError: "% %: invalid argument",
				want:      nil,
			},
			{
				files:     map[string][]fileSpec{},
				file:      "foo",
				wantError: "foo: invalid argument",
				want:      nil,
			},
			{
				files:     map[string][]fileSpec{"vfs2": {vfs2File1}},
				file:      "vfs2:///file1",
				fsType:    fsTypeReadFile,
				wantError: "vfs2:///file1: file is corrupted",
				want:      nil,
			},
			{
				files:     map[string][]fileSpec{"vfs2": {vfs2File1}},
				file:      "vfs2:///file1",
				fsType:    fsTypeFile,
				wantError: "vfs2:///file1: file is corrupted",
				want:      nil,
			},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the test case #%d - file %s", nc, tc.file), func() {
				Convey("and a bunch of mocked file systems", func() {
					vfss := make(map[string]fs.FS)
					for k, files := range tc.files {
						switch tc.fsType {
						case fsTypeFile:
							vfs := testutil.NewMockFS(ctrl)
							for _, f := range files {
								registerFileToFS(vfs, f.name, f.content, f.modTime, f.isCorrupted)
							}
							vfss[k] = vfs
						case fsTypeReadFile:
							vfs := testutil.NewMockReadFileFS(ctrl)
							for _, f := range files {
								registerFileToReadFileFS(vfs, f.name, f.content, f.modTime, f.isCorrupted)
							}
							vfss[k] = vfs
						default:
							t.Error("Unsupported fs type")
						}
					}

					Convey("and a composite file system under test", func() {
						compositeFS := NewFS()
						for mountPoint, vfs := range vfss {
							compositeFS.Mount(mountPoint, vfs)
						}

						Convey(fmt.Sprintf(`when the open("%s") is called`, tc.file), func() {
							result, err := compositeFS.Open(tc.file)

							Convey("then the result should be as expected", func() {
								if tc.wantError == "" {
									So(err, ShouldBeNil)

									stat, _ := result.Stat()
									So(stat.Name(), ShouldEqual, tc.file)
									So(stat.Size(), ShouldEqual, int64(len(tc.want.content)))
									So(stat.ModTime(), ShouldEqual, tc.want.modTime)
									So(stat.IsDir(), ShouldBeFalse)
								} else {
									So(err, ShouldNotBeNil)
									So(err.Error(), ShouldEqual, fmt.Sprintf("open %s", tc.wantError))
								}
							})
						})

						Convey(fmt.Sprintf(`when the readFile("%s") is called`, tc.file), func() {
							result, err := compositeFS.ReadFile(tc.file)

							Convey("Then the result should be as expected", func() {
								if tc.wantError == "" {
									So(err, ShouldBeNil)
									So(result, ShouldResemble, tc.want.content)
								} else {
									So(err, ShouldNotBeNil)
									So(err.Error(), ShouldEqual, fmt.Sprintf("readfile %s", tc.wantError))
								}
							})
						})

						Convey(`when the ListMounts() is called`, func() {
							result := compositeFS.ListMounts()

							Convey("Then the result should be as expected", func() {
								vfssNames := make([]string, 0, len(vfss))
								for k := range vfss {
									vfssNames = append(vfssNames, k)
								}
								sort.Strings(vfssNames)
								So(result, ShouldResemble, vfssNames)
							})
						})
					})
				})
			})
		}
	})
}

func registerFileToFS(vfs *testutil.MockFS, name string, content []byte, modTime time.Time, corrupted bool) {
	vfs.EXPECT().Open(name).AnyTimes().
		DoAndReturn(func(file string) (fs.File, error) {
			if corrupted {
				return nil, &fs.PathError{Op: "open", Path: name, Err: fmt.Errorf("file is corrupted")}
			}
			return wasm.NewVirtualFile(
				file,
				content,
				modTime), nil
		})
}

func registerFileToReadFileFS(vfs *testutil.MockReadFileFS, name string, content []byte, modTime time.Time, corrupted bool) {
	vfs.EXPECT().Open(name).AnyTimes().
		DoAndReturn(func(file string) (fs.File, error) {
			if corrupted {
				return nil, &fs.PathError{Op: "open", Path: name, Err: fmt.Errorf("file is corrupted")}
			}
			return wasm.NewVirtualFile(
				file,
				content,
				modTime), nil
		})
	vfs.EXPECT().ReadFile(name).AnyTimes().
		DoAndReturn(func(_ string) ([]byte, error) {
			if corrupted {
				return nil, &fs.PathError{Op: "open", Path: name, Err: fmt.Errorf("file is corrupted")}
			}
			return content, nil
		})
}

func TestGetUnderlyingError(t *testing.T) {
	Convey("Given an error", t, func() {
		Convey("when the error is fs.PathError", func() {
			underlyingErr := errors.New("underlying error")
			pathErr := &fs.PathError{Err: underlyingErr}

			Convey("then getUnderlyingError should return the underlying error", func() {
				So(getUnderlyingError(pathErr), ShouldEqual, underlyingErr)
			})
		})

		Convey("when the error is not fs.PathError", func() {
			genericErr := errors.New("generic error")

			Convey("then getUnderlyingError should return the error itself", func() {
				So(getUnderlyingError(genericErr), ShouldEqual, genericErr)
			})
		})
	})
}
