package tg

import (
	"context"

	"github.com/mrbanja/tg/v2/model"
)

type Handler func(ctx context.Context, u *model.Update)

type handler struct {
	handle  Handler
	filters []Filter
}

func (h *handler) isValid(ctx context.Context, u *model.Update) bool {
	for _, f := range h.filters {
		if !f(ctx, u) {
			return false
		}
	}
	return true
}
