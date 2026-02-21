package prologterm

import (
	"context"
	"strings"

	"github.com/axone-protocol/prolog/v3/engine"

	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
)

// Render serializes a Prolog term and appends ".\n" for file-based framing.
func Render(term engine.Term, quoted bool) ([]byte, error) {
	var sb strings.Builder
	stream := engine.NewOutputTextStream(&sb)
	var vm engine.VM

	_, err := engine.WriteTerm(
		&vm,
		stream,
		term,
		writeOptions(quoted),
		engine.Success,
		nil,
	).Force(context.Background())
	if err != nil {
		return nil, err
	}

	sb.WriteString(".\n")

	return []byte(sb.String()), nil
}

func writeOptions(quoted bool) engine.Term {
	quotedValue := prolog.AtomFalse
	if quoted {
		quotedValue = prolog.AtomTrue
	}

	return engine.List(prolog.AtomQuoted.Apply(quotedValue))
}
