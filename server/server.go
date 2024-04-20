package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/mrbanja/snet/v2"
)

type Server struct {
	opt    *Opt
	logger *slog.Logger

	h []Handler
}

func New(
	opt *Opt,
	logger *slog.Logger,
) *Server {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil,
		}))
	}
	return &Server{
		opt:    opt,
		logger: logger,

		h: make([]Handler, 0),
	}
}

func (s *Server) HttpServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+s.opt.whPath, s.webhookHandler)

	return &http.Server{
		Addr:    s.opt.addr,
		Handler: mux,
	}
}

func (s *Server) ListenAndServe(ctx context.Context, signals ...os.Signal) error {
	if err := s.SetWebhook(ctx); err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}
	server := s.HttpServer()
	server.BaseContext = func(listener net.Listener) context.Context { return ctx }
	return snet.ListenAndServe(ctx, server, s.logger, signals...)
}

// Register is not concurrent safe. Use it before calling ListenAndServe.
func (s *Server) Register(h HandlerFunc, f ...FilterFunc) {
	s.h = append(s.h, Handler{
		HandlerFunc: h,
		Filters:     f,
	})
}

func (s *Server) SetWebhook(ctx context.Context) error {
	if s.opt.tgToken == "" || s.opt.selfURL == "" {
		return nil
	}

	u, _ := url.Parse("https://api.telegram.org/bot" + s.opt.tgToken)
	u = u.JoinPath("/setWebhook")
	q := u.Query()
	q.Set("url", s.opt.selfURL+s.opt.whPath)
	if s.opt.whSecret != "" {
		q.Set("secret_token", s.opt.whSecret)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return snet.NewWrongStatusError(resp)
	}
	return nil
}
