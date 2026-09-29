package tg

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mrbanja/tg/v4/model"
)

// SecretTokenHeader carries the webhook secret token set via WithSecretToken.
const SecretTokenHeader = "X-Telegram-Bot-Api-Secret-Token"

type Server struct {
	addr                string
	webhookPath         string
	discoverableBaseURL string

	opts *serverOptions

	handlers []handler

	// inflight tracks detached handlers (WithDetachedHandlers).
	inflight sync.WaitGroup
}

type serverOptions struct {
	debug bool

	secretToken     string
	allowedUpdates  []string
	replaceWebhook  bool
	detached        bool
	shutdownTimeout time.Duration
}

type ServerOptionFunc func(*serverOptions)

func WithDebug(on bool) ServerOptionFunc {
	return func(o *serverOptions) {
		o.debug = on
	}
}

// WithSecretToken sets the webhook secret_token and rejects webhook requests
// whose X-Telegram-Bot-Api-Secret-Token header does not match.
func WithSecretToken(secret string) ServerOptionFunc {
	return func(o *serverOptions) {
		o.secretToken = secret
	}
}

// WithAllowedUpdates sets the update types passed as allowed_updates in setWebhook.
func WithAllowedUpdates(types ...string) ServerOptionFunc {
	return func(o *serverOptions) {
		o.allowedUpdates = types
	}
}

// WithReplaceWebhook makes SetWebhook replace a webhook that points to a
// different URL instead of returning ErrWebhookAlreadySetToDiffAddress.
// Pending updates are kept.
func WithReplaceWebhook(on bool) ServerOptionFunc {
	return func(o *serverOptions) {
		o.replaceWebhook = on
	}
}

// WithDetachedHandlers acknowledges every update immediately and runs the
// handler in its own goroutine. Handler contexts are derived from the server
// context (cancelled on shutdown) instead of the webhook request context.
// On shutdown the server waits for running handlers up to the shutdown timeout.
func WithDetachedHandlers() ServerOptionFunc {
	return func(o *serverOptions) {
		o.detached = true
	}
}

// WithShutdownTimeout sets how long graceful shutdown waits (default 5s).
func WithShutdownTimeout(d time.Duration) ServerOptionFunc {
	return func(o *serverOptions) {
		o.shutdownTimeout = d
	}
}

func NewServer(addr string, discoverableBaseURL string, optsFn ...ServerOptionFunc) *Server {
	opts := &serverOptions{shutdownTimeout: 5 * time.Second}
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
	// The secret token and allowed updates are not reported by getWebhookInfo,
	// so the webhook is always (re)set when they are configured.
	needsParams := s.opts.secretToken != "" || len(s.opts.allowedUpdates) != 0
	switch {
	case info.URL == whURL && !needsParams:
		slog.Info("[*] webhook is already set")
		return nil
	case info.URL != "" && info.URL != whURL && !s.opts.replaceWebhook:
		slog.Error("[*] webhook is already set to different url", "url", info.URL)
		return ErrWebhookAlreadySetToDiffAddress
	}
	return SetWebhook(ctx, model.SetWebhookRequest{
		URL:            whURL,
		AllowedUpdates: s.opts.allowedUpdates,
		SecretToken:    s.opts.secretToken,
	})
}

func (s *Server) Server(ctx context.Context, handler http.Handler) {
	if s.webhookPath == "" {
		slog.Error("[*] webhook path is empty: call SetWebhook before Server")
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := &http.Server{
		Addr:    s.addr,
		Handler: s.updateHandler(ctx, handler),
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

		ctx, cancel := context.WithTimeout(context.Background(), s.opts.shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("[*] server Shutdown: ", "err", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("[*] ListenAndServe: ", "err", err)
	}

	if s.opts.detached {
		s.waitInflight()
	}
}

func (s *Server) waitInflight() {
	done := make(chan struct{})
	go func() {
		s.inflight.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(s.opts.shutdownTimeout):
		slog.Error("[*] detached handlers did not finish before shutdown timeout")
	}
}

func (s *Server) updateHandler(ctx context.Context, handler http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(s.webhookPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if s.opts.secretToken != "" &&
			subtle.ConstantTimeCompare([]byte(r.Header.Get(SecretTokenHeader)), []byte(s.opts.secretToken)) != 1 {
			slog.Warn("[*] webhook request with invalid secret token", "remote_addr", r.RemoteAddr)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

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

		if s.opts.detached {
			s.inflight.Add(1)
			go func() {
				defer s.inflight.Done()
				defer func() {
					if rec := recover(); rec != nil {
						logger.Error("[*] handler panicked", "panic", rec)
					}
				}()
				s.dispatch(ctx, request, logger)
			}()
			return
		}
		s.dispatch(r.Context(), request, logger)
	})

	if handler != nil {
		mux.Handle("/", handler)
	}
	return mux
}

func (s *Server) dispatch(ctx context.Context, request *Request, logger *slog.Logger) {
	for _, h := range s.handlers {
		if h.isValid(ctx, request) {
			logger.Debug("[*] picked handler")
			h.handle(ctx, request)
			return
		}
	}
	logger.Debug("[*] no handler picked for an update")
}
