package tg

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrbanja/tg/v4/model"
)

type Server struct {
	addr                string
	webhookPath         string
	discoverableBaseURL string

	opts *serverOptions

	handlers []handler
}

type serverOptions struct {
	debug bool
}

type ServerOptionFunc func(*serverOptions)

func WithDebug(on bool) ServerOptionFunc {
	return func(o *serverOptions) {
		o.debug = on
	}
}

func NewServer(addr string, discoverableBaseURL string, optsFn ...ServerOptionFunc) *Server {
	opts := &serverOptions{}
	for _, fn := range optsFn {
		fn(opts)
	}

	return &Server{
		addr:                addr,
		discoverableBaseURL: discoverableBaseURL,
		opts:                opts,
	}
}

// Register is not thread-safe
func (s *Server) Register(h Handler, f ...Filter) {
	s.handlers = append(s.handlers, handler{handle: h, filters: f})
}

func (s *Server) SetWebhook(ctx context.Context, path string) error {
	s.webhookPath = path
	whURL, err := url.JoinPath(s.discoverableBaseURL, path)
	if err != nil {
		slog.Error("[*] url.JoinPath failed: ", "err", err)
		return err
	}

	info, err := GetWebhookInfo(ctx)
	if err != nil {
		return err
	}
	if info.URL == whURL {
		slog.Info("[*] webhook is already set")
		return nil
	} else if info.URL != "" {
		slog.Error("[*] webhook is already set to different url", "url", info.URL)
		return ErrWebhookAlreadySetToDiffAddress
	}
	return SetWebhook(ctx, model.SetWebhookRequest{URL: whURL})
}

func (s *Server) Server(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := &http.Server{
		Addr:    s.addr,
		Handler: s.updateHandler(),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		var sigCh = make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case <-ctx.Done():
		case <-sigCh:
			cancel()
		}

		slog.Info("[*] server is gracefully shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("[*] server Shutdown: ", "err", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("[*] ListenAndServe: ", "err", err)
	}
}

func (s *Server) updateHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(s.webhookPath, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("[*] read body failed: ", "err", err)
			http.Error(w, "read body failed", http.StatusBadRequest)
			return
		}

		if s.opts.debug { //> Hack for better printing text
			var d any
			json.Unmarshal(body, &d)
			slog.Debug("[*] update received", slog.Any("body", d))
		}

		var update model.Update
		if err = json.Unmarshal(body, &update); err != nil {
			slog.Error("[*] unmarshal body failed: ", "err", err)
			http.Error(w, "unmarshal body failed", http.StatusBadRequest)
			return
		}
		logger := slog.With(slog.Int64("update_id", update.ID))

		var request = newRequest(&update)
		logger = logger.With(slog.String("request_id", request.ID))

		for _, h := range s.handlers {
			if h.isValid(r.Context(), request) {
				logger.Debug("[*] picked handler")
				h.handle(r.Context(), request)
				return
			}
		}
		logger.Debug("[*] no handler picked for an update")
	})
	return mux
}
