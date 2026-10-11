package content

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// AssetURLPrefix is where published assets are served.
const AssetURLPrefix = "/assets/"

// Asset is a published file referenced by at least one published article.
type Asset struct {
	Key         string // stable identity: content-relative path or remote URL
	URL         string // /assets/<hash>/<name>
	Name        string // base file name
	Hash        string // hex SHA-256 of the content
	Size        int64
	Kind        string // image, drawio, excalidraw, openapi, mp4
	ContentType string
	LocalPath   string // file to copy from (content file or downloaded temp file)
	users       map[string]bool
}

// Resolver resolves asset references for articles and tracks which articles use them.
type Resolver struct {
	root    string
	assets  map[string]*Asset // by Key
	fetcher *Fetcher
}

// NewResolver returns a resolver for content under root. fetcher may be nil to forbid remote assets.
func NewResolver(root string, fetcher *Fetcher) *Resolver {
	return &Resolver{root: root, assets: map[string]*Asset{}, fetcher: fetcher}
}

// Resolve returns the public URL for src as referenced by the article at articleRel.
func (r *Resolver) Resolve(articleRel, kind, src string) (string, error) {
	if kind == "image" && strings.HasPrefix(strings.ToLower(src), "data:") {
		return src, nil // inline images need no file
	}
	if isRemote(src) {
		return r.resolveRemote(articleRel, kind, src)
	}
	rel, err := localRel(articleRel, src)
	if err != nil {
		return "", err
	}
	if kind == EmbedMP4 && !strings.EqualFold(path.Ext(rel), ".mp4") {
		return "", fmt.Errorf("video %q must be an .mp4 file", rel)
	}
	abs, err := ResolveInside(r.root, rel)
	if err != nil {
		return "", err
	}
	return r.register(articleRel, rel, kind, abs)
}

func (r *Resolver) register(articleRel, key, kind, file string) (string, error) {
	if a, ok := r.assets[key]; ok {
		a.users[articleRel] = true
		return a.URL, nil
	}
	sum, size, err := hashFile(file)
	if err != nil {
		return "", err
	}
	name := path.Base(filepath.ToSlash(key))
	if isRemote(key) {
		name = remoteName(key, kind)
	}
	a := &Asset{
		Key:         key,
		Name:        name,
		Hash:        sum,
		Size:        size,
		Kind:        kind,
		ContentType: contentType(name),
		LocalPath:   file,
		users:       map[string]bool{articleRel: true},
	}
	a.URL = AssetURLPrefix + sum[:16] + "/" + name
	r.assets[key] = a
	return a.URL, nil
}

// Published returns the assets used by at least one of the given articles, sorted by URL.
func (r *Resolver) Published(articles map[string]bool) []*Asset {
	var out []*Asset
	for _, a := range r.assets {
		for user := range a.users {
			if articles[user] {
				out = append(out, a)
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].URL < out[j].URL })
	return out
}

// localRel turns an asset reference into a content-relative path. A leading "/"
// means the content root; anything else is relative to the article's folder.
func localRel(articleRel, src string) (string, error) {
	src = strings.SplitN(strings.SplitN(src, "#", 2)[0], "?", 2)[0]
	if src == "" {
		return "", fmt.Errorf("empty asset path")
	}
	if strings.Contains(src, "://") || strings.HasPrefix(src, "//") {
		return "", fmt.Errorf("unsupported URL scheme in %q (only http and https are allowed)", src)
	}
	var rel string
	if strings.HasPrefix(src, "/") {
		rel = path.Clean(strings.TrimPrefix(src, "/"))
	} else {
		rel = path.Join(path.Dir(articleRel), src)
	}
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return "", fmt.Errorf("path %q is outside the content folder", src)
	}
	return rel, nil
}

func isRemote(src string) bool {
	l := strings.ToLower(src)
	return strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://")
}

func hashFile(file string) (string, int64, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", 0, err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}

var contentTypes = map[string]string{
	".png": "image/png", ".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".gif": "image/gif",
	".webp": "image/webp", ".avif": "image/avif", ".svg": "image/svg+xml", ".ico": "image/x-icon",
	".mp4": "video/mp4", ".json": "application/json", ".yaml": "application/yaml", ".yml": "application/yaml",
	".drawio": "application/xml", ".xml": "application/xml", ".excalidraw": "application/json",
}

func contentType(name string) string {
	if ct, ok := contentTypes[strings.ToLower(path.Ext(name))]; ok {
		return ct
	}
	return "application/octet-stream"
}
