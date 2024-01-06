package predicate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/prolog"
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
// Examples:
//
//	# Given a file `foo.txt` that contains `Hello World`:
//
//	file_to_string(File, String, Length) :-
//
//	open(File, read, In),
//	read_string(In, Length, String),
//	close(Stream).
//
//	# It gives:
//	?- file_to_string('path/file/foo.txt', String, Length).
//
//	String = 'Hello World'
//	Length = 11
func ReadString(vm *engine.VM, stream, length, result engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var s *engine.Stream
		switch st := env.Resolve(stream).(type) {
		case engine.Variable:
			return engine.Error(fmt.Errorf("read_string/3: stream cannot be a variable"))
		case *engine.Stream:
			s = st
		default:
			return engine.Error(fmt.Errorf("read_string/3: invalid domain for given stream"))
		}

		var maxLength uint64
		if maxLen, ok := env.Resolve(length).(engine.Integer); ok {
			maxLength = uint64(maxLen)
		}

		var builder strings.Builder
		var totalLen uint64
		for {
			r, l, err := s.ReadRune()
			if err != nil || (maxLength != 0 && totalLen >= maxLength) {
				if errors.Is(err, io.EOF) || totalLen >= maxLength {
					break
				}
				return engine.Error(fmt.Errorf("read_string/3: couldn't read stream: %w", err))
			}
			totalLen += uint64(l)
			_, err = builder.WriteRune(r)
			if err != nil {
				return engine.Error(fmt.Errorf("read_string/3: couldn't write string: %w", err))
			}
		}

		return engine.Unify(
			vm, prolog.Tuple(result, length),
			prolog.Tuple(prolog.StringToTerm(builder.String()), engine.Integer(totalLen)), cont, env)
	})
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
// - 'utf8' considers the string as a sequence of UTF-8 characters.
// - '<encoding>' considers the string as a sequence of characters in the given encoding.
//
// At least one of String or Bytes must be instantiated.
//
// Examples:
//
//	# Convert a string to a list of bytes.
//	- string_bytes('Hello World', Bytes, octet).
//
//	# Convert a list of bytes to a string.
//	- string_bytes(String, [72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100], octet).
func StringBytes(
	_ *engine.VM, str, bts, encoding engine.Term, cont engine.Cont, env *engine.Env,
) *engine.Promise {
	encodingAtom, err := prolog.AssertAtom(env, encoding)
	if err != nil {
		return engine.Error(err)
	}
	forwardConverter := func(value []engine.Term, options engine.Term, env *engine.Env) ([]engine.Term, error) {
		bs, err := prolog.StringTermToBytes(value[0], encodingAtom.String(), env)
		if err != nil {
			return nil, err
		}
		result, err := prolog.BytesToCodepointListTerm(bs, "text")
		if err != nil {
			return nil, err
		}
		return []engine.Term{result}, nil
	}
	backwardConverter := func(value []engine.Term, options engine.Term, env *engine.Env) ([]engine.Term, error) {
		if _, err := prolog.AssertList(env, value[0]); err != nil {
			return nil, err
		}
		bs, err := prolog.StringTermToBytes(value[0], "text", env)
		if err != nil {
			return nil, err
		}
		result, err := prolog.BytesToAtomListTerm(bs, encodingAtom.String())
		if err != nil {
			return nil, err
		}
		return []engine.Term{result}, nil
	}

	ok, env, err := prolog.UnifyFunctional([]engine.Term{str}, []engine.Term{bts}, encoding, forwardConverter, backwardConverter, env)
	if err != nil {
		return engine.Error(err)
	}
	if !ok {
		return engine.Bool(false)
	}
	return cont(env)
}
