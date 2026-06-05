package codec

import (
	"errors"

	"github.com/axone-protocol/prolog/v3/engine"

	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
	"github.com/axone-protocol/axoned/v15/x/logic/util"
)

var (
	atomText          = engine.NewAtom("text")
	atomBytes         = engine.NewAtom("bytes")
	atomTypeCharset   = engine.NewAtom("charset")
	atomValidEncoding = engine.NewAtom("encoding")
)

type textCodec struct{}

func init() {
	Register(&textCodec{})
}

func (c *textCodec) Name() string {
	return "text"
}

func (c *textCodec) Encode(payload []byte) engine.Term {
	term, err := parseCodecTerm(payload)
	if err != nil {
		return prolog.AtomError.Apply(atomSyntaxError.Apply(atomPrologSyntax.Apply(atomMalformedTerm)))
	}

	request, ok := term.(engine.Compound)
	if !ok || request.Functor() != atomText || request.Arity() != 2 {
		return errInvalidRequest
	}

	env := engine.NewEnv()
	encoding, err := prolog.AssertAtom(request.Arg(0), env)
	if err != nil {
		return prolog.AtomError.Apply(exceptionFormal(err))
	}

	text, err := prolog.TextTermToString(request.Arg(1), env)
	if err != nil {
		return prolog.AtomError.Apply(exceptionFormal(err))
	}

	var bs []byte
	switch encoding {
	case prolog.AtomText:
		bs = []byte(text)
	case prolog.AtomOctet:
		term, err := prolog.StringToOctetListTerm(text, env)
		if err != nil {
			return prolog.AtomError.Apply(exceptionFormal(err))
		}
		return atomOK.Apply(term)
	default:
		bs, err = encodeText(request.Arg(1), text, encoding, env)
		if err != nil {
			return prolog.AtomError.Apply(exceptionFormal(err))
		}
	}

	return atomOK.Apply(prolog.BytesToByteListTerm(bs))
}

func (c *textCodec) Decode(payload []byte) engine.Term {
	term, err := parseCodecTerm(payload)
	if err != nil {
		return prolog.AtomError.Apply(atomSyntaxError.Apply(atomPrologSyntax.Apply(atomMalformedTerm)))
	}

	request, ok := term.(engine.Compound)
	if !ok || request.Functor() != atomBytes || request.Arity() != 2 {
		return errInvalidRequest
	}

	env := engine.NewEnv()
	encoding, err := prolog.AssertAtom(request.Arg(0), env)
	if err != nil {
		return prolog.AtomError.Apply(exceptionFormal(err))
	}

	var text string
	switch encoding {
	case prolog.AtomText:
		bs, err := prolog.ByteListTermToBytes(request.Arg(1), env)
		if err != nil {
			return prolog.AtomError.Apply(exceptionFormal(err))
		}
		text = string(bs)
	case prolog.AtomOctet:
		text, err = prolog.OctetListTermToString(request.Arg(1), env)
		if err != nil {
			return prolog.AtomError.Apply(exceptionFormal(err))
		}
	default:
		bs, err := prolog.ByteListTermToBytes(request.Arg(1), env)
		if err != nil {
			return prolog.AtomError.Apply(exceptionFormal(err))
		}
		text, err = decodeText(request.Arg(1), bs, encoding, env)
		if err != nil {
			return prolog.AtomError.Apply(exceptionFormal(err))
		}
	}

	return atomOK.Apply(prolog.StringToCharacterListTerm(text))
}

func encodeText(value engine.Term, text string, encoding engine.Atom, env *engine.Env) ([]byte, error) {
	bs, err := util.Encode(text, encoding.String())
	if err != nil {
		switch {
		case errors.Is(err, util.ErrInvalidCharset):
			return nil, engine.TypeError(atomTypeCharset, encoding, env)
		default:
			return nil, prolog.WithError(engine.DomainError(validEncoding(encoding.String()), value, env), err, env)
		}
	}
	return bs, nil
}

func decodeText(value engine.Term, bs []byte, encoding engine.Atom, env *engine.Env) (string, error) {
	text, err := util.Decode(bs, encoding.String())
	if err != nil {
		switch {
		case errors.Is(err, util.ErrInvalidCharset):
			return "", engine.TypeError(atomTypeCharset, encoding, env)
		default:
			return "", prolog.WithError(engine.DomainError(validEncoding(encoding.String()), value, env), err, env)
		}
	}
	return text, nil
}

func validEncoding(encoding string) engine.Term {
	return atomValidEncoding.Apply(engine.NewAtom(encoding))
}
