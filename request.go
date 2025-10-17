package tg

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/mrbanja/tg/model"
)

type Request struct {
	ID     string
	Update *model.Update

	data any
}

func newRequest(update *model.Update) *Request {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return &Request{
		ID:     hex.EncodeToString(buf),
		Update: update,
	}
}

func (r *Request) SetData(v any) {
	r.data = v
}

func GetDataFromRequest[T any](r *Request) (T, bool) {
	d, ok := r.data.(T)
	return d, ok
}
