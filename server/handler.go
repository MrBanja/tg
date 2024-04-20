package server

import (
	"context"

	"github.com/mrbanja/tg/tgmodel"
)

type Handler struct {
	HandlerFunc HandlerFunc
	Filters     []FilterFunc
}

type HandlerFunc func(ctx context.Context, u *tgmodel.Update) error
type FilterFunc func(u *tgmodel.Update) bool
