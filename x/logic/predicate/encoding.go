package predicate

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/axone-protocol/prolog/v2/engine"

	"github.com/axone-protocol/axoned/v11/x/logic/prolog"
)

var (
	atomClassic = engine.NewAtom("classic")
	atomURL     = engine.NewAtom("url")
	atomString  = engine.NewAtom("string")
	atomAtom    = engine.NewAtom("atom")
)

// Base64Encoded is a predicate that unifies a string to a base64 encoded string as specified by [RFC 4648].
//
// The signature is as follows:
//
//	base64_encoded(+Plain, -Encoded, +Options) is det
//	base64_encoded(-Plain, +Encoded, +Options) is det
//	base64_encoded(+Plain, +Encoded, +Options) is det
//
// Where:
//   - Plain is an atom, a list of character codes, or list of characters containing the unencoded (plain) text.
//   - Encoded is an atom or string containing the base64 encoded text.
//   - Options is a list of options that can be used to control the encoding process.
//
// # Options
//
// The following options are supported:
//
//   - padding(+Boolean)
//
// If true (default), the output is padded with = characters.
//
//   - charset(+Charset)
//
// Define the encoding character set to use. The (default) 'classic' uses the classical rfc2045 characters. The value 'url'
// uses URL and file name friendly characters.
//
//   - as(+Type)
//
// Defines the type of the output. One of string (default) or atom.
//
//   - encoding(+Encoding)
//
// Encoding to use for translation between (Unicode) text and bytes (Base64 is an encoding for bytes). Default is utf8.
//
// [RFC 4648]: https://rfc-editor.org/rfc/rfc4648.html
func Base64Encoded(_ *engine.VM, plain, encoded, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	encoding, err := getBase64Encoding(options, env)
	if err != nil {
		return engine.Error(err)
	}

	asOpt, err := prolog.GetOptionWithDefault(prolog.AtomAs, options, atomString, env)
	if err != nil {
		return engine.Error(err)
	}

	forwardConverter := func(value []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		src, err := prolog.TextTermToString(value[0], env)
		if err != nil {
			return nil, err
		}
		dst := encoding.EncodeToString([]byte(src))

		switch asOpt {
		case atomString:
			return []engine.Term{prolog.StringToCharacterListTerm(dst)}, nil
		case atomAtom:
			return []engine.Term{engine.NewAtom(dst)}, nil
		default:
			return nil, engine.DomainError(prolog.AtomAs, asOpt, env)
		}
	}
	backwardConverter := func(value []engine.Term, options engine.Term, env *engine.Env) ([]engine.Term, error) {
		src, err := prolog.TextTermToString(value[0], env)
		if err != nil {
			return nil, err
		}

		dst, err := encoding.DecodeString(src)
		if err != nil {
			return nil,
				prolog.WithError(
					engine.DomainError(prolog.ValidEncoding("base64"), value[0], env), err, env)
		}

		encodingOpt, err := prolog.GetOptionAsAtomWithDefault(prolog.AtomEncoding, options, prolog.AtomUtf8, env)
		if err != nil {
			return nil, err
		}
		dstStr, err := prolog.Decode(value[0], dst, encodingOpt, env)
		if err != nil {
			return nil, err
		}

		switch asOpt {
		case atomString:
			return []engine.Term{prolog.StringToCharacterListTerm(dstStr)}, nil
		case atomAtom:
			return []engine.Term{engine.NewAtom(dstStr)}, nil
		default:
			return nil, engine.DomainError(prolog.AtomAs, asOpt, env)
		}
	}

	return prolog.UnifyFunctionalPredicate(
		[]engine.Term{plain}, []engine.Term{encoded}, options, forwardConverter, backwardConverter, cont, env)
}

// getBase64Encoding returns the base64 encoding based on the options provided.
func getBase64Encoding(options engine.Term, env *engine.Env) (*base64.Encoding, error) {
	var encoding *base64.Encoding

	charsetOpt, err := prolog.GetOptionAsAtomWithDefault(prolog.AtomCharset, options, atomClassic, env)
	if err != nil {
		return nil, err
	}
	switch charsetOpt {
	case atomClassic:
		encoding = base64.StdEncoding
	case atomURL:
		encoding = base64.URLEncoding
	default:
		return nil, engine.DomainError(prolog.AtomCharset, charsetOpt, env)
	}

	paddingOpt, err := prolog.GetOptionWithDefault(prolog.AtomPadding, options, prolog.AtomTrue, env)
	if err != nil {
		return nil, err
	}
	switch paddingOpt {
	case prolog.AtomTrue:
		encoding = encoding.WithPadding(base64.StdPadding)
	case prolog.AtomFalse:
		encoding = encoding.WithPadding(base64.NoPadding)
	default:
		return nil, engine.DomainError(prolog.AtomPadding, paddingOpt, env)
	}
	return encoding, nil
}

// HexBytes is a predicate that unifies hexadecimal encoded bytes to a list of bytes.
//
// The signature is as follows:
//
//	hex_bytes(?Hex, ?Bytes) is det
//
// Where:
//   - Hex is an Atom, string or list of characters in hexadecimal encoding.
//   - Bytes is the list of numbers between 0 and 255 that represent the sequence of bytes.
//
// # Examples:
//
//	# Convert hexadecimal atom to list of bytes.
//	- hex_bytes('2c26b46b68ffc68ff99b453c1d3041341342d706483bfa0f98a5e886266e7ae', Bytes).
func HexBytes(vm *engine.VM, hexa, bts engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	var result []byte

	switch h := env.Resolve(hexa).(type) {
	case engine.Variable:
	case engine.Atom:
		src := []byte(h.String())
		result = make([]byte, hex.DecodedLen(len(src)))
		_, err := hex.Decode(result, src)
		if err != nil {
			return engine.Error(
				prolog.WithError(
					engine.DomainError(prolog.ValidEncoding("hex"), hexa, env), err, env))
		}
	default:
		return engine.Error(engine.TypeError(prolog.AtomTypeAtom, hexa, env))
	}

	switch b := env.Resolve(bts).(type) {
	case engine.Variable:
		if result == nil {
			return engine.Error(engine.InstantiationError(env))
		}
		return engine.Unify(vm, bts, prolog.BytesToByteListTerm(result), cont, env)
	case engine.Compound:
		src, err := prolog.ByteListTermToBytes(b, env)
		if err != nil {
			return engine.Error(err)
		}
		dst := hex.EncodeToString(src)
		var r engine.Term = engine.NewAtom(dst)
		return engine.Unify(vm, hexa, r, cont, env)
	default:
		return engine.Error(engine.TypeError(prolog.AtomTypeText, bts, env))
	}
}
