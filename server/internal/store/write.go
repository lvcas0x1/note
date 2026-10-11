package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"notes/internal/content"
)

// Layout of the data folder.
const (
	BuildsDir   = "builds"
	CurrentLink = "current"
	DBFile      = "content.db"
	AssetsDir   = "assets"
	KeepBuilds  = 1
)

// WriteBuild writes res into a new folder under dataDir/builds and then points
// dataDir/current at it with an atomic rename. On any error the partial build
// is removed and current is left unchanged. It returns the new build folder.
func WriteBuild(dataDir string, res *content.Result) (dir string, err error) {
	builds := filepath.Join(dataDir, BuildsDir)
	if err := os.MkdirAll(builds, 0o750); err != nil {
		return "", err
	}
	name := time.Now().UTC().Format("20060102T150405.000000000Z")
	buildDir := filepath.Join(builds, name)
	if err := os.Mkdir(buildDir, 0o750); err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(buildDir)
		}
	}()
	if err := writeDB(filepath.Join(buildDir, DBFile), res); err != nil {
		return "", err
	}
	if err := copyAssets(filepath.Join(buildDir, AssetsDir), res.Assets); err != nil {
		return "", err
	}
	if err := swapCurrent(dataDir, filepath.Join(BuildsDir, name)); err != nil {
		return "", err
	}
	pruneBuilds(builds, name)
	return buildDir, nil
}

// CurrentDir returns the folder that dataDir/current points to.
func CurrentDir(dataDir string) (string, error) {
	target, err := os.Readlink(filepath.Join(dataDir, CurrentLink))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrNoBuild
		}
		return "", err
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(dataDir, target)
	}
	return target, nil
}

// ErrNoBuild means no successful build exists yet.
var ErrNoBuild = errors.New("no successful build yet; run `notes build` first")

func writeDB(file string, res *content.Result) error {
	db, err := Create(file)
	if err != nil {
		return err
	}
	defer func() { _ = Close(db) }()
	categories := map[string][]string{}
	sidebarCategories(res.Sidebar, nil, categories)
	sidebarJSON, err := json.Marshal(res.Sidebar)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&SidebarDoc{ID: 1, JSON: string(sidebarJSON)}).Error; err != nil {
			return err
		}
		for _, a := range res.Articles {
			tags := a.Tags
			if tags == nil {
				tags = []string{}
			}
			tagsJSON, _ := json.Marshal(tags)
			cat, inSidebar := categories[a.Path]
			if cat == nil {
				cat = []string{}
			}
			catJSON, _ := json.Marshal(cat)
			row := Article{ID: a.ID, Path: a.Path, Title: a.Title, Date: a.Date, Summary: a.Summary,
				Tags: string(tagsJSON), Category: string(catJSON), InSidebar: inSidebar, HTML: a.HTML}
			if err := tx.Create(&row).Error; err != nil {
				return fmt.Errorf("article %s: %w", a.ID, err)
			}
			if err := tx.Exec(`INSERT INTO articles_fts (id, title, summary, tags, body) VALUES (?, ?, ?, ?, ?)`,
				a.ID, a.Title, a.Summary, strings.Join(a.Tags, " "), a.Text).Error; err != nil {
				return fmt.Errorf("index %s: %w", a.ID, err)
			}
		}
		for _, as := range res.Assets {
			row := Asset{URL: as.URL, Name: as.Name, Hash: as.Hash, Size: as.Size, Kind: as.Kind,
				ContentType: as.ContentType, File: assetFile(as)}
			if err := tx.Create(&row).Error; err != nil {
				return fmt.Errorf("asset %s: %w", as.URL, err)
			}
		}
		return nil
	})
}

// sidebarCategories maps each article path to the titles of its ancestor nodes.
func sidebarCategories(nodes []content.Node, parents []string, out map[string][]string) {
	for _, n := range nodes {
		if n.Path != "" {
			out[n.Path] = append([]string{}, parents...)
		}
		sidebarCategories(n.Children, append(append([]string{}, parents...), n.Title), out)
	}
}

func assetFile(a *content.Asset) string {
	return filepath.ToSlash(strings.TrimPrefix(a.URL, content.AssetURLPrefix))
}

func copyAssets(dir string, assets []*content.Asset) error {
	for _, a := range assets {
		dst := filepath.Join(dir, filepath.FromSlash(assetFile(a)))
		if err := os.MkdirAll(filepath.Dir(dst), 0o750); err != nil {
			return err
		}
		if err := copyFile(a.LocalPath, dst); err != nil {
			return fmt.Errorf("asset %s: %w", a.Key, err)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

// swapCurrent atomically replaces dataDir/current with a symlink to target.
func swapCurrent(dataDir, target string) error {
	tmp := filepath.Join(dataDir, fmt.Sprintf(".%s.%d.tmp", CurrentLink, time.Now().UnixNano()))
	if err := os.Symlink(target, tmp); err != nil {
		return err
	}
	if err := os.Rename(tmp, filepath.Join(dataDir, CurrentLink)); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// pruneBuilds removes every build except the newest KeepBuilds (the current one).
// The server keeps an old build's database open for a short grace period;
// removing its folder does not break requests in flight on Unix, and asset
// URLs are content-hashed, so the new build serves the same files.
func pruneBuilds(builds, current string) {
	entries, err := os.ReadDir(builds)
	if err != nil {
		return
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	for i, n := range names {
		if i >= KeepBuilds && n != current {
			_ = os.RemoveAll(filepath.Join(builds, n))
		}
	}
}
