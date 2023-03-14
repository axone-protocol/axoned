package fs

import (
	"context"
	"fmt"
	"net/url"
)

type URIHandler interface {
	Scheme() string
	Open(ctx context.Context, uri *url.URL) ([]byte, error)
}

type Router struct {
	handlers map[string]URIHandler
}

func NewRouter() Router {
	return Router{
		handlers: make(map[string]URIHandler),
	}
}
func (r *Router) Open(ctx context.Context, name string) ([]byte, error) {
	uri, err := url.Parse(name)
	if err != nil {
		return nil, err
	}

	handler, ok := r.handlers[uri.Scheme]

	if !ok {
		return nil, fmt.Errorf("could not find handler for load %s file", name)
	}

	return handler.Open(ctx, uri)
}

func (r *Router) RegisterHandler(handler URIHandler) {
	r.handlers[handler.Scheme()] = handler
}
