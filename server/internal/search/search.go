// Package search runs keyword full-text search over the FTS5 index.
package search

import (
	"context"
	"errors"
	"html"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"

	"notes/internal/store"
)

// Query limits shared with the web app's static search (fixtures/search/cases.json).
const (
	MinQueryRunes = 3
	MaxQueryRunes = 200
	MaxResults    = 20
	Timeout       = 2 * time.Second
)

var (
	// ErrTooShort means the trimmed query has fewer than MinQueryRunes characters.
	ErrTooShort = errors.New("query must have at least 3 characters")
	// ErrTooLong means the query has more than MaxQueryRunes characters.
	ErrTooLong = errors.New("query must have at most 200 characters")
)

// Hit is one search result.
type Hit struct {
	store.ArticleSummary
	Snippet string `json:"snippet" doc:"HTML-escaped excerpt; matched terms are wrapped in <mark>"`
}

// Snippet markers are control characters that cannot appear in article text
// after compilation; they are replaced by <mark> after HTML escaping.
const (
	markOpen  = "\x02"
	markClose = "\x03"
)

// MatchExpression turns user input into an FTS5 MATCH expression. Every term is
// a quoted string, so FTS5 operators in the input are never interpreted. Terms
// shorter than three characters are dropped, because the trigram index cannot
// match them, unless no other term remains.
func MatchExpression(q string) (string, error) {
	if utf8.RuneCountInString(q) > MaxQueryRunes {
		return "", ErrTooLong
	}
	q = strings.TrimSpace(q)
	if utf8.RuneCountInString(q) < MinQueryRunes {
		return "", ErrTooShort
	}
	var terms []string
	for _, t := range strings.Fields(q) {
		if utf8.RuneCountInString(t) >= MinQueryRunes {
			terms = append(terms, `"`+strings.ReplaceAll(t, `"`, `""`)+`"`)
		}
	}
	if len(terms) == 0 {
		return "", ErrTooShort
	}
	return strings.Join(terms, " "), nil
}

// Search returns up to limit articles matching q, best first.
func Search(ctx context.Context, db *gorm.DB, q string, limit int) ([]Hit, error) {
	match, err := MatchExpression(q)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > MaxResults {
		limit = MaxResults
	}
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	var rows []struct {
		ID       string
		Title    string
		Date     string
		Summary  string
		Category string
		Snippet  string
	}
	err = db.WithContext(ctx).Raw(`
		SELECT a.id, a.title, a.date, a.summary, a.category,
		       snippet(articles_fts, -1, ?, ?, '…', 16) AS snippet
		FROM articles_fts JOIN articles a ON a.id = articles_fts.id
		WHERE articles_fts MATCH ?
		ORDER BY bm25(articles_fts, 0.0, 10.0, 3.0, 3.0, 1.0), a.id
		LIMIT ?`, markOpen, markClose, match, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	hits := make([]Hit, 0, len(rows))
	for _, r := range rows {
		cat := store.ArticleSummary{ID: r.ID, Title: r.Title, Date: r.Date, Summary: r.Summary}
		cat.Category = decodeCategory(r.Category)
		hits = append(hits, Hit{ArticleSummary: cat, Snippet: SnippetHTML(r.Snippet)})
	}
	return hits, nil
}

// SnippetHTML escapes a raw FTS5 snippet and turns its markers into <mark> tags.
func SnippetHTML(raw string) string {
	s := html.EscapeString(raw)
	s = strings.ReplaceAll(s, markOpen, "<mark>")
	return strings.ReplaceAll(s, markClose, "</mark>")
}
