package predicate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/x/logic/util"
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
			vm, util.Tuple(result, length),
			util.Tuple(util.StringToTerm(builder.String()), engine.Integer(totalLen)), cont, env)
	})
}
