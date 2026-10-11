package search

import (
	"context"
	"errors"
	"strings"
	"testing"

	"notes/internal/content"
	"notes/internal/store"
	"notes/internal/store/storetest"
)

func openBuild(t *testing.T, res *content.Result) *store.Reader {
	t.Helper()
	data := t.TempDir()
	if _, err := store.WriteBuild(data, res); err != nil {
		t.Fatal(err)
	}
	r, err := store.OpenCurrent(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = r.Close() })
	return r
}

func art(id, title, text string) *content.Article {
	return &content.Article{ID: id, Path: id + ".md", Title: title, Date: "2026-10-05", Summary: "", HTML: "<p></p>", Text: title + " " + text}
}

func TestTitleMatchRanksFirst(t *testing.T) {
	r := openBuild(t, &content.Result{Articles: []*content.Article{
		art("b", "Other topic", "We tried kafka once, kafka twice, and kafka again in the body."),
		art("a", "Kafka deploy", "Notes about brokers."),
	}})
	hits, err := Search(context.Background(), r.DB(), "kafka", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 2 || hits[0].ID != "a" || hits[1].ID != "b" {
		t.Fatalf("hits = %+v, want a before b", hits)
	}
}

func TestJapaneseSubstringAndHighlight(t *testing.T) {
	r := openBuild(t, &content.Result{Articles: []*content.Article{art("jp", "メモ", "これは日本語のテキストです")}})
	hits, err := Search(context.Background(), r.DB(), "本語の", 0)
	if err != nil || len(hits) != 1 {
		t.Fatalf("hits = %+v err = %v", hits, err)
	}
	if !strings.Contains(hits[0].Snippet, "<mark>本語の</mark>") {
		t.Errorf("snippet = %q", hits[0].Snippet)
	}
}

func TestCaseInsensitive(t *testing.T) {
	r := openBuild(t, &content.Result{Articles: []*content.Article{art("k", "KAFKA", "x")}})
	if hits, err := Search(context.Background(), r.DB(), "kafka", 0); err != nil || len(hits) != 1 {
		t.Fatalf("hits = %+v err = %v", hits, err)
	}
}

func TestOperatorCharactersAreSafe(t *testing.T) {
	r := openBuild(t, &content.Result{Articles: []*content.Article{art("x", "foo bar", "OR (bar")}})
	for _, q := range []string{`"foo* OR (bar`, `NEAR(foo bar)`, `title:foo`, `foo AND NOT bar`, `"""`, `^foo`} {
		if _, err := Search(context.Background(), r.DB(), q, 0); err != nil && !errors.Is(err, ErrTooShort) {
			t.Errorf("query %q: %v", q, err)
		}
	}
}

func TestSnippetIsEscaped(t *testing.T) {
	r := openBuild(t, &content.Result{Articles: []*content.Article{art("x", "t", "use <script>alert(1)</script> carefully")}})
	hits, err := Search(context.Background(), r.DB(), "alert", 0)
	if err != nil || len(hits) != 1 {
		t.Fatalf("hits = %+v err = %v", hits, err)
	}
	if strings.Contains(hits[0].Snippet, "<script>") || !strings.Contains(hits[0].Snippet, "&lt;script&gt;") {
		t.Errorf("snippet not escaped: %q", hits[0].Snippet)
	}
}

func TestLengthRulesAndCap(t *testing.T) {
	r := openBuild(t, storetest.Result(storetest.Articles(2, 60)))
	ctx := context.Background()
	if _, err := Search(ctx, r.DB(), "  ab ", 0); !errors.Is(err, ErrTooShort) {
		t.Errorf("err = %v, want ErrTooShort", err)
	}
	if _, err := Search(ctx, r.DB(), strings.Repeat("a", 201), 0); !errors.Is(err, ErrTooLong) {
		t.Errorf("err = %v, want ErrTooLong", err)
	}
	if _, err := Search(ctx, r.DB(), strings.Repeat("あ", 200), 0); err != nil {
		t.Errorf("200 characters must be allowed: %v", err)
	}
	hits, err := Search(ctx, r.DB(), "kafka", 100)
	if err != nil || len(hits) != MaxResults {
		t.Fatalf("got %d hits err = %v, want the cap of %d", len(hits), err, MaxResults)
	}
}

func TestDraftsNeverMatch(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := writeFile(root, name, body); err != nil {
			t.Fatal(err)
		}
	}
	write("sidebar.json", `[{"title":"P","path":"p.md"}]`)
	write("p.md", "---\ndate: 2026-01-01\n---\npublished words\n")
	write("d.md", "---\ndate: 2026-01-01\ndraft: true\n---\n# zebracorn secret\n")
	res, err := content.Build(root, content.Options{})
	if err != nil {
		t.Fatal(err)
	}
	r := openBuild(t, res)
	hits, err := Search(context.Background(), r.DB(), "zebracorn", 0)
	if err != nil || len(hits) != 0 {
		t.Fatalf("draft matched: %+v err = %v", hits, err)
	}
}
