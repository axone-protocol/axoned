package fs

import (
	"bytes"
	"io/fs"
	"net/url"
	"time"
)

type VirtualFile struct {
	reader  *bytes.Reader
	uri     *url.URL
	modTime time.Time
}

var (
	_ fs.File     = (*VirtualFile)(nil)
	_ fs.FileInfo = (*VirtualFile)(nil)
)

func NewVirtualFile(src []byte, uri *url.URL, modTime time.Time) VirtualFile {
	return VirtualFile{
		reader:  bytes.NewReader(src),
		uri:     uri,
		modTime: modTime,
	}
}

func (o VirtualFile) Name() string {
	return o.uri.String()
}

func (o VirtualFile) Size() int64 {
	return o.reader.Size()
}

func (o VirtualFile) Mode() fs.FileMode {
	return fs.ModeIrregular
}

func (o VirtualFile) ModTime() time.Time {
	return o.modTime
}

func (o VirtualFile) IsDir() bool {
	return false
}

func (o VirtualFile) Sys() any {
	return nil
}

func (o VirtualFile) Stat() (fs.FileInfo, error) {
	return o, nil
}

func (o VirtualFile) Read(b []byte) (int, error) {
	return o.reader.Read(b)
}

func (o VirtualFile) Close() error {
	return nil
}
