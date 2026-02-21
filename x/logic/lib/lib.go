package lib

import "embed"

// Files contains embedded Prolog library files mounted under /v1/lib.
//
//go:embed .keepme
var Files embed.FS
