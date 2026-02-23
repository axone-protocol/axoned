package lib

import "embed"

// Files contains embedded Prolog library files mounted under /v1/lib.
// These are optional libraries that can be explicitly loaded by user programs.
//
//go:embed *.pl
var Files embed.FS
