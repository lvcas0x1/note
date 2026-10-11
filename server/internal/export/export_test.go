package export

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"notes/internal/content"
	"notes/internal/store"
)

var webFS = fstest.MapFS{
	"index.html":                         {Data: []byte(`<!doctype html><html><head><meta charset="utf-8"><script src="/config.js"></script><script type="module" src="/static/app.js"></script><link rel="stylesheet" href="/static/app.css"></head><body><div id="root"></div></body></html>`)},
	"static/app.js":                      {Data: []byte("console.log(1)")},
	"static/app.css":                     {Data: []byte("body{}")},
	"vendor/drawio/viewer-static.min.js": {Data: []byte("/* viewer */")},
	"vendor/drawio/LICENSE":              {Data: []byte("Apache-2.0")},
}

// buildContent builds a small content folder with a draft, an image, and a video.
func buildContent(t *testing.T, extra map[string]string) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"sidebar.json":       `[{"title":"Cloud","children":[{"title":"S3","path":"aws/s3/s3.md"},{"title":"Secret","path":"secret.md"}]}]`,
		"aws/s3/s3.md":       "---\ndate: 2026-10-05\n---\n# S3\n\nBuckets ![logo](/img/logo.png)\n\n![mp4](../../videos/demo.mp4)\n",
		"secret.md":          "---\ndate: 2026-10-05\ndraft: true\n---\n# Zebracorn secret\n\nHidden draft body ![x](img/secret.png)\n",
		"img/logo.png":       "png",
		"img/secret.png":     "secret-png",
		"videos/demo.mp4":    "video-bytes",
		"aws/s3/../../go.md": "",
	}
	delete(files, "aws/s3/../../go.md")
	for k, v := range extra {
		files[k] = v
	}
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		_ = os.MkdirAll(filepath.Dir(p), 0o750)
		if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	res, err := content.Build(root, content.Options{})
	if err != nil {
		t.Fatal(err)
	}
	data := t.TempDir()
	if _, err := store.WriteBuild(data, res); err != nil {
		t.Fatal(err)
	}
	return data
}

func runExport(t *testing.T, data string, mutate func(*Options)) (string, *Summary, error) {
	t.Helper()
	out := filepath.Join(t.TempDir(), "site")
	o := Options{DataDir: data, Out: out, BasePath: "/notes/", VideoBaseURL: "https://videos.example.com/notes", Web: webFS}
	if mutate != nil {
		mutate(&o)
	}
	s, err := Run(context.Background(), o)
	return out, s, err
}

func TestExportWritesStaticSite(t *testing.T) {
	out, s, err := runExport(t, buildContent(t, nil), nil)
	if err != nil {
		t.Fatal(err)
	}
	if s.Articles != 1 || s.Assets != 1 || s.Videos != 1 {
		t.Fatalf("summary = %+v", s)
	}
	for _, f := range []string{"index.html", "404.html", "config.js", "static/app.js", "vendor/drawio/viewer-static.min.js",
		"data/sidebar.json", "data/calendar.json", "data/days/2026-10-05.json", "data/articles/aws/s3/s3.json",
		"data/search.db", "a/aws/s3/s3/index.html", "d/2026-10-05/index.html"} {
		if _, err := os.Stat(filepath.Join(out, f)); err != nil {
			t.Errorf("missing %s", f)
		}
	}
	var a store.ArticleDetail
	raw, _ := os.ReadFile(filepath.Join(out, "data/articles/aws/s3/s3.json"))
	if err := json.Unmarshal(raw, &a); err != nil || a.ID != "aws/s3/s3" || a.Title != "S3" {
		t.Fatalf("article json = %s", raw)
	}
	if !strings.Contains(a.HTML, `src="/notes/assets/`) {
		t.Errorf("image URL not under the base path: %s", a.HTML)
	}
}

func TestNoBuild(t *testing.T) {
	_, _, err := runExport(t, t.TempDir(), nil)
	if err == nil || !strings.Contains(err.Error(), "notes build") {
		t.Fatalf("err = %v, want run-build-first", err)
	}
}

func TestDraftsStayOut(t *testing.T) {
	out, _, err := runExport(t, buildContent(t, nil), nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = filepath.WalkDir(out, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if strings.HasSuffix(p, ".md") || strings.HasSuffix(p, ".mp4") {
			t.Errorf("site output contains %s", p)
		}
		raw, _ := os.ReadFile(p)
		for _, secret := range []string{"Zebracorn", "Hidden draft body", "secret-png", "secret.png"} {
			if strings.Contains(string(raw), secret) {
				t.Errorf("%s contains draft content %q", p, secret)
			}
		}
		return nil
	})
}

func TestVideoURLRewrite(t *testing.T) {
	out, _, err := runExport(t, buildContent(t, nil), nil)
	if err != nil {
		t.Fatal(err)
	}
	var m Manifest
	raw, err := os.ReadFile(filepath.Join(VideosDir(out), "videos-manifest.json"))
	if err != nil || json.Unmarshal(raw, &m) != nil || len(m.Files) != 1 {
		t.Fatalf("manifest = %s err = %v", raw, err)
	}
	f := m.Files[0]
	got, err := os.ReadFile(filepath.Join(VideosDir(out), filepath.FromSlash(f.Path)))
	if err != nil || string(got) != "video-bytes" || f.Size != int64(len("video-bytes")) || len(f.SHA256) != 64 {
		t.Fatalf("video file %s: %q err = %v", f.Path, got, err)
	}
	article, _ := os.ReadFile(filepath.Join(out, "data/articles/aws/s3/s3.json"))
	want := `https://videos.example.com/notes/` + f.Path
	if !strings.Contains(string(article), want) {
		t.Fatalf("article does not point to %s: %s", want, article)
	}
}

func TestMissingVideoBaseURL(t *testing.T) {
	_, _, err := runExport(t, buildContent(t, nil), func(o *Options) { o.VideoBaseURL = "" })
	if err == nil || !strings.Contains(err.Error(), "aws/s3/s3") || !strings.Contains(err.Error(), "video base URL") {
		t.Fatalf("err = %v, want the article named", err)
	}
}

func TestSizeLimits(t *testing.T) {
	data := buildContent(t, map[string]string{"img/logo.png": strings.Repeat("x", 4096)})
	_, _, err := runExport(t, data, func(o *Options) { o.MaxFileBytes = 1024 })
	if err == nil || !strings.Contains(err.Error(), "logo.png") {
		t.Fatalf("err = %v, want the large file named", err)
	}
	_, _, err = runExport(t, data, func(o *Options) { o.MaxSiteBytes = 2048 })
	if err == nil || !strings.Contains(err.Error(), "bytes") {
		t.Fatalf("err = %v, want the total size", err)
	}
}

func TestOutputMustBeEmpty(t *testing.T) {
	data := buildContent(t, nil)
	out := filepath.Join(t.TempDir(), "site")
	_ = os.MkdirAll(out, 0o750)
	_ = os.WriteFile(filepath.Join(out, "keep.txt"), []byte("x"), 0o600)
	if _, err := Run(context.Background(), Options{DataDir: data, Out: out, Web: webFS, VideoBaseURL: "https://v.example.com"}); err == nil {
		t.Fatal("export into a non-empty folder must fail")
	}
}

func TestOutputIsWorldReadable(t *testing.T) {
	out, _, err := runExport(t, buildContent(t, nil), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{out, VideosDir(out)} {
		_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			info, _ := d.Info()
			if info.Mode().Perm()&0o044 != 0o044 {
				t.Errorf("%s is not readable by others: %v", p, info.Mode().Perm())
			}
			return nil
		})
	}
}
