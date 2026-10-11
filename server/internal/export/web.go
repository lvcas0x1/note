package export

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ConfigFile is the runtime configuration script the web app loads first.
const ConfigFile = "config.js"

// RuntimeConfig is what config.js sets on window.__NOTES_CONFIG__.
type RuntimeConfig struct {
	Mode         string `json:"mode"` // "server" or "static"
	BasePath     string `json:"basePath"`
	VideoBaseURL string `json:"videoBaseUrl,omitempty"`
}

// ConfigScript renders config.js.
func ConfigScript(c RuntimeConfig) []byte {
	raw, _ := json.Marshal(c)
	return []byte("window.__NOTES_CONFIG__ = " + string(raw) + ";\n")
}

// StaticCSP is the policy for static pages (design D10). frame-ancestors is not
// included because browsers ignore it in a <meta> element. 'wasm-unsafe-eval'
// is not needed: SQLite WASM runs in a dedicated worker, which this document
// policy does not govern (verified by apps/web/e2e/static-search.mjs).
func StaticCSP(videoOrigin string) string {
	media := "'self'"
	if videoOrigin != "" {
		media += " " + videoOrigin
	}
	return "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data: blob:; media-src " + media + "; connect-src 'self'; font-src 'self'; " +
		"worker-src 'self' blob:; object-src 'none'; base-uri 'self'; form-action 'none'"
}

func (w *writer) writeWeb() error {
	videoOrigin := ""
	if w.o.VideoBaseURL != "" {
		u, err := url.Parse(w.o.VideoBaseURL)
		if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
			return fmt.Errorf("video base URL %q must be an absolute http(s) URL", w.o.VideoBaseURL)
		}
		videoOrigin = u.Scheme + "://" + u.Host
	}
	err := fs.WalkDir(w.o.Web, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || p == "index.html" || strings.HasPrefix(pathBase(p), ".") {
			return err
		}
		raw, err := fs.ReadFile(w.o.Web, p)
		if err != nil {
			return err
		}
		return writeFile(filepath.Join(w.o.Out, filepath.FromSlash(p)), raw)
	})
	if err != nil {
		return err
	}
	index, err := fs.ReadFile(w.o.Web, "index.html")
	if err != nil {
		return err
	}
	shell := []byte(w.shellHTML(string(index), StaticCSP(videoOrigin)))
	cfg := RuntimeConfig{Mode: "static", BasePath: w.o.BasePath, VideoBaseURL: w.o.VideoBaseURL}
	if err := writeFile(filepath.Join(w.o.Out, ConfigFile), ConfigScript(cfg)); err != nil {
		return err
	}
	shells := []string{"index.html", "404.html"}
	all, err := w.r.AllArticles(w.ctx)
	if err != nil {
		return err
	}
	for _, a := range all {
		shells = append(shells, filepath.Join("a", filepath.FromSlash(a.ID), "index.html"))
	}
	days, err := w.r.Calendar(w.ctx, "", "")
	if err != nil {
		return err
	}
	for _, d := range days {
		shells = append(shells, filepath.Join("d", d.Date, "index.html"))
	}
	for _, s := range shells {
		if err := writeFile(filepath.Join(w.o.Out, s), shell); err != nil {
			return err
		}
	}
	return nil
}

// shellHTML prefixes root-relative URLs with the base path and adds the CSP meta tag.
func (w *writer) shellHTML(index, csp string) string {
	base := w.o.BasePath
	for _, attr := range []string{`src="/`, `href="/`} {
		index = strings.ReplaceAll(index, attr, attr[:len(attr)-1]+base)
	}
	meta := `<meta http-equiv="Content-Security-Policy" content="` + csp + `">`
	if i := strings.Index(strings.ToLower(index), "<head>"); i >= 0 {
		return index[:i+len("<head>")] + "\n    " + meta + index[i+len("<head>"):]
	}
	return meta + index
}

func writeFile(file string, raw []byte) error {
	if err := os.MkdirAll(filepath.Dir(file), publicDirMode); err != nil { //nolint:gosec // public website output
		return err
	}
	return os.WriteFile(file, raw, publicFileMode) //nolint:gosec // public website output
}

func pathBase(p string) string {
	if i := strings.LastIndex(p, "/"); i >= 0 {
		return p[i+1:]
	}
	return p
}
