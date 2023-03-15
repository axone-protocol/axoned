package fs

import (
	"bytes"
	"io/fs"
	"net/url"
	"time"
)

type Object struct {
	reader *bytes.Reader
	uri    *url.URL
}

func NewObject(src []byte, uri *url.URL) Object {
	return Object{
		reader: bytes.NewReader(src),
		uri:    uri,
	}
}

func (o Object) Name() string {
	return o.uri.String()
}

func (o Object) Size() int64 {
	return o.reader.Size()
}

func (o Object) Mode() fs.FileMode {
	return fs.ModeIrregular
}

func (o Object) ModTime() time.Time {
	return time.Now() // TODO: change time
}

func (o Object) IsDir() bool {
	return false
}

func (o Object) Sys() any {
	return nil
}

func (o Object) Stat() (fs.FileInfo, error) {
	return o, nil
}

func (o Object) Read(b []byte) (int, error) {
	return o.reader.Read(b)
}

func (o Object) Close() error {
	return nil
}
