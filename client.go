package tg

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/mrbanja/snet/v2"

	"github.com/mrbanja/tg/tgmodel"
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

const baseURL = "https://api.telegram.org/bot%s/%s"

func (c *Client) SendMessage(ctx context.Context, body tgmodel.SendMessage) (*tgmodel.Message, error) {
	url := c.buildURl("sendMessage")
	req, err := snet.NewRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, snet.NewWrongStatusError(resp)
	}

	tgResp, err := snet.UnmarshalResp[tgmodel.Response[tgmodel.Message]](resp)
	if err != nil {
		return nil, err
	}
	return tgResp.Result, err
}

func (c *Client) buildURl(method string) string {
	return fmt.Sprintf(baseURL, c.token, method)
}
