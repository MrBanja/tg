package server

import (
	"context"
	"net/http"

	"github.com/mrbanja/snet/v2"

	"github.com/mrbanja/tg/tgmodel"
)

const SecretTokenHeader = "X-Telegram-Bot-Api-Secret-Token"

func (s *Server) webhookHandler(w http.ResponseWriter, r *http.Request) {
	if s.opt.whSecret != "" && r.Header.Get(SecretTokenHeader) != s.opt.whSecret {
		s.logger.Warn("invalid secret token")
		return
	}

	update, err := snet.UnmarshalReq[tgmodel.Update](r)
	if err != nil {
		s.logger.Error("unmarshal update", "error", err)
		return
	}

	s.logger.Debug("received update", "update", update)
	if err := s.processUpdate(r.Context(), update); err != nil {
		s.logger.Error("process update", "error", err)
	}
}

func (s *Server) processUpdate(ctx context.Context, u *tgmodel.Update) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	for _, r := range s.h {
		var isMatch = true
		for _, f := range r.Filters {
			if !f(u) {
				isMatch = false
				break
			}
		}

		if isMatch {
			if err := r.HandlerFunc(ctx, u); err != nil {
				return err
			}
			break
		}
	}

	return nil
}
