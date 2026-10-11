package content

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Remote download limits for embed files (content-build spec: remote embed download).
const (
	DefaultFetchTimeout  = 30 * time.Second
	DefaultFetchMaxBytes = 10 << 20 // 10 MB
)

// remoteKinds are the reference kinds that may point to http(s) URLs.
// Videos must always be local files.
var remoteKinds = map[string]bool{"image": true, EmbedDrawio: true, EmbedExcalidraw: true, EmbedOpenAPI: true}

// Fetcher downloads remote embed files at build time into Dir.
type Fetcher struct {
	Client   *http.Client
	Dir      string
	Timeout  time.Duration
	MaxBytes int64
	count    int
}

// NewFetcher returns a fetcher with the default limits that stores files in dir.
func NewFetcher(dir string) *Fetcher {
	return &Fetcher{Client: http.DefaultClient, Dir: dir, Timeout: DefaultFetchTimeout, MaxBytes: DefaultFetchMaxBytes}
}

// Download fetches rawURL into a new file under f.Dir and returns its path.
func (f *Fetcher) Download(ctx context.Context, rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("unsupported URL %q (only http and https are allowed)", rawURL)
	}
	ctx, cancel := context.WithTimeout(ctx, f.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := f.Client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("download timed out after %s", f.Timeout)
		}
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}
	if err := os.MkdirAll(f.Dir, 0o750); err != nil {
		return "", err
	}
	f.count++
	dst := filepath.Join(f.Dir, fmt.Sprintf("remote-%04d", f.count))
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	n, err := io.Copy(out, io.LimitReader(resp.Body, f.MaxBytes+1))
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("download timed out after %s", f.Timeout)
		}
		return "", fmt.Errorf("download failed: %w", err)
	}
	if n > f.MaxBytes {
		_ = os.Remove(dst)
		return "", fmt.Errorf("download exceeds the %d MB limit", f.MaxBytes>>20)
	}
	return dst, nil
}

func (r *Resolver) resolveRemote(articleRel, kind, src string) (string, error) {
	if kind == EmbedMP4 {
		return "", fmt.Errorf("videos must be local files, not URLs")
	}
	if !remoteKinds[kind] {
		return "", fmt.Errorf("remote %s references are not supported", kind)
	}
	if a, ok := r.assets[src]; ok {
		a.users[articleRel] = true
		return a.URL, nil
	}
	if r.fetcher == nil {
		return "", fmt.Errorf("remote reference %q is not allowed here", src)
	}
	file, err := r.fetcher.Download(context.Background(), src)
	if err != nil {
		return "", err
	}
	return r.register(articleRel, src, kind, file)
}

// remoteName picks a file name for a downloaded asset from its URL.
func remoteName(rawURL, kind string) string {
	name := "file"
	if u, err := url.Parse(rawURL); err == nil {
		if b := path.Base(u.Path); b != "." && b != "/" && b != "" {
			name = b
		}
	}
	name = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r < 0x20 {
			return '_'
		}
		return r
	}, name)
	if path.Ext(name) == "" {
		switch kind {
		case EmbedDrawio:
			name += ".drawio"
		case EmbedExcalidraw:
			name += ".excalidraw"
		case EmbedOpenAPI:
			name += ".json"
		}
	}
	return name
}
