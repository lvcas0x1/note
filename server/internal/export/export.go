// Package export writes a static copy of the current build for static hosts
// such as GitHub Pages.
package export

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"notes/internal/store"
)

// Static host limits (GitHub Pages: 1 GB per site; GitHub: 100 MiB per file).
const (
	MaxSiteBytes = 1 << 30
	MaxFileBytes = 100 << 20
)

// Options configures an export.
type Options struct {
	DataDir      string // folder holding builds/ and current
	Out          string // site output folder; must be empty or missing
	BasePath     string // URL path the site is served under, e.g. "/" or "/notes/"
	VideoBaseURL string // absolute URL where the videos folder is uploaded
	Web          fs.FS  // built web app (index.html, static/, vendor/)

	// Limits default to MaxSiteBytes and MaxFileBytes.
	MaxSiteBytes int64
	MaxFileBytes int64
}

// VideosDir returns the folder that receives videos for out.
func VideosDir(out string) string {
	return strings.TrimRight(out, string(os.PathSeparator)) + "-videos"
}

// Manifest lists the exported videos.
type Manifest struct {
	BaseURL string          `json:"baseUrl"`
	Files   []ManifestEntry `json:"files"`
}

// ManifestEntry is one video file, relative to the videos folder.
type ManifestEntry struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

// Summary describes a finished export.
type Summary struct {
	Articles  int
	Assets    int
	Videos    int
	SiteBytes int64
}

// Run writes the static site.
func Run(ctx context.Context, o Options) (*Summary, error) {
	base, err := normalizeBase(o.BasePath)
	if err != nil {
		return nil, err
	}
	o.BasePath = base
	if o.VideoBaseURL != "" {
		o.VideoBaseURL = strings.TrimRight(o.VideoBaseURL, "/") + "/"
	}
	r, err := store.OpenCurrent(o.DataDir)
	if errors.Is(err, store.ErrNoBuild) {
		return nil, fmt.Errorf("no successful build to export; run `notes build` first")
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	if _, err := fs.Stat(o.Web, "index.html"); err != nil {
		return nil, fmt.Errorf("the web app is not built (index.html missing); run `npm run build -w @notes/web` first")
	}
	if err := ensureEmpty(o.Out); err != nil {
		return nil, err
	}
	if err := ensureEmpty(VideosDir(o.Out)); err != nil {
		return nil, err
	}

	w := &writer{o: o, r: r, ctx: ctx}
	if err := w.run(); err != nil {
		return nil, err
	}
	if o.MaxSiteBytes == 0 {
		o.MaxSiteBytes = MaxSiteBytes
	}
	if o.MaxFileBytes == 0 {
		o.MaxFileBytes = MaxFileBytes
	}
	size, err := checkLimits(o.Out, o.MaxSiteBytes, o.MaxFileBytes)
	if err != nil {
		return nil, err
	}
	return &Summary{Articles: w.articles, Assets: w.assets, Videos: len(w.manifest.Files), SiteBytes: size}, nil
}

type writer struct {
	o        Options
	r        *store.Reader
	ctx      context.Context
	videos   map[string]string // asset URL -> public video URL
	manifest Manifest
	articles int
	assets   int
}

func (w *writer) run() error {
	if err := w.copyAssets(); err != nil {
		return err
	}
	if err := w.writeData(); err != nil {
		return err
	}
	if err := w.writeSearchDB(); err != nil {
		return err
	}
	if err := w.writeWeb(); err != nil {
		return err
	}
	if len(w.manifest.Files) > 0 {
		w.manifest.BaseURL = w.o.VideoBaseURL
		if err := writeJSON(filepath.Join(VideosDir(w.o.Out), "videos-manifest.json"), w.manifest); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) copyAssets() error {
	assets, err := w.r.Assets(w.ctx)
	if err != nil {
		return err
	}
	w.videos = map[string]string{}
	for _, a := range assets {
		_, src, err := w.r.Asset(w.ctx, a.URL)
		if err != nil {
			return err
		}
		if a.Kind == "mp4" {
			if err := w.copyVideo(a, src); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(src, filepath.Join(w.o.Out, "assets", filepath.FromSlash(a.File))); err != nil {
			return err
		}
		w.assets++
	}
	return nil
}

func (w *writer) copyVideo(a store.Asset, src string) error {
	if w.o.VideoBaseURL == "" {
		users, err := w.articlesUsing(a.URL)
		if err != nil {
			return err
		}
		return fmt.Errorf("article %s embeds a video, but no video base URL was given (--video-base-url)", strings.Join(users, ", "))
	}
	dst := filepath.Join(VideosDir(w.o.Out), filepath.FromSlash(a.File))
	if err := copyFile(src, dst); err != nil {
		return err
	}
	sum, size, err := hashFile(dst)
	if err != nil {
		return err
	}
	w.manifest.Files = append(w.manifest.Files, ManifestEntry{Path: a.File, Size: size, SHA256: sum})
	w.videos[a.URL] = w.o.VideoBaseURL + a.File
	return nil
}

func (w *writer) articlesUsing(assetURL string) ([]string, error) {
	all, err := w.r.AllArticles(w.ctx)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, a := range all {
		if strings.Contains(a.HTML, `"`+assetURL+`"`) {
			ids = append(ids, a.ID)
		}
	}
	if len(ids) == 0 {
		ids = []string{"(unknown)"}
	}
	return ids, nil
}

// rewriteHTML points asset references at the base path and videos at the video host.
func (w *writer) rewriteHTML(html string) string {
	for from, to := range w.videos {
		html = strings.ReplaceAll(html, `"`+from+`"`, `"`+to+`"`)
	}
	return strings.ReplaceAll(html, `="/assets/`, `="`+w.o.BasePath+`assets/`)
}

func (w *writer) writeData() error {
	dataDir := filepath.Join(w.o.Out, "data")
	sidebar, err := w.r.Sidebar(w.ctx)
	if err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(dataDir, "sidebar.json"), sidebar); err != nil {
		return err
	}
	days, err := w.r.Calendar(w.ctx, "", "")
	if err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(dataDir, "calendar.json"), days); err != nil {
		return err
	}
	for _, d := range days {
		list, err := w.r.Day(w.ctx, d.Date)
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(dataDir, "days", d.Date+".json"), list); err != nil {
			return err
		}
	}
	all, err := w.r.AllArticles(w.ctx)
	if err != nil {
		return err
	}
	for _, a := range all {
		a.HTML = w.rewriteHTML(a.HTML)
		file, err := ArticleDataPath(dataDir, a.ID)
		if err != nil {
			return err
		}
		if err := writeJSON(file, a); err != nil {
			return err
		}
		w.articles++
	}
	return nil
}

// ArticleDataPath is where an article's JSON goes under dataDir.
func ArticleDataPath(dataDir, id string) (string, error) {
	clean := path.Clean("/" + id)
	if clean == "/" || clean != "/"+id {
		return "", fmt.Errorf("article ID %q cannot be used as a file path", id)
	}
	return filepath.Join(dataDir, "articles", filepath.FromSlash(id)+".json"), nil
}

func normalizeBase(b string) (string, error) {
	if b == "" {
		b = "/"
	}
	if !strings.HasPrefix(b, "/") || strings.Contains(b, "//") || strings.ContainsAny(b, `"'<> ?#\`) {
		return "", fmt.Errorf("base path %q must be an absolute URL path such as / or /notes/", b)
	}
	if !strings.HasSuffix(b, "/") {
		b += "/"
	}
	return b, nil
}

func ensureEmpty(dir string) error {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("output folder %s is not empty; remove it or choose another", dir)
	}
	return nil
}

func checkLimits(out string, maxSite, maxFile int64) (int64, error) {
	var total int64
	var big []string
	err := filepath.WalkDir(out, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		total += info.Size()
		if info.Size() > maxFile {
			rel, _ := filepath.Rel(out, p)
			big = append(big, fmt.Sprintf("%s (%d MiB)", filepath.ToSlash(rel), info.Size()>>20))
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	sort.Strings(big)
	if len(big) > 0 {
		return total, fmt.Errorf("files larger than %d MiB cannot be published to GitHub: %s", maxFile>>20, strings.Join(big, ", "))
	}
	if total > maxSite {
		return total, fmt.Errorf("the site output is %d bytes (%.2f GB), over the %.2f GB GitHub Pages limit",
			total, float64(total)/(1<<30), float64(maxSite)/(1<<30))
	}
	return total, nil
}

// The export is a public website: files must be readable by whatever web
// server or host serves them, so they are world-readable (0644/0755).
const (
	publicDirMode  = 0o755
	publicFileMode = 0o644
)

func writeJSON(file string, v any) error {
	if err := os.MkdirAll(filepath.Dir(file), publicDirMode); err != nil { //nolint:gosec // public website output
		return err
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(file, raw, publicFileMode) //nolint:gosec // public website output
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), publicDirMode); err != nil { //nolint:gosec // public website output
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, publicFileMode) //nolint:gosec // public website output
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func hashFile(file string) (string, int64, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", 0, err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	n, err := io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), n, err
}
