package sys

import (
	"bytes"
	"io/fs"
	"time"
)

type file struct {
	reader *bytes.Reader
	info   *fileInfo
}

type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

var (
	_ fs.File     = (*file)(nil)
	_ fs.FileInfo = (*fileInfo)(nil)
)

func NewVirtualFile(name string, content []byte, modTime time.Time) fs.File {
	reader := bytes.NewReader(content)
	return &file{
		reader: reader,
		info: &fileInfo{
			name:    name,
			size:    int64(len(content)),
			modTime: modTime,
		},
	}
}

func (i fileInfo) Name() string {
	return i.name
}

func (i fileInfo) Size() int64 {
	return i.size
}

func (i fileInfo) Mode() fs.FileMode {
	return fs.ModeIrregular
}

func (i fileInfo) ModTime() time.Time {
	return i.modTime
}

func (i fileInfo) IsDir() bool {
	return false
}

func (i fileInfo) Sys() any {
	return nil
}

func (o file) Stat() (fs.FileInfo, error) {
	return o.info, nil
}

func (o file) Read(b []byte) (int, error) {
	return o.reader.Read(b)
}

func (o file) Close() error {
	return nil
}
