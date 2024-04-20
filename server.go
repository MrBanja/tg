package tg

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/mrbanja/snet/v2"

	"github.com/mrbanja/tg/server"
)

func NewServer(opt *server.Opt, logger *slog.Logger) *server.Server {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil,
		}))
	}
	return server.New(opt, logger)
}

func ServerOptions() *server.Opt {
	return server.Options()
}

func ListenWebhook(
	ctx context.Context,
	addr string,
	path string,
	whFunc http.HandlerFunc,
	logger *slog.Logger,
) error {
	if whFunc == nil {
		return errors.New("whFunc is required")
	}

	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil,
		}))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST "+path, whFunc)

	s := &http.Server{
		Addr:        addr,
		Handler:     mux,
		BaseContext: func(listener net.Listener) context.Context { return ctx },
	}
	return snet.ListenAndServe(ctx, s, logger)
}
