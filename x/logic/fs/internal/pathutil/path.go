package pathutil

import (
	"errors"
	"io/fs"
	"path"
	"strings"
)

// NormalizeSubpath normalizes an input path into a relative subpath safe to pass to an fs.FS.
// It accepts absolute and relative paths, rejects escaping traversal, and returns "." for root-like inputs.
func NormalizeSubpath(name string) (string, error) {
	trimmed := strings.TrimPrefix(name, "/")
	if trimmed == "" || trimmed == "." {
		return ".", nil
	}

	for _, segment := range strings.Split(trimmed, "/") {
		if segment == ".." {
			return "", fs.ErrPermission
		}
	}

	cleaned := path.Clean(trimmed)
	if cleaned == "." {
		return ".", nil
	}

	if !fs.ValidPath(cleaned) {
		return "", fs.ErrInvalid
	}

	return cleaned, nil
}

// UnwrapPathError returns the underlying error when err is an *fs.PathError.
func UnwrapPathError(err error) error {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return pathErr.Err
	}

	return err
}
