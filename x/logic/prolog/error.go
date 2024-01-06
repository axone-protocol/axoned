package prolog

import "github.com/ichiban/prolog/engine"

// Error atoms
var (
	AtomEncodingError = engine.NewAtom("encoding_error") // AtomEncodingError is the term used to indicate the encoding error.
)

// EncodingError returns the compound term error(encoding_error(Encoding, cause)).
func EncodingError(encoding string, cause error, env *engine.Env) engine.Exception {
	return engine.NewException(AtomError.Apply(AtomEncodingError.Apply(engine.NewAtom(encoding)), StringToTerm(cause.Error())), env)
}
