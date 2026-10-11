package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"notes/internal/content"
	"notes/internal/httpx"
	"notes/internal/logx"
	"notes/internal/store"
	"notes/internal/store/storetest"
)

type fixedBuild struct {
	r   *store.Reader
	err error
}

func (f fixedBuild) Reader() (*store.Reader, error) { return f.r, f.err }

func newServer(t *testing.T) (*httptest.Server, []*content.Article) {
	t.Helper()
	data := t.TempDir()
	articles := storetest.Articles(7, 40)
	articles[0].ID = "aws/s3/s3"
	articles[0].Path = "aws/s3/s3.md"
	if _, err := store.WriteBuild(data, storetest.Result(articles)); err != nil {
		t.Fatal(err)
	}
	r, err := store.OpenCurrent(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = r.Close() })
	e, _ := httpx.NewEngine(logx.New(io.Discard, slog.LevelError))
	Register(e, fixedBuild{r: r})
	srv := httptest.NewServer(e)
	t.Cleanup(srv.Close)
	return srv, articles
}

func get(t *testing.T, srv *httptest.Server, path string) (*http.Response, []byte) {
	t.Helper()
	return do(t, srv, http.MethodGet, path)
}

func do(t *testing.T, srv *httptest.Server, method, path string) (*http.Response, []byte) {
	t.Helper()
	req, _ := http.NewRequest(method, srv.URL+path, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	if len(resp.Cookies()) > 0 || resp.Header.Get("Set-Cookie") != "" {
		t.Errorf("%s %s set a cookie", method, path)
	}
	return resp, body
}

func expectProblem(t *testing.T, resp *http.Response, body []byte, status int, contains string) {
	t.Helper()
	if resp.StatusCode != status {
		t.Fatalf("status = %d, want %d (%s)", resp.StatusCode, status, body)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/problem+json") {
		t.Errorf("content type = %q", ct)
	}
	var p map[string]any
	if err := json.Unmarshal(body, &p); err != nil {
		t.Fatalf("body is not JSON: %s", body)
	}
	if p["status"] != float64(status) || p["title"] == nil || p["detail"] == nil {
		t.Errorf("problem = %s", body)
	}
	if contains != "" && !strings.Contains(string(body), contains) {
		t.Errorf("problem %s does not mention %q", body, contains)
	}
}

func TestEndpoints(t *testing.T) {
	srv, articles := newServer(t)

	resp, body := get(t, srv, "/api/v1/sidebar")
	if resp.StatusCode != 200 || !strings.Contains(string(body), `"title":"Generated"`) {
		t.Fatalf("sidebar %d %s", resp.StatusCode, body)
	}

	resp, body = get(t, srv, "/api/v1/articles/"+url.PathEscape("aws/s3/s3"))
	var a store.ArticleDetail
	if resp.StatusCode != 200 || json.Unmarshal(body, &a) != nil || a.ID != "aws/s3/s3" || a.Title != articles[0].Title || a.HTML == "" {
		t.Fatalf("article %d %s", resp.StatusCode, body)
	}
	for _, field := range []string{`"date"`, `"summary"`, `"tags"`, `"category"`, `"html"`} {
		if !strings.Contains(string(body), field) {
			t.Errorf("article response missing %s", field)
		}
	}

	resp, body = get(t, srv, "/api/v1/articles/does%2Fnot%2Fexist")
	expectProblem(t, resp, body, 404, "does/not/exist")

	resp, body = get(t, srv, "/api/v1/calendar?from=2026-10-01&to=2026-10-31")
	var days []store.DayCount
	if resp.StatusCode != 200 || json.Unmarshal(body, &days) != nil || len(days) == 0 {
		t.Fatalf("calendar %d %s", resp.StatusCode, body)
	}

	resp, body = get(t, srv, "/api/v1/days/2026-10-05")
	if resp.StatusCode != 200 || !strings.HasPrefix(string(body), "[") {
		t.Fatalf("day %d %s", resp.StatusCode, body)
	}

	resp, body = get(t, srv, "/api/v1/search?q=kafka")
	if resp.StatusCode != 200 || !strings.Contains(string(body), `"snippet"`) {
		t.Fatalf("search %d %s", resp.StatusCode, body)
	}

	resp, body = get(t, srv, "/api/v1/openapi.json")
	if resp.StatusCode != 200 || !strings.Contains(string(body), `"openapi":"3.1`) || strings.Contains(string(body), "securitySchemes") {
		t.Fatalf("openapi %d %.200s", resp.StatusCode, body)
	}
	if strings.Contains(string(body), "$schema") {
		t.Error("responses advertise $schema links")
	}
}

func TestBadRequests(t *testing.T) {
	srv, _ := newServer(t)
	cases := []struct{ path, mention string }{
		{"/api/v1/days/2026-13-40", "date"},
		{"/api/v1/calendar?from=2026-10-01&to=2027-10-03", "query.to"},
		{"/api/v1/calendar?from=2026-10-31&to=2026-10-01", "query.from"},
		{"/api/v1/calendar?from=2026-10-01", "to"},
		{"/api/v1/search?q=" + strings.Repeat("a", 201), "q"},
		{"/api/v1/search?q=ab", "q"},
	}
	for _, c := range cases {
		resp, body := get(t, srv, c.path)
		expectProblem(t, resp, body, 400, c.mention)
	}
	if resp, body := get(t, srv, "/api/v1/calendar?from=2026-01-01&to=2027-01-01"); resp.StatusCode != 200 {
		t.Errorf("a 366-day range must be accepted: %d %s", resp.StatusCode, body)
	}
}

func TestWriteMethodsAreRejected(t *testing.T) {
	srv, _ := newServer(t)
	for _, m := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		resp, body := do(t, srv, m, "/api/v1/sidebar")
		expectProblem(t, resp, body, 405, "")
	}
}

func TestNoBuildIs503(t *testing.T) {
	e, _ := httpx.NewEngine(logx.New(io.Discard, slog.LevelError))
	Register(e, fixedBuild{err: store.ErrNoBuild})
	srv := httptest.NewServer(e)
	defer srv.Close()
	resp, body := get(t, srv, "/api/v1/sidebar")
	expectProblem(t, resp, body, 503, "built")
}
