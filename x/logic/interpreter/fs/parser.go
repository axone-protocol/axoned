package fs

import (
	"context"
	"fmt"
	"net/url"
)

type URIHandler interface {
	CanOpen(ctx context.Context, uri *url.URL) bool
	Open(ctx context.Context, uri *url.URL) ([]byte, error)
}

type Parser struct {
	Handlers []URIHandler
}

func (p *Parser) Parse(ctx context.Context, name string) ([]byte, error) {
	uri, err := url.Parse(name)
	if err != nil {
		return nil, err
	}

	if uri.Scheme != "okp4" {
		return nil, fmt.Errorf("incompatible schema '%s' for %s", uri.Scheme, name)
	}

	for _, handler := range p.Handlers {
		if handler.CanOpen(ctx, uri) {
			return handler.Open(ctx, uri)
		}
	}

	return nil, fmt.Errorf("could not find handler for load %s file", name)
}
