package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"notes/internal/content"
	"notes/internal/store"
)

// sharedDir holds the search cases shared with the web app's static search.
var sharedDir = filepath.Join("..", "..", "..", "fixtures", "search")

type fixtureFile struct {
	Articles []struct {
		ID, Title, Date, Summary, Text string
		Tags                           []string
	} `json:"articles"`
}

type caseFile struct {
	Comment string      `json:"comment"`
	Cases   []caseEntry `json:"cases"`
}

type caseEntry struct {
	Name     string            `json:"name"`
	Query    string            `json:"query"`
	IDs      []string          `json:"ids,omitempty"` // missing means no results
	Error    string            `json:"error,omitempty"`
	Snippets map[string]string `json:"snippets,omitempty"`
}

// FixtureResult loads fixtures/search/fixture.json as a build result.
func FixtureResult(t testing.TB) *content.Result {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(sharedDir, "fixture.json"))
	if err != nil {
		t.Fatal(err)
	}
	var f fixtureFile
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	res := &content.Result{}
	for _, a := range f.Articles {
		res.Articles = append(res.Articles, &content.Article{ID: a.ID, Path: a.ID + ".md", Title: a.Title, Date: a.Date,
			Summary: a.Summary, Tags: a.Tags, Text: a.Text, HTML: "<p></p>"})
	}
	return res
}

func TestSharedSearchCases(t *testing.T) {
	path := filepath.Join(sharedDir, "cases.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var cases caseFile
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	r := openBuild(t, FixtureResult(t))
	changed := false
	for i, c := range cases.Cases {
		t.Run(c.Name, func(t *testing.T) {
			hits, err := Search(context.Background(), r.DB(), c.Query, 0)
			switch c.Error {
			case "too_short":
				if !errors.Is(err, ErrTooShort) {
					t.Fatalf("err = %v, want too_short", err)
				}
				return
			case "too_long":
				if !errors.Is(err, ErrTooLong) {
					t.Fatalf("err = %v, want too_long", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			ids := []string{}
			snippets := map[string]string{}
			for _, h := range hits {
				ids = append(ids, h.ID)
				snippets[h.ID] = h.Snippet
			}
			want := c.IDs
			if want == nil {
				want = []string{}
			}
			if !slices.Equal(ids, want) {
				t.Fatalf("ids = %v, want %v", ids, want)
			}
			if *updateCases {
				cases.Cases[i].Snippets = snippets
				changed = true
				return
			}
			for id, want := range c.Snippets {
				if snippets[id] != want {
					t.Errorf("snippet for %s = %q, want %q", id, snippets[id], want)
				}
			}
		})
	}
	if changed {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(cases); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

var updateCases = flag.Bool("update-cases", false, "record snippets in fixtures/search/cases.json")

var updateSearchDB = flag.Bool("update-searchdb", false, "rewrite fixtures/search/fixture.search.db")

// TestSharedCasesOnStrippedDB writes (with -update-searchdb) and checks the
// search-only database that the static site downloads. The web app's Vitest
// suite runs the same cases against this file with SQLite WASM.
func TestSharedCasesOnStrippedDB(t *testing.T) {
	file := filepath.Join(sharedDir, "fixture.search.db")
	if *updateSearchDB {
		data := t.TempDir()
		dir, err := store.WriteBuild(data, FixtureResult(t))
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(filepath.Join(dir, store.DBFile))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, raw, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := store.StripForSearch(file); err != nil {
			t.Fatal(err)
		}
	}
	db, err := store.OpenReadOnly(file)
	if err != nil {
		t.Fatalf("open %s (run with -update-searchdb): %v", file, err)
	}
	defer func() { _ = store.Close(db) }()
	hits, err := Search(context.Background(), db, "kafka", 0)
	if err != nil || len(hits) != 2 || hits[0].ID != "kafka/deploy" {
		t.Fatalf("hits = %+v err = %v", hits, err)
	}
}
