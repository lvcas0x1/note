package httpx

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"

	"notes/internal/logx"
)

func searchServer(t *testing.T, perMinute int, trusted ...string) http.Handler {
	t.Helper()
	e, err := NewEngine(logx.New(io.Discard, slog.LevelError), trusted...)
	if err != nil {
		t.Fatal(err)
	}
	api := humagin.New(e, huma.DefaultConfig("t", "1"))
	l := NewRateLimiter(perMinute)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/api/v1/search",
		Middlewares: huma.Middlewares{l.HumaMiddleware}},
		func(context.Context, *struct{}) (*struct{}, error) { return nil, nil })
	return e
}

func hit(h http.Handler, remote, forwarded string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search", nil)
	req.RemoteAddr = remote
	if forwarded != "" {
		req.Header.Set("X-Forwarded-For", forwarded)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestSixtyFirstRequestIs429(t *testing.T) {
	h := searchServer(t, 60)
	for i := range 60 {
		if w := hit(h, "203.0.113.5:1234", ""); w.Code >= 400 {
			t.Fatalf("request %d: %d", i+1, w.Code)
		}
	}
	w := hit(h, "203.0.113.5:1234", "")
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("61st request: %d, want 429", w.Code)
	}
	if ra, err := strconv.Atoi(w.Header().Get("Retry-After")); err != nil || ra < 1 {
		t.Errorf("Retry-After = %q", w.Header().Get("Retry-After"))
	}
	if w.Header().Get("Content-Type") != "application/problem+json" {
		t.Errorf("content type = %q", w.Header().Get("Content-Type"))
	}
	if w := hit(h, "203.0.113.6:1234", ""); w.Code >= 400 {
		t.Errorf("another IP was limited: %d", w.Code)
	}
}

func TestForwardedForFromUntrustedPeerIsIgnored(t *testing.T) {
	h := searchServer(t, 2)
	hit(h, "198.51.100.7:1", "1.1.1.1")
	hit(h, "198.51.100.7:1", "2.2.2.2")
	// A spoofed header cannot buy a fresh bucket: the connection IP is used.
	if w := hit(h, "198.51.100.7:1", "3.3.3.3"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("spoofed X-Forwarded-For bypassed the limit: %d", w.Code)
	}
}

func TestForwardedForFromTrustedProxyIsHonored(t *testing.T) {
	h := searchServer(t, 2, "10.0.0.0/8")
	hit(h, "10.0.0.1:1", "1.1.1.1")
	hit(h, "10.0.0.1:1", "1.1.1.1")
	if w := hit(h, "10.0.0.1:1", "1.1.1.1"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("client behind proxy not limited: %d", w.Code)
	}
	if w := hit(h, "10.0.0.1:1", "2.2.2.2"); w.Code >= 400 {
		t.Fatalf("a different client behind the trusted proxy was limited: %d", w.Code)
	}
}

func TestIdleClientsAreEvicted(t *testing.T) {
	l := NewRateLimiter(60)
	l.maxClients = 3
	now := time.Unix(0, 0)
	l.now = func() time.Time { return now }
	for _, ip := range []string{"a", "b", "c"} {
		l.Allow(ip)
	}
	now = now.Add(time.Hour)
	l.Allow("d")
	if n := l.Len(); n != 1 {
		t.Fatalf("tracked clients = %d, want only the new one", n)
	}
	for _, ip := range []string{"e", "f", "g", "h"} {
		l.Allow(ip)
	}
	if n := l.Len(); n > 3 {
		t.Fatalf("tracked clients = %d, want at most 3", n)
	}
}
