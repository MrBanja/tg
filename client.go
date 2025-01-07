package tg

import (
	"log/slog"
	"os"
)

type Client struct {
	token string

	logger *slog.Logger
}

func New(token string, logger *slog.Logger) *Client {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil,
		}))
	}
	return &Client{
		token:  token,
		logger: logger,
	}
}
