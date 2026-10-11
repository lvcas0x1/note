package content

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

// Stamp is a date written into an article's front matter.
type Stamp struct {
	Path   string // content-relative
	Date   string // YYYY-MM-DD
	Source string // "created" or "modified"
}

// StampDates writes `date: YYYY-MM-DD` into the front matter of every
// Markdown file under root that has no date yet. The date is the file's
// creation time when the OS records it, otherwise its modification time, in
// the local time zone. Files that already have a date are not touched. With
// dryRun, nothing is written.
func StampDates(root string, dryRun bool) ([]Stamp, error) {
	files, err := markdownFiles(root)
	if err != nil {
		return nil, err
	}
	var stamps []Stamp
	for _, rel := range files {
		path := joinRel(root, rel)
		source, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		has, err := hasDate(source)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", rel, err)
		}
		if has {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		when, from := info.ModTime(), "modified"
		if bt, ok := birthTime(info); ok && !bt.IsZero() && bt.Before(when.Add(time.Second)) {
			when, from = bt, "created"
		}
		s := Stamp{Path: rel, Date: when.In(time.Local).Format(DateLayout), Source: from}
		stamps = append(stamps, s)
		if dryRun {
			continue
		}
		if err := os.WriteFile(path, withDate(source, s.Date), info.Mode().Perm()); err != nil {
			return nil, err
		}
	}
	return stamps, nil
}

// hasDate reports whether the front matter sets a date key.
func hasDate(source []byte) (bool, error) {
	p := parseMarkdown(newMarkdown(), source)
	fm := map[string]any{}
	if data := frontmatterData(p); data != nil {
		if err := data.Decode(&fm); err != nil {
			return false, fmt.Errorf("front matter: %w", err)
		}
	}
	_, ok := fm["date"]
	return ok, nil
}

// withDate adds a date line to existing YAML front matter, or adds front matter.
func withDate(source []byte, date string) []byte {
	nl := []byte("\n")
	if bytes.Contains(source, []byte("\r\n")) {
		nl = []byte("\r\n")
	}
	line := append([]byte("date: "+date), nl...)
	opening := append([]byte("---"), nl...)
	if bytes.HasPrefix(source, opening) {
		out := append([]byte{}, opening...)
		out = append(out, line...)
		return append(out, source[len(opening):]...)
	}
	out := append([]byte{}, opening...)
	out = append(out, line...)
	out = append(out, opening...)
	return append(out, source...)
}
