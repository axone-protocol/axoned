package predicate

import (
	"context"
	"fmt"

	"github.com/ichiban/prolog/engine"
)

type Component string

const (
	QueryComponent    Component = "query"
	FragmentComponent Component = "fragment"
	PathComponent     Component = "path"
	SegmentComponent  Component = "segment"
)

func NewComponent(v string) (Component, error) {
	switch v {
	case string(QueryComponent):
		return QueryComponent, nil
	case string(FragmentComponent):
		return FragmentComponent, nil
	case string(PathComponent):
		return PathComponent, nil
	case string(SegmentComponent):
		return SegmentComponent, nil
	default:
		return "", fmt.Errorf("invalid component name %s, expected `query`, `fragment`, `path` or `segment`", v)
	}
}

func URIEncoded(vm *engine.VM, component, decoded, encoded engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var comp Component
		switch c := env.Resolve(component).(type) {
		case engine.Atom:
			_, err := NewComponent(c.String())
			if err != nil {
				return engine.Error(fmt.Errorf("uri_encoded/3: %w", err))
			}
		default:
			return engine.Error(fmt.Errorf("uri_encoded/3: invalid component type: %T, should be Atom", component))
		}

		fmt.Printf("%s", comp)

		return engine.Bool(true)
	})
}
