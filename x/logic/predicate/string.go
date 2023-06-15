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

// ReadString is a predicate that will read a given Stream and unify it to String.
// Optionally give a max length of reading, when stream reach the given length, the reading is stop.
// If Length is unbound, Stream is read to the end and Length is unified with the number of characters read.
//
// read_string(+Stream, ?Length, -String) is det
//
// Where
//   - `Stream`: represent a stream
//   - `Length`: is the max length to read
//   - `String`: represent the unified read stream as string
//
// Example:
//
// # Given a file `foo.txt` that contains `Hello World`:
// ```
// file_to_string(File, String, Length) :-
//
//	open(File, read, In),
//	read_string(In, Length, String),
//	close(Stream).
//
// ```
//
// Result :
//
// ```
//
//	?- file_to_string('path/file/foo.txt', String, Length).
//
//	String = 'Hello World'
//	Length = 11
//
// ```.
func ReadString(vm *engine.VM, stream, length, result engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var s *engine.Stream
		switch st := env.Resolve(stream).(type) {
		case engine.Variable:
			return engine.Error(fmt.Errorf("read_string/3: stream could not be a variable"))
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
				return engine.Error(fmt.Errorf("read_string/3: error occurs reading stream: %w", err))
			}
			totalLen += uint64(l)
			_, err = builder.WriteRune(r)
			if err != nil {
				return engine.Error(fmt.Errorf("read_string/3: failed write string: %w", err))
			}
		}

		return engine.Unify(vm, Tuple(result, length), Tuple(util.StringToTerm(builder.String()), engine.Integer(totalLen)), cont, env)
	})
}
