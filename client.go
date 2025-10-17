package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/mrbanja/tg/model"
)

func GetWebhookInfo(ctx context.Context) (*model.WebhookInfo, error) {
	resp, err := send[*model.WebhookInfo](ctx, "getWebhookInfo", nil)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func SetWebhook(ctx context.Context, req model.SetWebhookRequest) error {
	resp, err := send[any](ctx, "setWebhook", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("setWebhook failed: %v", resp.Result)
	}
	return nil
}

func DeleteWebhook(ctx context.Context, req model.DeleteWebhookRequest) error {
	resp, err := send[any](ctx, "deleteWebhook", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("deleteWebhook failed: %v", resp.Result)
	}
	return nil
}

func SendPhoto(ctx context.Context, chatID int64, reader io.Reader) (*model.Response[any], error) {
	log := slog.With(slog.Int64("chat_id", chatID))
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if err := w.WriteField("chat_id", fmt.Sprintf("%d", chatID)); err != nil {
		log.Error("[*] writeField failed: ", "err", err)
		return nil, err
	}
	fileW, err := w.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		log.Error("[*] createFormFile failed: ", "err", err)
		return nil, err
	}
	if _, err := io.Copy(fileW, reader); err != nil {
		log.Error("[*] copy failed: ", "err", err)
	}
	if err := w.Close(); err != nil {
		log.Error("[*] close failed: ", "err", err)
		return nil, err
	}

	return do[any](ctx, "sendPhoto", &b, w.FormDataContentType())
}

func DeleteMessage(ctx context.Context, req model.DeleteMessageRequest) error {
	resp, err := send[any](ctx, "deleteMessage", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("deleteMessage failed: %v", resp.Result)
	}
	return nil
}

func SetMessageReaction(ctx context.Context, req model.SetMessageReactionRequest) error {
	resp, err := send[any](ctx, "setMessageReaction", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("setMessageReaction failed: %v", resp.Result)
	}
	return nil
}

func SendMessage(ctx context.Context, req model.SendMessageRequest) (*model.Message, error) {
	resp, err := send[*model.Message](ctx, "sendMessage", req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("sendMessage failed: %v", resp.Result)
	}
	return resp.Result, nil
}

func EditMessageText(ctx context.Context, req model.EditMessageTextRequest) error {
	resp, err := send[any](ctx, "editMessageText", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("editMessageText failed: %v", resp.Result)
	}
	return nil
}

func send[T any](ctx context.Context, method string, obj any) (*model.Response[T], error) {
	if token == "" {
		slog.Error("[*] token is empty")
		return nil, fmt.Errorf("token is empty")
	}

	body, err := json.Marshal(obj)
	if err != nil {
		slog.Error("[*] marshal obj failed: ", "err", err)
		return nil, err
	}

	return do[T](ctx, method, bytes.NewReader(body), "application/json")
}

func do[T any](
	ctx context.Context,
	method string,
	body io.Reader,
	contentType string,
) (*model.Response[T], error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, buildURL(method), body)
	if err != nil {
		slog.Error("[*] new request failed: ", "err", err)
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("[*] do request failed: ", "err", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		slog.Error("[*] status code is not 200", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("status code is not 200")
	}

	var res model.Response[T]
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func buildURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)
}
