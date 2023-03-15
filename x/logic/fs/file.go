package fs

import (
	"bytes"
	"io/fs"
	"net/url"
	"time"
)

type VirtualFile struct {
	reader *bytes.Reader
	info   *VirtualFileInfo
}

type VirtualFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

var (
	_ fs.File     = (*VirtualFile)(nil)
	_ fs.FileInfo = (*VirtualFileInfo)(nil)
)

func NewVirtualFile(src []byte, uri *url.URL, modTime time.Time) VirtualFile {
	reader := bytes.NewReader(src)
	return VirtualFile{
		reader: reader,
		info: &VirtualFileInfo{
			name:    uri.String(),
			size:    reader.Size(),
			modTime: modTime,
		},
	}
}

func (i VirtualFileInfo) Name() string {
	return i.name
}

func (i VirtualFileInfo) Size() int64 {
	return i.size
}

func (i VirtualFileInfo) Mode() fs.FileMode {
	return fs.ModeIrregular
}

func (i VirtualFileInfo) ModTime() time.Time {
	return i.modTime
}

func (i VirtualFileInfo) IsDir() bool {
	return false
}

func (i VirtualFileInfo) Sys() any {
	return nil
}

func (o VirtualFile) Stat() (fs.FileInfo, error) {
	return o.info, nil
}

func (o VirtualFile) Read(b []byte) (int, error) {
	return o.reader.Read(b)
}

func (o VirtualFile) Close() error {
	return nil
}
