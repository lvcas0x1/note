package httpx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"notes/internal/logx"
)

func init() { gin.SetMode(gin.TestMode) }

func TestRequestLoggerDropsSearchText(t *testing.T) {
	var buf bytes.Buffer
	r := gin.New()
	r.Use(RequestLogger(logx.New(&buf, slog.LevelInfo)))
	r.GET("/api/v1/search", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/v1/search?q=secret+topic&limit=5", nil))

	out := buf.String()
	if strings.Contains(out, "secret") {
		t.Fatalf("log contains search text: %s", out)
	}
	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatal(err)
	}
	if rec["query"] != "limit=5" {
		t.Errorf("query = %v, want limit=5", rec["query"])
	}
	if rec["status"] != float64(200) {
		t.Errorf("status = %v", rec["status"])
	}
}

func TestRecoveryReturnsProblem(t *testing.T) {
	var buf bytes.Buffer
	r := gin.New()
	r.Use(Recovery(logx.New(&buf, slog.LevelInfo)))
	r.GET("/boom", func(*gin.Context) { panic("kaboom") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("content type = %q", ct)
	}
	var p Problem
	if err := json.Unmarshal(w.Body.Bytes(), &p); err != nil || p.Status != 500 || p.Title == "" {
		t.Errorf("body = %s", w.Body.String())
	}
	if strings.Contains(w.Body.String(), "kaboom") {
		t.Error("panic value leaked to the client")
	}
	if !strings.Contains(buf.String(), "kaboom") {
		t.Error("panic not logged")
	}
}
