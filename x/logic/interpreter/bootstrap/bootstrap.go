package bootstrap

import _ "embed"

//go:embed bootstrap.pl
var bootstrap string

// Bootstrap returns the default bootstrap program.
func Bootstrap() string {
	return bootstrap
}
