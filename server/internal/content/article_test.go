package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func compileOne(t *testing.T, files map[string]string, rel, sidebarTitle string) (*Article, error) {
	t.Helper()
	root := t.TempDir()
	writeFiles(t, root, files)
	return CompileArticle(root, rel, sidebarTitle)
}

func TestFrontMatterValues(t *testing.T) {
	a, err := compileOne(t, map[string]string{"a.md": "---\ntitle: Front Title\ndate: 2026-10-11\nsummary: Short.\ntags: [go, sqlite]\n---\n# H1 Title\n\nFirst paragraph.\n"}, "a.md", "Sidebar")
	if err != nil {
		t.Fatal(err)
	}
	if a.Title != "Front Title" || a.Date != "2026-10-11" || a.Summary != "Short." || strings.Join(a.Tags, ",") != "go,sqlite" || a.Draft {
		t.Fatalf("article = %+v", a)
	}
	if a.ID != "a" || len(a.Warnings) != 0 {
		t.Fatalf("id = %q warnings = %v", a.ID, a.Warnings)
	}
}

func TestQuotedDateAndDraft(t *testing.T) {
	a, err := compileOne(t, map[string]string{"x/b.md": "---\ndate: \"2025-03-14\"\ndraft: true\n---\nBody\n"}, "x/b.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if a.Date != "2025-03-14" || !a.Draft || a.ID != "x/b" {
		t.Fatalf("article = %+v", a)
	}
}

func TestDocsifyArticleFallbacks(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{"golang/basic.md": "# Hello World\n\nGo is a *small* language\nwith a `go` tool.\n\n## More\n"})
	mtime := time.Date(2024, 5, 6, 12, 0, 0, 0, time.Local)
	if err := os.Chtimes(filepath.Join(root, "golang/basic.md"), mtime, mtime); err != nil {
		t.Fatal(err)
	}
	a, err := CompileArticle(root, "golang/basic.md", "Basic")
	if err != nil {
		t.Fatal(err)
	}
	if a.Title != "Hello World" {
		t.Errorf("title = %q, want H1", a.Title)
	}
	if a.Date != "2024-05-06" {
		t.Errorf("date = %q, want mtime date", a.Date)
	}
	if a.Summary != "Go is a small language with a go tool." {
		t.Errorf("summary = %q", a.Summary)
	}
	if len(a.Warnings) != 1 || !strings.Contains(a.Warnings[0], "golang/basic.md") || !strings.Contains(a.Warnings[0], "no date") {
		t.Errorf("warnings = %v", a.Warnings)
	}
}

func TestTitleFallbackOrder(t *testing.T) {
	a, err := compileOne(t, map[string]string{"n.md": "---\ndate: 2026-01-01\n---\nNo heading here.\n"}, "n.md", "From Sidebar")
	if err != nil || a.Title != "From Sidebar" {
		t.Fatalf("title = %q err = %v, want sidebar title", a.Title, err)
	}
	a, err = compileOne(t, map[string]string{"dir/file-name.md": "---\ndate: 2026-01-01\n---\ntext\n"}, "dir/file-name.md", "")
	if err != nil || a.Title != "file-name" {
		t.Fatalf("title = %q err = %v, want file name", a.Title, err)
	}
}

func TestSummaryCutTo200Runes(t *testing.T) {
	long := strings.Repeat("日本語", 100)
	a, err := compileOne(t, map[string]string{"s.md": "---\ndate: 2026-01-01\n---\n" + long + "\n"}, "s.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len([]rune(a.Summary)); got != SummaryMaxRunes {
		t.Fatalf("summary has %d runes", got)
	}
}

func TestInvalidDateFails(t *testing.T) {
	for _, bad := range []string{"2026-02-30", "2026-13-01", "\"10/11/2026\"", "2026-10-11T10:00:00Z"} {
		_, err := compileOne(t, map[string]string{"bad.md": "---\ndate: " + bad + "\n---\nx\n"}, "bad.md", "")
		if err == nil {
			t.Errorf("date %s: expected error", bad)
			continue
		}
		if !strings.Contains(err.Error(), "bad.md") || !strings.Contains(err.Error(), "/date") {
			t.Errorf("date %s: error %q should name the file and the field", bad, err)
		}
	}
}

func TestUnknownKeyWarns(t *testing.T) {
	a, err := compileOne(t, map[string]string{"u.md": "---\ndate: 2026-01-01\nlayout: post\n---\nx\n"}, "u.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Warnings) != 1 || !strings.Contains(a.Warnings[0], `"layout"`) {
		t.Fatalf("warnings = %v", a.Warnings)
	}
}

func TestPlainTextForSearch(t *testing.T) {
	a, err := compileOne(t, map[string]string{"p.md": "---\ndate: 2026-01-01\n---\n# Title\n\nPara **bold**.\n\n```go\nfunc main() {}\n```\n"}, "p.md", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Title", "Para bold.", "func main() {}"} {
		if !strings.Contains(a.Text, want) {
			t.Errorf("text %q missing %q", a.Text, want)
		}
	}
}

func TestLeadingTitleHeadingIsDropped(t *testing.T) {
	a, err := compileOne(t, map[string]string{"t.md": "---\ndate: 2026-01-01\n---\n# Same Title\n\nBody.\n\n# Another H1\n"}, "t.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if a.Title != "Same Title" || strings.Contains(a.HTML, ">Same Title<") || !strings.Contains(a.HTML, "Another H1") {
		t.Fatalf("title = %q html = %s", a.Title, a.HTML)
	}
	b, err := compileOne(t, map[string]string{"u.md": "---\ntitle: Front\ndate: 2026-01-01\n---\n# Different\n\nBody.\n"}, "u.md", "")
	if err != nil || !strings.Contains(b.HTML, ">Different<") {
		t.Fatalf("a heading that differs from the title must stay: %s err = %v", b.HTML, err)
	}
}
