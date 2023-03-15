package fs

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"
)

type URIHandler interface {
	Scheme() string
	Open(ctx context.Context, uri *url.URL) (fs.File, error)
}

type Router struct {
	handlers map[string]URIHandler
}

func NewRouter() Router {
	return Router{
		handlers: make(map[string]URIHandler),
	}
}
func (r *Router) Open(ctx context.Context, name string) (fs.File, error) {
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
