package store

import (
	"path/filepath"
	"testing"
)

func TestCreateSchema(t *testing.T) {
	file := filepath.Join(t.TempDir(), "content.db")
	db, err := Create(file)
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	if err := db.Raw(`SELECT name FROM sqlite_master WHERE type IN ('table') ORDER BY name`).Scan(&tables).Error; err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"articles": false, "sidebar_docs": false, "assets": false, "articles_fts": false}
	for _, name := range tables {
		if _, ok := want[name]; ok {
			want[name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("table %s missing; have %v", name, tables)
		}
	}
	// The FTS table uses the trigram tokenizer: a 3-character CJK substring matches.
	if err := db.Exec(`INSERT INTO articles_fts (id, title, summary, tags, body) VALUES ('a', 'T', '', '', '日本語のテキスト')`).Error; err != nil {
		t.Fatal(err)
	}
	var n int64
	if err := db.Raw(`SELECT count(*) FROM articles_fts WHERE articles_fts MATCH ?`, `"本語の"`).Scan(&n).Error; err != nil || n != 1 {
		t.Fatalf("trigram match = %d err = %v", n, err)
	}
	if err := Close(db); err != nil {
		t.Fatal(err)
	}

	ro, err := OpenReadOnly(file)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = Close(ro) }()
	if err := ro.Exec(`INSERT INTO articles_fts (id, title, summary, tags, body) VALUES ('b','','','','x')`).Error; err == nil {
		t.Fatal("write through the read-only connection succeeded")
	}
}
