package content

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newRemoteServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/a.drawio", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("<mxfile/>")) })
	mux.HandleFunc("/big.json", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", 2048)))
	})
	mux.HandleFunc("/slow.json", func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
		case <-r.Context().Done():
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func buildWith(t *testing.T, body string, extra map[string]string, f func(*Fetcher)) (*Result, error) {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"sidebar.json": `[{"title":"A","path":"a.md"}]`,
		"a.md":         "---\ndate: 2026-01-01\n---\n" + body + "\n",
	}
	for k, v := range extra {
		files[k] = v
	}
	writeFiles(t, root, files)
	fetcher := NewFetcher(t.TempDir())
	if f != nil {
		f(fetcher)
	}
	return Build(root, Options{Fetcher: fetcher})
}

func TestRemoteDrawioIsDownloaded(t *testing.T) {
	srv := newRemoteServer(t)
	res, err := buildWith(t, "![drawio]("+srv.URL+"/a.drawio)", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	html := res.Articles[0].HTML
	if strings.Contains(html, srv.URL) || strings.Contains(html, "127.0.0.1") {
		t.Fatalf("compiled HTML still references the remote host: %s", html)
	}
	if !strings.Contains(html, `data-embed="drawio" data-src="/assets/`) || len(res.Assets) != 1 || res.Assets[0].Name != "a.drawio" {
		t.Fatalf("html = %s assets = %+v", html, res.Assets)
	}
	got, _ := os.ReadFile(res.Assets[0].LocalPath)
	if string(got) != "<mxfile/>" {
		t.Fatalf("downloaded content = %q", got)
	}
}

func TestRemoteFailures(t *testing.T) {
	srv := newRemoteServer(t)
	cases := []struct {
		name, body, want string
		tweak            func(*Fetcher)
	}{
		{"404", "![drawio](" + srv.URL + "/missing.drawio)", "HTTP 404", nil},
		{"oversize", "![openapi](" + srv.URL + "/big.json)", "exceeds", func(f *Fetcher) { f.MaxBytes = 1024 }},
		{"timeout", "![openapi](" + srv.URL + "/slow.json)", "timed out", func(f *Fetcher) { f.Timeout = 100 * time.Millisecond }},
		{"unreachable", "![drawio](http://127.0.0.1:1/x.drawio)", "download failed", nil},
		{"file scheme", "![openapi](file:///etc/hosts)", "unsupported URL scheme", nil},
		{"remote video", "![mp4](https://example.com/a.mp4)", "videos must be local files", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := buildWith(t, tc.body, nil, tc.tweak)
			if err == nil || !strings.Contains(err.Error(), tc.want) || !strings.Contains(err.Error(), "a.md") {
				t.Fatalf("err = %v, want %q naming a.md", err, tc.want)
			}
		})
	}
}

func TestLocalAssetRules(t *testing.T) {
	cases := []struct {
		name, body, want string
		extra            map[string]string
	}{
		{"missing local asset", "![excalidraw](./missing.excalidraw)", `"missing.excalidraw" does not exist`, nil},
		{"path traversal", "![drawio](../../etc/passwd)", "outside the content folder", nil},
		{"wrong video extension", "![mp4](./clip.mov)", ".mp4 file", map[string]string{"clip.mov": "x"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := buildWith(t, tc.body, tc.extra, nil)
			if err == nil || !strings.Contains(err.Error(), tc.want) || !strings.Contains(err.Error(), "a.md") {
				t.Fatalf("err = %v, want %q naming a.md", err, tc.want)
			}
		})
	}
}

func TestLargeLocalVideoHasNoSizeLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("hashes a 2 GB sparse file")
	}
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"sidebar.json": `[{"title":"A","path":"a.md"}]`,
		"a.md":         "---\ndate: 2026-01-01\n---\n![mp4](videos/big.mp4)\n",
	})
	if err := os.MkdirAll(filepath.Join(root, "videos"), 0o750); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(filepath.Join(root, "videos", "big.mp4"))
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(2 << 30); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	res, err := Build(root, Options{})
	if err != nil {
		t.Fatalf("2 GB video failed: %v", err)
	}
	if len(res.Assets) != 1 || res.Assets[0].Size != 2<<30 || res.Assets[0].ContentType != "video/mp4" {
		t.Fatalf("assets = %+v", res.Assets)
	}
}

func TestContentHashedNames(t *testing.T) {
	res, err := buildWith(t, "![a](img/one.png)\n\n![b](img/two.png)", map[string]string{"img/one.png": "1", "img/two.png": "2"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Assets) != 2 || res.Assets[0].URL == res.Assets[1].URL {
		t.Fatalf("assets = %+v", res.Assets)
	}
	for _, a := range res.Assets {
		if !strings.HasPrefix(a.URL, "/assets/"+a.Hash[:16]+"/") {
			t.Errorf("url %s is not content-hashed (%s)", a.URL, a.Hash)
		}
	}
}
