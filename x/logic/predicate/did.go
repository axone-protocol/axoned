package predicate

import (
	"strings"

	"github.com/axone-protocol/prolog/v2/engine"
	godid "github.com/nuts-foundation/go-did/did"
	"github.com/samber/lo"

	"github.com/axone-protocol/axoned/v11/x/logic/prolog"
)

// DIDPrefix is the prefix for a DID.
const DIDPrefix = "did"

// DIDComponents is a predicate which breaks down a DID into its components according to the [W3C DID] specification.
//
// The signature is as follows:
//
//	did_components(+DID, -Components) is det
//	did_components(-DID, +Components) is det
//
// where:
//   - DID represent DID URI, given as an Atom, compliant with [W3C DID] specification.
//   - Components is a compound Term in the format did(Method, ID, Path, Query, Fragment), aligned with the [DID syntax],
//     where: Method is the method name, ID is the method-specific identifier, Path is the path component, Query is the
//     query component and Fragment is the fragment component. Values are given as an Atom and are url encoded.
//     For any component not present, its value will be null and thus will be left as an uninstantiated variable.
//
// # Examples:
//
//	# Decompose a DID into its components.
//	- did_components('did:example:123456?versionId=1', did_components(Method, ID, Path, Query, Fragment)).
//
//	# Reconstruct a DID from its components.
//	- did_components(DID, did_components('example', '123456', _, 'versionId=1', _42)).
//
// [W3C DID]: https://w3c.github.io/did-core
// [DID syntax]: https://w3c.github.io/did-core/#did-syntax
//
//nolint:funlen
func DIDComponents(vm *engine.VM, did, components engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	switch t1 := env.Resolve(did).(type) {
	case engine.Variable:
	case engine.Atom:
		parsedDid, err := godid.ParseDIDURL(t1.String())
		if err != nil {
			return engine.Error(prolog.WithError(engine.DomainError(prolog.ValidEncoding("did"), did, env), err, env))
		}

		terms := lo.Map([]string{parsedDid.Method, parsedDid.ID, parsedDid.Path, parsedDid.Query.Encode(), parsedDid.Fragment},
			func(segment string, _ int) engine.Term {
				if segment == "" {
					return engine.NewVariable()
				}
				return engine.NewAtom(segment)
			})

		return engine.Unify(vm, components, prolog.AtomDIDComponents.Apply(terms...), cont, env)
	default:
		return engine.Error(engine.TypeError(prolog.AtomTypeAtom, did, env))
	}

	switch t2 := env.Resolve(components).(type) {
	case engine.Variable:
		return engine.Error(engine.InstantiationError(env))
	case engine.Compound:
		if t2.Functor() != prolog.AtomDIDComponents || t2.Arity() != 5 {
			return engine.Error(engine.DomainError(prolog.AtomDIDComponents, components, env))
		}

		buf := strings.Builder{}
		buf.WriteString(DIDPrefix)

		for i := 0; i < t2.Arity(); i++ {
			sep := ""
			switch i {
			case 0, 1:
				sep = ":"
			case 2:
				sep = "/"
			case 3:
				sep = "?"
			case 4:
				sep = "#"
			}
			switch segment := env.Resolve(t2.Arg(i)).(type) {
			case engine.Variable:
			default:
				atom, err := prolog.AssertAtom(segment, env)
				if err != nil {
					return engine.Error(err)
				}
				if !strings.HasPrefix(atom.String(), sep) {
					buf.WriteString(sep)
				}
				buf.WriteString(atom.String())
			}
		}

		return engine.Unify(vm, did, engine.NewAtom(buf.String()), cont, env)
	default:
		return engine.Error(engine.TypeError(prolog.AtomDIDComponents, components, env))
	}
}
