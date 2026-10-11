package store

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	"gorm.io/gorm"

	"notes/internal/content"
)

// ErrNotFound is returned when an article or asset does not exist in the build.
var ErrNotFound = errors.New("not found")

// ArticleSummary is an article as listed in day lists and search results.
type ArticleSummary struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Date     string   `json:"date" format:"date"`
	Summary  string   `json:"summary"`
	Category []string `json:"category" doc:"Titles of the sidebar categories that contain the article"`
}

// ArticleDetail is a full article.
type ArticleDetail struct {
	ArticleSummary
	Tags      []string `json:"tags"`
	InSidebar bool     `json:"inSidebar"`
	HTML      string   `json:"html" doc:"Compiled, sanitized HTML body"`
}

// DayCount is the number of articles on one date.
type DayCount struct {
	Date  string `json:"date" format:"date"`
	Count int    `json:"count" minimum:"1"`
}

// Reader runs queries against one build.
type Reader struct {
	db  *gorm.DB
	dir string
}

// OpenBuild opens the build in dir read-only.
func OpenBuild(dir string) (*Reader, error) {
	db, err := OpenReadOnly(filepath.Join(dir, DBFile))
	if err != nil {
		return nil, err
	}
	return &Reader{db: db, dir: dir}, nil
}

// OpenCurrent opens the build that dataDir/current points to.
func OpenCurrent(dataDir string) (*Reader, error) {
	dir, err := CurrentDir(dataDir)
	if err != nil {
		return nil, err
	}
	return OpenBuild(dir)
}

// Dir is the build folder.
func (r *Reader) Dir() string { return r.dir }

// DB exposes the read-only connection for search.
func (r *Reader) DB() *gorm.DB { return r.db }

// Close releases the connection.
func (r *Reader) Close() error { return Close(r.db) }

// Ping checks that the database is readable.
func (r *Reader) Ping() error {
	var n int64
	return r.db.Model(&SidebarDoc{}).Count(&n).Error
}

// Sidebar returns the pruned sidebar tree.
func (r *Reader) Sidebar(ctx context.Context) ([]content.Node, error) {
	var doc SidebarDoc
	if err := r.db.WithContext(ctx).First(&doc, 1).Error; err != nil {
		return nil, err
	}
	nodes := []content.Node{}
	err := json.Unmarshal([]byte(doc.JSON), &nodes)
	return nodes, err
}

// Article returns one article by ID.
func (r *Reader) Article(ctx context.Context, id string) (*ArticleDetail, error) {
	var row Article
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return detail(row), nil
}

// AllArticles returns every article with its HTML, sorted by ID.
func (r *Reader) AllArticles(ctx context.Context) ([]*ArticleDetail, error) {
	var rows []Article
	if err := r.db.WithContext(ctx).Order("id").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*ArticleDetail, 0, len(rows))
	for _, row := range rows {
		out = append(out, detail(row))
	}
	return out, nil
}

// Calendar returns per-day counts for days in [from, to] that have articles.
func (r *Reader) Calendar(ctx context.Context, from, to string) ([]DayCount, error) {
	out := []DayCount{}
	q := r.db.WithContext(ctx).Model(&Article{}).Select("date, count(*) AS count").Group("date").Order("date")
	if from != "" {
		q = q.Where("date >= ? AND date <= ?", from, to)
	}
	err := q.Scan(&out).Error
	return out, err
}

// Day returns the articles dated date, sorted by title.
func (r *Reader) Day(ctx context.Context, date string) ([]ArticleSummary, error) {
	var rows []Article
	if err := r.db.WithContext(ctx).Where("date = ?", date).Order("title COLLATE NOCASE, id").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]ArticleSummary, 0, len(rows))
	for _, row := range rows {
		out = append(out, summary(row))
	}
	return out, nil
}

// Asset looks up a published asset by URL and returns it with its file path.
func (r *Reader) Asset(ctx context.Context, url string) (*Asset, string, error) {
	var a Asset
	err := r.db.WithContext(ctx).Where("url = ?", url).Take(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", ErrNotFound
	}
	if err != nil {
		return nil, "", err
	}
	return &a, filepath.Join(r.dir, AssetsDir, filepath.FromSlash(a.File)), nil
}

// Assets returns every published asset.
func (r *Reader) Assets(ctx context.Context) ([]Asset, error) {
	var out []Asset
	err := r.db.WithContext(ctx).Order("url").Find(&out).Error
	return out, err
}

func summary(row Article) ArticleSummary {
	return ArticleSummary{ID: row.ID, Title: row.Title, Date: row.Date, Summary: row.Summary, Category: decodeList(row.Category)}
}

func detail(row Article) *ArticleDetail {
	return &ArticleDetail{ArticleSummary: summary(row), Tags: decodeList(row.Tags), InSidebar: row.InSidebar, HTML: row.HTML}
}

func decodeList(s string) []string {
	out := []string{}
	if strings.TrimSpace(s) != "" {
		_ = json.Unmarshal([]byte(s), &out)
	}
	return out
}
