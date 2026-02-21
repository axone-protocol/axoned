package bootstrap

import (
	_ "embed"

	logiclib "github.com/axone-protocol/axoned/v14/x/logic/lib"
)

//go:embed bootstrap.pl
var bootstrap string

// Bootstrap returns the default bootstrap program.
func Bootstrap() string {
	return bootstrap + "\n\n" + logiclib.Stdlib()
}
