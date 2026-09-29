package tg

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mrbanja/tg/v4/model"
)

func setupTestAPI(t *testing.T, h http.HandlerFunc) {
	t.Helper()
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	prevBase, prevToken := apiBaseURL, token
	SetAPIBaseURL(ts.URL)
	SetToken("TEST")
	t.Cleanup(func() {
		SetAPIBaseURL(prevBase)
		SetToken(prevToken)
	})
}

func TestSendVoiceUploadSendsMultipart(t *testing.T) {
	setupTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botTEST/sendVoice" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if got := r.FormValue("chat_id"); got != "42" {
			t.Errorf("chat_id = %q", got)
		}
		if got := r.FormValue("duration"); got != "7" {
			t.Errorf("duration = %q", got)
		}
		if got := r.FormValue("reply_parameters"); got != `{"message_id":3}` {
			t.Errorf("reply_parameters = %q", got)
		}
		if _, ok := r.MultipartForm.Value["voice"]; ok {
			t.Errorf("voice must be a file part, not a field")
		}
		f, hdr, err := r.FormFile("voice")
		if err != nil {
			t.Fatalf("voice file: %v", err)
		}
		defer f.Close()
		data, _ := io.ReadAll(f)
		if string(data) != "OggS-data" || hdr.Filename != "guide.ogg" {
			t.Errorf("unexpected file %q (%s)", data, hdr.Filename)
		}
		_, _ = io.WriteString(w, `{"ok":true,"result":{"message_id":10,"date":1,"voice":{"file_id":"fid","file_unique_id":"u","duration":7}}}`)
	})

	msg, err := SendVoiceUpload(context.Background(), model.SendVoiceRequest{
		ChatID:          42,
		Duration:        7,
		ReplyParameters: &model.ReplyParameters{MessageID: 3},
	}, "guide.ogg", strings.NewReader("OggS-data"))
	if err != nil {
		t.Fatalf("SendVoiceUpload: %v", err)
	}
	if msg.MessageID != 10 || msg.Voice == nil || msg.Voice.FileID != "fid" {
		t.Fatalf("unexpected message: %+v", msg)
	}
}

func TestSendRichMessageUploadAttachesMedia(t *testing.T) {
	setupTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		var rm map[string]any
		if err := json.Unmarshal([]byte(r.FormValue("rich_message")), &rm); err != nil {
			t.Fatalf("rich_message is not JSON: %v", err)
		}
		if _, _, err := r.FormFile("audio1"); err != nil {
			t.Fatalf("attached file missing: %v", err)
		}
		_, _ = io.WriteString(w, `{"ok":true,"result":{"message_id":11,"date":1}}`)
	})

	_, err := SendRichMessageUpload(context.Background(), model.SendRichMessageRequest{
		ChatID: 1,
		RichMessage: model.InputRichMessage{Blocks: []model.InputRichBlock{
			model.InputRichBlockVoiceNote{VoiceNote: model.InputMediaVoiceNote{Media: model.AttachURI("audio1")}},
		}},
	}, InputFile{Field: "audio1", FileName: "a.ogg", Reader: strings.NewReader("x")})
	if err != nil {
		t.Fatalf("SendRichMessageUpload: %v", err)
	}
}

func TestAPIErrorExposesRetryAfter(t *testing.T) {
	setupTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"ok":false,"error_code":429,"description":"Too Many Requests: retry after 3","parameters":{"retry_after":3}}`)
	})

	_, err := SendMessage(context.Background(), model.SendMessageRequest{ChatID: 1, Text: "x"})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.ErrorCode != 429 || apiErr.Description == "" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
	if d, ok := RetryAfter(err); !ok || d != 3*time.Second {
		t.Fatalf("RetryAfter = %v, %v", d, ok)
	}
	if !strings.Contains(err.Error(), "sendMessage: status 429") {
		t.Fatalf("unexpected error text: %v", err)
	}
}

func TestGetFileAndDownloadFile(t *testing.T) {
	setupTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/botTEST/getFile":
			var req model.GetFileRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.FileID != "abc" {
				t.Errorf("file_id = %q", req.FileID)
			}
			_, _ = io.WriteString(w, `{"ok":true,"result":{"file_id":"abc","file_unique_id":"u","file_path":"photos/file_1.jpg"}}`)
		case "/file/botTEST/photos/file_1.jpg":
			_, _ = io.WriteString(w, "JPEGDATA")
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	})

	f, err := GetFile(context.Background(), "abc")
	if err != nil {
		t.Fatalf("GetFile: %v", err)
	}
	data, err := DownloadFile(context.Background(), f)
	if err != nil {
		t.Fatalf("DownloadFile: %v", err)
	}
	if string(data) != "JPEGDATA" {
		t.Fatalf("unexpected data %q", data)
	}
}

func TestCallDecodesResult(t *testing.T) {
	setupTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botTEST/getMe" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"Guide","username":"guide_bot"}}`)
	})

	me, err := Call[model.User](context.Background(), "getMe", struct{}{})
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if me.Username != "guide_bot" {
		t.Fatalf("unexpected user %+v", me)
	}
}
