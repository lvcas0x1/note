package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStampDates(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"plain.md":       "# Plain\n\nNo front matter.\n",
		"fm.md":          "---\ntitle: Has FM\ntags: [a]\n---\n# Body\n",
		"dated.md":       "---\ndate: 2025-01-02\n---\nAlready dated.\n",
		"crlf.md":        "---\r\ntitle: CRLF\r\n---\r\nbody\r\n",
		".hidden/x.md":   "# skipped\n",
		"nested/deep.md": "Deep\n",
	})
	when := time.Date(2024, 3, 9, 12, 0, 0, 0, time.Local)
	for _, f := range []string{"plain.md", "fm.md", "crlf.md", "nested/deep.md"} {
		if err := os.Chtimes(filepath.Join(root, f), when, when); err != nil {
			t.Fatal(err)
		}
	}

	dry, err := StampDates(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(dry) != 4 {
		t.Fatalf("dry run stamps = %+v", dry)
	}
	if raw, _ := os.ReadFile(filepath.Join(root, "plain.md")); strings.Contains(string(raw), "date:") {
		t.Fatal("dry run wrote a file")
	}

	stamps, err := StampDates(root, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range stamps {
		// The creation time of a just-written file is today, so the result is
		// the creation date where the OS records one, else the modification date.
		if s.Source == "modified" && s.Date != "2024-03-09" {
			t.Errorf("%s: %+v", s.Path, s)
		}
	}
	read := func(name string) string {
		raw, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(raw)
	}
	if got := read("plain.md"); !strings.HasPrefix(got, "---\ndate: ") || !strings.Contains(got, "\n---\n# Plain\n") {
		t.Errorf("plain.md = %q", got)
	}
	if got := read("fm.md"); !strings.HasPrefix(got, "---\ndate: ") || !strings.Contains(got, "title: Has FM\ntags: [a]\n---\n# Body\n") {
		t.Errorf("fm.md = %q", got)
	}
	if got := read("dated.md"); got != "---\ndate: 2025-01-02\n---\nAlready dated.\n" {
		t.Errorf("dated.md changed: %q", got)
	}
	if got := read("crlf.md"); !strings.HasPrefix(got, "---\r\ndate: ") || strings.Contains(strings.ReplaceAll(got, "\r\n", ""), "\n") {
		t.Errorf("crlf.md line endings broken: %q", got)
	}
	if got := read(".hidden/x.md"); strings.Contains(got, "date:") {
		t.Error("hidden file was stamped")
	}

	// Every stamped article now builds without a missing-date warning.
	for _, f := range []string{"plain.md", "fm.md", "crlf.md", "nested/deep.md"} {
		a, err := CompileArticle(root, f, "")
		if err != nil {
			t.Fatal(err)
		}
		for _, w := range a.Warnings {
			if strings.Contains(w, "no date") {
				t.Errorf("%s still warns: %s", f, w)
			}
		}
	}
	again, err := StampDates(root, false)
	if err != nil || len(again) != 0 {
		t.Fatalf("second run stamped %+v err = %v", again, err)
	}
}

func TestStampPrefersCreationTime(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{"a.md": "x\n"})
	info, err := os.Stat(filepath.Join(root, "a.md"))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := birthTime(info); !ok {
		t.Skip("this OS does not expose creation times")
	}
	// Push the modification time into the future: the creation date must win.
	future := time.Now().Add(72 * time.Hour)
	if err := os.Chtimes(filepath.Join(root, "a.md"), future, future); err != nil {
		t.Fatal(err)
	}
	stamps, err := StampDates(root, true)
	if err != nil || len(stamps) != 1 {
		t.Fatalf("stamps = %+v err = %v", stamps, err)
	}
	if stamps[0].Source != "created" || stamps[0].Date != time.Now().Format(DateLayout) {
		t.Fatalf("stamp = %+v, want today's creation date", stamps[0])
	}
}
