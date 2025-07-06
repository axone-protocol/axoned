package prolog

import (
	"errors"

	"github.com/axone-protocol/prolog/v2/engine"

	"github.com/axone-protocol/axoned/v12/x/logic/util"
)

// Encode encodes the given string with the given encoding.
// value is used for error reporting (the culprit).
func Encode(value engine.Term, str string, encoding engine.Atom, env *engine.Env) ([]byte, error) {
	bs, err := util.Encode(str, encoding.String())
	if err != nil {
		switch {
		case errors.Is(err, util.ErrInvalidCharset):
			return nil, engine.TypeError(AtomTypeCharset, encoding, env)
		default:
			return nil, WithError(
				engine.DomainError(ValidEncoding(encoding.String()), value, env), err, env)
		}
	}
	return bs, nil
}

// Decode decodes the given byte slice with the given encoding.
// value is used for error reporting (the culprit).
func Decode(value engine.Term, bs []byte, encoding engine.Atom, env *engine.Env) (string, error) {
	str, err := util.Decode(bs, encoding.String())
	if err != nil {
		switch {
		case errors.Is(err, util.ErrInvalidCharset):
			return "", engine.TypeError(AtomTypeCharset, encoding, env)
		default:
			return "", WithError(
				engine.DomainError(ValidEncoding(encoding.String()), value, env), err, env)
		}
	}
	return str, nil
}
