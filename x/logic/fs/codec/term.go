package codec

import (
	"fmt"
	"io"
	"strings"

	goprolog "github.com/axone-protocol/prolog/v3"
	"github.com/axone-protocol/prolog/v3/engine"
)

func parseCodecTerm(payload []byte) (engine.Term, error) {
	interpreter := goprolog.New(strings.NewReader(""), io.Discard)
	parser := engine.NewParser(&interpreter.VM, strings.NewReader(string(payload)))
	term, err := parser.Term()
	if err != nil {
		return nil, err
	}
	if parser.More() {
		return nil, fmt.Errorf("unexpected trailing term")
	}
	return term, nil
}
