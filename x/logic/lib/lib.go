package lib

import "embed"

// Files contains embedded Prolog library files mounted under /v1/lib.
//
//go:embed *.pl
var Files embed.FS
