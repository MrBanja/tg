package tg

import (
	"context"
)

type Handler func(ctx context.Context, req *Request)

type handler struct {
	handle  Handler
	filters []Filter
}

func (h *handler) isValid(ctx context.Context, req *Request) bool {
	for _, f := range h.filters {
		if !f(ctx, req) {
			return false
		}
	}
	return true
}
