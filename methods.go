package tg

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mrbanja/snet/v2"

	"github.com/mrbanja/tg/tgmodel"
)

const baseURL = "https://api.telegram.org/bot%s/%s"

func (c *Client) SendMessage(ctx context.Context, body tgmodel.SendMessage) (*tgmodel.Message, error) {
	return Call[tgmodel.Message](ctx, c.token, "sendMessage", body)
}

func (c *Client) GetChat(ctx context.Context, chatID string) (*tgmodel.ChatFullInfo, error) {
	type getChatBody struct {
		ChatID string `json:"chat_id"`
	}
	return Call[tgmodel.ChatFullInfo](ctx, c.token, "getChat", getChatBody{ChatID: chatID})
}

func Call[T any](
	ctx context.Context,
	token string,
	method string,
	body any,
) (*T, error) {
	url := fmt.Sprintf(baseURL, token, method)
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

	tgResp, err := snet.UnmarshalResp[tgmodel.Response[T]](resp)
	if err != nil {
		return nil, err
	}
	return tgResp.Result, err
}
