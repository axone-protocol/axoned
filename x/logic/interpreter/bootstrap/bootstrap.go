package bootstrap

import (
	_ "embed"
)

//go:embed bootstrap.pl
var bootstrap string

//go:embed stdlib.pl
var stdlib string

// Bootstrap returns the default bootstrap program.
func Bootstrap() string {
	return bootstrap + "\n\n" + stdlib
}
