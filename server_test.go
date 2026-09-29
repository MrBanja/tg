package tg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const testUpdate = `{"update_id":1,"message":{"message_id":2,"date":1,"chat":{"id":3,"type":"private"},"text":"hi"}}`

func newTestServer(opts ...ServerOptionFunc) *Server {
	s := NewServer(":0", "https://example.com", opts...)
	s.webhookPath = "/hook"
	return s
}

func TestWebhookRejectsInvalidSecretToken(t *testing.T) {
	s := newTestServer(WithSecretToken("s3cret"))
	called := false
	s.Register(func(ctx context.Context, req *Request) { called = true })
	h := s.updateHandler(context.Background(), nil)

	req := httptest.NewRequest(http.MethodPost, "/hook", strings.NewReader(testUpdate))
	req.Header.Set(SecretTokenHeader, "wrong")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized || called {
		t.Fatalf("expected 401 without dispatch, got %d (called=%v)", rec.Code, called)
	}

	req = httptest.NewRequest(http.MethodPost, "/hook", strings.NewReader(testUpdate))
	req.Header.Set(SecretTokenHeader, "s3cret")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !called {
		t.Fatalf("expected 200 with dispatch, got %d (called=%v)", rec.Code, called)
	}
}

func TestWebhookRejectsNonPost(t *testing.T) {
	s := newTestServer()
	h := s.updateHandler(context.Background(), nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/hook", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestDetachedHandlersOutliveRequest(t *testing.T) {
	s := newTestServer(WithDetachedHandlers())
	release := make(chan struct{})
	done := make(chan error, 1)
	s.Register(func(ctx context.Context, req *Request) {
		<-release
		done <- ctx.Err()
	})

	serverCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h := s.updateHandler(serverCtx, nil)

	reqCtx, cancelReq := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodPost, "/hook", strings.NewReader(testUpdate)).WithContext(reqCtx)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req) // must return before the handler finishes
	cancelReq()
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	close(release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("handler context cancelled with the request: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("handler did not run")
	}
	s.inflight.Wait()
}
