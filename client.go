package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/mrbanja/tg/v4/model"
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

func SendRichMessage(ctx context.Context, req model.SendRichMessageRequest) (*model.Message, error) {
	resp, err := send[*model.Message](ctx, "sendRichMessage", req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("sendRichMessage failed: %v", resp.Result)
	}
	return resp.Result, nil
}

func SendRichMessageDraft(ctx context.Context, req model.SendRichMessageDraftRequest) error {
	resp, err := send[any](ctx, "sendRichMessageDraft", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("sendRichMessageDraft failed: %v", resp.Result)
	}
	return nil
}

// SendRichMessageUpload sends a rich message whose media reference new files
// uploaded in the same request: set InputMedia*.Media to model.AttachURI(f.Field)
// for every file f.
func SendRichMessageUpload(ctx context.Context, req model.SendRichMessageRequest, files ...InputFile) (*model.Message, error) {
	resp, err := sendMultipart[*model.Message](ctx, "sendRichMessage", req, files)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("sendRichMessage failed: %v", resp.Result)
	}
	return resp.Result, nil
}

// SendMessageDraft streams a partial plain-text message to a private chat.
func SendMessageDraft(ctx context.Context, req model.SendMessageDraftRequest) error {
	resp, err := send[any](ctx, "sendMessageDraft", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("sendMessageDraft failed: %v", resp.Result)
	}
	return nil
}

// SendVoice sends a voice note referenced by file_id or HTTP URL (req.Voice).
func SendVoice(ctx context.Context, req model.SendVoiceRequest) (*model.Message, error) {
	resp, err := send[*model.Message](ctx, "sendVoice", req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("sendVoice failed: %v", resp.Result)
	}
	return resp.Result, nil
}

// SendVoiceUpload uploads and sends a new voice note (OGG/OPUS, MP3 or M4A).
// req.Voice is ignored.
func SendVoiceUpload(ctx context.Context, req model.SendVoiceRequest, fileName string, voice io.Reader) (*model.Message, error) {
	req.Voice = ""
	files := []InputFile{{Field: "voice", FileName: fileName, Reader: voice}}
	resp, err := sendMultipart[*model.Message](ctx, "sendVoice", req, files)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("sendVoice failed: %v", resp.Result)
	}
	return resp.Result, nil
}

func SendChatAction(ctx context.Context, req model.SendChatActionRequest) error {
	resp, err := send[any](ctx, "sendChatAction", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("sendChatAction failed: %v", resp.Result)
	}
	return nil
}

func AnswerCallbackQuery(ctx context.Context, req model.AnswerCallbackQueryRequest) error {
	resp, err := send[any](ctx, "answerCallbackQuery", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("answerCallbackQuery failed: %v", resp.Result)
	}
	return nil
}

func SetMyCommands(ctx context.Context, req model.SetMyCommandsRequest) error {
	resp, err := send[any](ctx, "setMyCommands", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("setMyCommands failed: %v", resp.Result)
	}
	return nil
}

func EditMessageReplyMarkup(ctx context.Context, req model.EditMessageReplyMarkupRequest) error {
	resp, err := send[any](ctx, "editMessageReplyMarkup", req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("editMessageReplyMarkup failed: %v", resp.Result)
	}
	return nil
}

// GetFile prepares a file for downloading with DownloadFile.
func GetFile(ctx context.Context, fileID string) (*model.File, error) {
	resp, err := send[*model.File](ctx, "getFile", model.GetFileRequest{FileID: fileID})
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("getFile failed: %v", resp.Result)
	}
	return resp.Result, nil
}

// MaxDownloadSize is the largest file the Bot API lets bots download.
const MaxDownloadSize = 20 << 20

// DownloadFile downloads a file prepared by GetFile. The download URL contains
// the bot token, so it is never exposed.
func DownloadFile(ctx context.Context, file *model.File) ([]byte, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	if file == nil || file.FilePath == "" {
		return nil, fmt.Errorf("download file: empty file path")
	}

	fileURL, err := url.JoinPath(apiBaseURL, "file", "bot"+token, file.FilePath)
	if err != nil {
		return nil, fmt.Errorf("download file: build url: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("download file: new request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		// url.Error embeds the URL, which contains the token.
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return nil, fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download file: status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, MaxDownloadSize+1))
	if err != nil {
		return nil, fmt.Errorf("download file: read body: %w", err)
	}
	if len(data) > MaxDownloadSize {
		return nil, fmt.Errorf("download file: larger than %d bytes", MaxDownloadSize)
	}
	return data, nil
}

// Call invokes any Bot API method with a JSON request and decodes its result.
// Use it for methods that have no dedicated wrapper.
func Call[T any](ctx context.Context, method string, req any) (T, error) {
	var zero T
	resp, err := send[T](ctx, method, req)
	if err != nil {
		return zero, err
	}
	if !resp.Ok {
		return zero, fmt.Errorf("%s failed: %v", method, resp.Result)
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

	resp, err := httpClient.Do(req)
	if err != nil {
		// url.Error embeds the request URL, which contains the token.
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			err = fmt.Errorf("%s: %w", method, urlErr.Err)
		}
		slog.Error("[*] do request failed: ", "err", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		slog.Error("[*] status code is not 200", "status", resp.StatusCode, "body", string(respBody))
		return nil, newAPIError(method, resp.StatusCode, respBody)
	}

	var res model.Response[T]
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func buildURL(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", apiBaseURL, token, method)
}

// InputFile is a file uploaded with multipart/form-data. Field is the
// multipart part name: the method parameter (e.g. "voice") for top-level
// files, or the name referenced as model.AttachURI(Field) inside InputMedia.
type InputFile struct {
	Field    string
	FileName string
	Reader   io.Reader
}

// sendMultipart sends req as multipart/form-data: every top-level JSON field
// becomes a form field (objects and arrays stay JSON-serialized), followed by files.
func sendMultipart[T any](ctx context.Context, method string, req any, files []InputFile) (*model.Response[T], error) {
	if token == "" {
		slog.Error("[*] token is empty")
		return nil, fmt.Errorf("token is empty")
	}

	body, contentType, err := buildMultipart(req, files)
	if err != nil {
		return nil, fmt.Errorf("%s: build multipart body: %w", method, err)
	}
	return do[T](ctx, method, body, contentType)
}

func buildMultipart(req any, files []InputFile) (*bytes.Buffer, string, error) {
	raw, err := json.Marshal(req)
	if err != nil {
		return nil, "", fmt.Errorf("marshal request: %w", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, "", fmt.Errorf("request must be a JSON object: %w", err)
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for _, k := range keys {
		v := fields[k]
		if string(v) == "null" {
			continue
		}
		value := string(v)
		if strings.HasPrefix(value, `"`) {
			if err := json.Unmarshal(v, &value); err != nil {
				return nil, "", fmt.Errorf("field %s: %w", k, err)
			}
		}
		if err := w.WriteField(k, value); err != nil {
			return nil, "", fmt.Errorf("write field %s: %w", k, err)
		}
	}
	for _, f := range files {
		if f.Field == "" || f.Reader == nil {
			return nil, "", fmt.Errorf("file %q: field name and reader are required", f.FileName)
		}
		fileName := f.FileName
		if fileName == "" {
			fileName = f.Field
		}
		part, err := w.CreateFormFile(f.Field, fileName)
		if err != nil {
			return nil, "", fmt.Errorf("create form file %s: %w", f.Field, err)
		}
		if _, err := io.Copy(part, f.Reader); err != nil {
			return nil, "", fmt.Errorf("copy file %s: %w", f.Field, err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", fmt.Errorf("close multipart writer: %w", err)
	}
	return &b, w.FormDataContentType(), nil
}

// APIError is returned when the Bot API responds with a non-200 status.
type APIError struct {
	Method      string
	StatusCode  int
	ErrorCode   int
	Description string
	// RetryAfter is set when flood control was exceeded (HTTP 429).
	RetryAfter time.Duration
	body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: status %d: %s", e.Method, e.StatusCode, bytes.TrimSpace(e.body))
}

func newAPIError(method string, status int, body []byte) *APIError {
	apiErr := &APIError{Method: method, StatusCode: status, body: body}
	var resp model.Response[json.RawMessage]
	if err := json.Unmarshal(body, &resp); err == nil {
		apiErr.ErrorCode = resp.ErrorCode
		apiErr.Description = resp.Description
		if resp.Parameters != nil && resp.Parameters.RetryAfter > 0 {
			apiErr.RetryAfter = time.Duration(resp.Parameters.RetryAfter) * time.Second
		}
	}
	return apiErr
}

// RetryAfter reports how long to wait before repeating a request that hit
// Telegram flood control.
func RetryAfter(err error) (time.Duration, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.RetryAfter > 0 {
		return apiErr.RetryAfter, true
	}
	return 0, false
}
