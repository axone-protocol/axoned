package sys

import (
	"io/fs"
	"time"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/virtualfile"
)

func NewVirtualFile(name string, content []byte, modTime time.Time) fs.File {
	return virtualfile.New(name, content, modTime)
}
