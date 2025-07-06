package predicate

import (
	"errors"
	"io"
	"strings"

	"github.com/axone-protocol/prolog/v2/engine"

	"github.com/axone-protocol/axoned/v12/x/logic/prolog"
)

// ReadString is a predicate that reads characters from the provided Stream and unifies them with String.
// Users can optionally specify a maximum length for reading; if the stream reaches this length, the reading stops.
// If Length remains unbound, the entire Stream is read, and upon completion, Length is unified with the count of characters read.
//
// The signature is as follows:
//
//	read_string(+Stream, ?Length, -String) is det
//
// Where:
//   - Stream is the input stream to read from.
//   - Length is the optional maximum number of characters to read from the Stream. If unbound, denotes the full length of Stream.
//   - String is the resultant string after reading from the Stream.
//
// # Examples:
//
//	# Given a file `foo.txt` that contains `Hello World`:
//
//	 file_to_string(File, String, Length) :-
//	 open(File, read, In),
//	 read_string(In, Length, String),
//	 close(Stream).
//
//	# It gives:
//	?- file_to_string('path/file/foo.txt', String, Length).
//
//	String = 'Hello World'
//	Length = 11
func ReadString(vm *engine.VM, stream, length, result engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	is, err := prolog.AssertStream(stream, env)
	if err != nil {
		return engine.Error(err)
	}

	var maxLength uint64
	if maxLen, ok := env.Resolve(length).(engine.Integer); ok {
		maxLength = uint64(maxLen) //nolint:gosec // disable G115
	}

	var builder strings.Builder
	var totalLen uint64
	for {
		r, l, err := is.ReadRune()
		if err != nil || (maxLength != 0 && totalLen >= maxLength) {
			if errors.Is(err, io.EOF) || totalLen >= maxLength {
				break
			}
			return engine.Error(engine.SyntaxError(prolog.ErrorTerm(err), env))
		}
		totalLen += uint64(l) //nolint:gosec // disable G115
		_, err = builder.WriteRune(r)
		if err != nil {
			return engine.Error(engine.SyntaxError(prolog.ErrorTerm(err), env))
		}
	}

	var r engine.Term = engine.NewAtom(builder.String())
	return engine.Unify(
		vm, prolog.Tuple(result, length),
		prolog.Tuple(r, engine.Integer(totalLen)), cont, env) //nolint:gosec // disable G115
}

// StringBytes is a predicate that unifies a string with a list of bytes, returning true when the (Unicode) String is
// represented by Bytes in Encoding.
//
// The signature is as follows:
//
//	string_bytes(?String, ?Bytes, +Encoding)
//
// Where:
//   - String is the string to convert to bytes. It can be an Atom, string or list of characters codes.
//   - Bytes is the list of numbers between 0 and 255 that represent the sequence of bytes.
//   - Encoding is the encoding to use for the conversion.
//
// Encoding can be one of the following:
// - 'text' considers the string as a sequence of Unicode characters.
// - 'octet' considers the string as a sequence of bytes.
// - '<encoding>' considers the string as a sequence of characters in the given encoding.
//
// At least one of String or Bytes must be instantiated.
//
// # Examples:
//
//	# Convert a string to a list of bytes.
//	- string_bytes('Hello World', Bytes, octet).
//
//	# Convert a list of bytes to a string.
//	- string_bytes(String, [72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100], octet).
func StringBytes(
	_ *engine.VM, str, bts, encodingTerm engine.Term, cont engine.Cont, env *engine.Env,
) *engine.Promise {
	encoding, err := prolog.AssertAtom(encodingTerm, env)
	if err != nil {
		return engine.Error(err)
	}
	forwardConverter := func(value []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		str, err := prolog.TextTermToString(value[0], env)
		if err != nil {
			return nil, err
		}

		switch encoding {
		case prolog.AtomText:
			return []engine.Term{prolog.StringToByteListTerm(str)}, nil
		case prolog.AtomOctet:
			term, err := prolog.StringToOctetListTerm(str, env)
			if err != nil {
				return nil, err
			}
			return []engine.Term{term}, nil
		default:
			bs, err := prolog.Encode(value[0], str, encoding, env)
			if err != nil {
				return nil, err
			}

			return []engine.Term{prolog.BytesToByteListTerm(bs)}, nil
		}
	}
	backwardConverter := func(value []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		var result string
		switch encoding {
		case prolog.AtomText:
			bs, err := prolog.ByteListTermToBytes(value[0], env)
			if err != nil {
				return nil, err
			}
			result = string(bs)
		case prolog.AtomOctet:
			result, err = prolog.OctetListTermToString(value[0], env)
			if err != nil {
				return nil, err
			}
		default:
			bs, err := prolog.ByteListTermToBytes(value[0], env)
			if err != nil {
				return nil, err
			}
			result, err = prolog.Decode(value[0], bs, encoding, env)
			if err != nil {
				return nil, err
			}
		}
		return []engine.Term{prolog.StringToCharacterListTerm(result)}, nil
	}

	return prolog.UnifyFunctionalPredicate(
		[]engine.Term{str}, []engine.Term{bts}, encoding, forwardConverter, backwardConverter, cont, env)
}
