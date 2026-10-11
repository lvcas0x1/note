package export

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// pagesServer mimics a static host such as GitHub Pages: it serves files and
// folder index.html files under prefix, and 404.html with status 404 otherwise.
func pagesServer(t *testing.T, root, prefix string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, prefix) {
			http.NotFound(w, r)
			return
		}
		rel := path.Clean("/" + strings.TrimPrefix(r.URL.Path, prefix))
		file := filepath.Join(root, filepath.FromSlash(rel))
		if info, err := os.Stat(file); err == nil && info.IsDir() {
			file = filepath.Join(file, "index.html")
		}
		if raw, err := os.ReadFile(file); err == nil {
			if strings.HasSuffix(file, ".html") {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
			}
			_, _ = w.Write(raw)
			return
		}
		raw, _ := os.ReadFile(filepath.Join(root, "404.html"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(raw)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func getBody(t *testing.T, url string) (int, string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

func TestStaticSiteUnderSubPath(t *testing.T) {
	out, _, err := runExport(t, buildContent(t, nil), nil)
	if err != nil {
		t.Fatal(err)
	}
	srv := pagesServer(t, out, "/notes/")

	// Direct load and reload of an article URL get the app shell.
	for range 2 {
		status, body := getBody(t, srv.URL+"/notes/a/aws/s3/s3/")
		if status != http.StatusOK || !strings.Contains(body, `<div id="root">`) {
			t.Fatalf("article URL: %d %.120s", status, body)
		}
	}
	// Every URL in the shell is under the base path and loads.
	_, shell := getBody(t, srv.URL+"/notes/")
	for _, ref := range []string{`src="/notes/config.js"`, `src="/notes/static/app.js"`, `href="/notes/static/app.css"`} {
		if !strings.Contains(shell, ref) {
			t.Errorf("shell missing %s", ref)
		}
	}
	for _, p := range []string{"config.js", "static/app.js", "static/app.css", "data/sidebar.json", "data/articles/aws/s3/s3.json", "data/search.db"} {
		if status, _ := getBody(t, srv.URL+"/notes/"+p); status != http.StatusOK {
			t.Errorf("/notes/%s: %d", p, status)
		}
	}
	_, article := getBody(t, srv.URL+"/notes/data/articles/aws/s3/s3.json")
	i := strings.Index(article, `/notes/assets/`)
	if i < 0 {
		t.Fatalf("article has no asset URL under the base path: %s", article)
	}
	assetURL := article[i : i+strings.IndexAny(article[i:], `"\`)]
	if status, _ := getBody(t, srv.URL+assetURL); status != http.StatusOK {
		t.Errorf("asset %s: %d", assetURL, status)
	}
	// config.js carries the static mode, base path, and video host.
	_, cfg := getBody(t, srv.URL+"/notes/config.js")
	if !strings.Contains(cfg, `"mode":"static"`) || !strings.Contains(cfg, `"basePath":"/notes/"`) || !strings.Contains(cfg, `videos.example.com`) {
		t.Errorf("config.js = %s", cfg)
	}
	// Unknown paths get the app shell from 404.html, which renders the not-found view.
	status, body := getBody(t, srv.URL+"/notes/no/such/page")
	if status != http.StatusNotFound || !strings.Contains(body, `<div id="root">`) {
		t.Errorf("unknown path: %d %.120s", status, body)
	}
}

func TestCSPMetaOnEveryPage(t *testing.T) {
	out, _, err := runExport(t, buildContent(t, nil), nil)
	if err != nil {
		t.Fatal(err)
	}
	pages := 0
	_ = filepath.Walk(out, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".html") {
			return err
		}
		pages++
		raw, _ := os.ReadFile(p)
		s := string(raw)
		if !strings.Contains(s, `<meta http-equiv="Content-Security-Policy" content="default-src 'self';`) {
			t.Errorf("%s has no CSP meta tag", p)
		}
		if !strings.Contains(s, "media-src 'self' https://videos.example.com;") {
			t.Errorf("%s: media-src does not allow the video origin", p)
		}
		if strings.Contains(s, "unsafe-eval") || strings.Contains(s, "script-src 'self' 'unsafe-inline'") {
			t.Errorf("%s allows inline or eval scripts", p)
		}
		if strings.Index(s, "Content-Security-Policy") > strings.Index(s, "<script") {
			t.Errorf("%s: CSP meta must come before the first script", p)
		}
		return nil
	})
	if pages < 4 {
		t.Fatalf("only %d HTML pages", pages)
	}
}
