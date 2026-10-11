package store

import (
	"context"
	"errors"
	"sort"
	"testing"

	"notes/internal/content"
	"notes/internal/store/storetest"
)

func openFixture(t *testing.T, res *content.Result) *Reader {
	t.Helper()
	data := t.TempDir()
	if _, err := WriteBuild(data, res); err != nil {
		t.Fatal(err)
	}
	r, err := OpenCurrent(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = r.Close() })
	return r
}

func TestQueries(t *testing.T) {
	ctx := context.Background()
	articles := storetest.Articles(1, 50)
	r := openFixture(t, storetest.Result(articles))

	nodes, err := r.Sidebar(ctx)
	if err != nil || len(nodes) != 1 || len(nodes[0].Children[0].Children) != 3 {
		t.Fatalf("sidebar = %+v err = %v", nodes, err)
	}

	a, err := r.Article(ctx, articles[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if a.Title != articles[0].Title || a.HTML != articles[0].HTML || !a.InSidebar || len(a.Category) != 2 || a.Category[1] != "Nested" {
		t.Fatalf("article = %+v", a)
	}
	if other, _ := r.Article(ctx, articles[10].ID); other.InSidebar || len(other.Category) != 0 {
		t.Fatalf("article outside sidebar = %+v", other)
	}
	if _, err := r.Article(ctx, "missing/id"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}

	want := map[string]int{}
	for _, a := range articles {
		if a.Date >= "2026-10-01" && a.Date <= "2026-10-31" {
			want[a.Date]++
		}
	}
	counts, err := r.Calendar(ctx, "2026-10-01", "2026-10-31")
	if err != nil {
		t.Fatal(err)
	}
	if len(counts) != len(want) {
		t.Fatalf("counts = %+v, want %v", counts, want)
	}
	for _, c := range counts {
		if want[c.Date] != c.Count {
			t.Errorf("%s: count %d, want %d", c.Date, c.Count, want[c.Date])
		}
	}

	day, err := r.Day(ctx, "2026-10-05")
	if err != nil || len(day) == 0 {
		t.Fatalf("day = %v err = %v", day, err)
	}
	if !sort.SliceIsSorted(day, func(i, j int) bool { return lower(day[i].Title) < lower(day[j].Title) }) {
		t.Errorf("day list not sorted by title: %+v", day)
	}
	if empty, err := r.Day(ctx, "2026-10-06"); err != nil || len(empty) != 0 {
		t.Errorf("empty day = %v err = %v", empty, err)
	}
}

func TestAssetLookupAndReadOnly(t *testing.T) {
	ctx := context.Background()
	r := openFixture(t, sampleResult(t))
	a, file, err := r.Asset(ctx, "/assets/0123456789abcdef/a.png")
	if err != nil || a.ContentType != "image/png" || file == "" {
		t.Fatalf("asset = %+v file = %s err = %v", a, file, err)
	}
	if _, _, err := r.Asset(ctx, "/assets/0000000000000000/none.png"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if err := r.DB().Exec(`DELETE FROM articles`).Error; err == nil {
		t.Fatal("write through the served connection succeeded")
	}
}

func lower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
