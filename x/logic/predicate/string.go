package predicate

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
)

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
		var builder strings.Builder
		var len uint64 = 0
		for {
			r, l, err := s.ReadRune()
			len += uint64(l)
			if err != nil {
				if err == io.EOF {
					break
				}
				return engine.Error(fmt.Errorf("read_string/3: error occurs reading stream: %w", err))
			}
			_, err = builder.WriteRune(r)
			if err != nil {
				return engine.Error(fmt.Errorf("read_string/3: failed write string: %w", err))
			}
		}

		return engine.Unify(vm, Tuple(result, length), Tuple(util.StringToTerm(builder.String()), engine.Integer(len)), cont, env)
	})
}
