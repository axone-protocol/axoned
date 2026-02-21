package lib

import "embed"

// Files contains embedded Prolog library files mounted under /v1/lib.
//
//go:embed *.pl
var Files embed.FS

//go:embed stdlib.pl
var stdlib string

// Stdlib returns the auto-loaded Prolog stdlib program.
func Stdlib() string {
	return stdlib
}
