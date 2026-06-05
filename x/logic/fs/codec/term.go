package codec

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	goprolog "github.com/axone-protocol/prolog/v3"
	"github.com/axone-protocol/prolog/v3/engine"
)

func parseCodecTerm(payload []byte) (engine.Term, error) {
	interpreter := goprolog.New(strings.NewReader(""), io.Discard)
	parser := engine.NewParser(&interpreter.VM, bytes.NewReader(payload))
	term, err := parser.Term()
	if err != nil {
		return nil, err
	}
	if parser.More() {
		return nil, fmt.Errorf("unexpected trailing term")
	}
	return term, nil
}
