package main

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"notes/internal/config"
	"notes/internal/logx"
	"notes/internal/store"
	"notes/internal/store/storetest"
)

func startServer(t *testing.T, data string) string {
	t.Helper()
	cfg := config.Server{Listen: "127.0.0.1:0", SearchRate: 60, DataDir: data}
	srv, closeFn, err := newHTTPServer(cfg, logx.New(io.Discard, slog.LevelError))
	if err != nil {
		t.Fatal(err)
	}
	ln, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = srv.Serve(ln) }()
	t.Cleanup(func() { _ = srv.Close(); closeFn() })
	return "http://" + ln.Addr().String()
}

func fetch(t *testing.T, method, url string) (*http.Response, string) {
	t.Helper()
	req, _ := http.NewRequest(method, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, _ := io.ReadAll(resp.Body)
	return resp, string(b)
}

func TestSecurityHeadersEverywhere(t *testing.T) {
	data := t.TempDir()
	if _, err := store.WriteBuild(data, storetest.Result(storetest.Articles(1, 5))); err != nil {
		t.Fatal(err)
	}
	base := startServer(t, data)
	for _, p := range []string{"/", "/a/some/article", "/api/v1/sidebar", "/healthz", "/api/nope"} {
		resp, _ := fetch(t, http.MethodGet, base+p)
		h := resp.Header
		csp := h.Get("Content-Security-Policy")
		if !strings.Contains(csp, "default-src 'self'") || !strings.Contains(csp, "frame-ancestors 'none'") || strings.Contains(csp, "http") {
			t.Errorf("%s: CSP = %q", p, csp)
		}
		if h.Get("Strict-Transport-Security") == "" || h.Get("X-Content-Type-Options") != "nosniff" || h.Get("Referrer-Policy") != "strict-origin-when-cross-origin" {
			t.Errorf("%s: headers = %v", p, h)
		}
	}
}

func TestRoutingRules(t *testing.T) {
	data := t.TempDir()
	if _, err := store.WriteBuild(data, storetest.Result(storetest.Articles(1, 5))); err != nil {
		t.Fatal(err)
	}
	base := startServer(t, data)
	// The web app may not be built in this checkout; the fallback then answers 503.
	resp, _ := fetch(t, http.MethodGet, base+"/a/aws/s3/s3")
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("SPA fallback status = %d", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusOK && !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		t.Errorf("SPA fallback content type = %q", resp.Header.Get("Content-Type"))
	}
	if resp, _ := fetch(t, http.MethodGet, base+"/api/unknown"); resp.StatusCode != http.StatusNotFound {
		t.Errorf("unknown API route = %d, want 404", resp.StatusCode)
	}
	for _, m := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		for _, p := range []string{"/", "/a/x", "/api/v1/sidebar", "/healthz"} {
			if resp, _ := fetch(t, m, base+p); resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("%s %s = %d, want 405", m, p, resp.StatusCode)
			}
		}
	}
	if resp, body := fetch(t, http.MethodGet, base+"/config.js"); resp.StatusCode != http.StatusOK || !strings.Contains(body, `"mode":"server"`) {
		t.Errorf("config.js = %d %s", resp.StatusCode, body)
	}
	if resp, _ := fetch(t, http.MethodHead, base+"/healthz"); resp.StatusCode != http.StatusOK {
		t.Errorf("HEAD /healthz = %d", resp.StatusCode)
	}
}

func TestHealthAndHotReload(t *testing.T) {
	data := t.TempDir()
	base := startServer(t, data)
	if resp, _ := fetch(t, http.MethodGet, base+"/healthz"); resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("health before any build = %d, want 503", resp.StatusCode)
	}
	first := storetest.Articles(1, 3)
	if _, err := store.WriteBuild(data, storetest.Result(first)); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool {
		resp, _ := fetch(t, http.MethodGet, base+"/healthz")
		return resp.StatusCode == http.StatusOK
	})

	second := storetest.Articles(9, 3)
	second[0].Title = "Freshly rebuilt title"
	if _, err := store.WriteBuild(data, storetest.Result(second)); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool {
		_, body := fetch(t, http.MethodGet, base+"/api/v1/sidebar")
		return strings.Contains(body, "Freshly rebuilt title")
	})
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("condition not met within 5s")
}

func TestDefaultListenIsLoopbackOnly(t *testing.T) {
	cfg, err := config.Load(func(string) string { return "" })
	if err != nil || cfg.Listen != "127.0.0.1:8080" {
		t.Fatalf("default listen = %q err = %v", cfg.Listen, err)
	}
	host, _, _ := net.SplitHostPort(cfg.Listen)
	ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			_ = c.Close()
		}
	}()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	ip := nonLoopbackIP()
	if ip == "" {
		t.Skip("no non-loopback interface")
	}
	if c, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), time.Second); err == nil {
		_ = c.Close()
		t.Fatalf("server bound to the default address is reachable on %s", ip)
	}
	if c, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), time.Second); err != nil {
		t.Fatalf("not reachable on loopback: %v", err)
	} else {
		_ = c.Close()
	}
}

func nonLoopbackIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if n, ok := a.(*net.IPNet); ok && !n.IP.IsLoopback() && n.IP.To4() != nil {
			return n.IP.String()
		}
	}
	return ""
}
