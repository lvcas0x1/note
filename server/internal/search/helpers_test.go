package search

import (
	"os"
	"path/filepath"
	"testing"

	"notes/internal/store"
	"notes/internal/store/storetest"
)

func writeFile(root, name, body string) error {
	p := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(body), 0o600)
}

func openBuildB(b *testing.B) *store.Reader {
	b.Helper()
	data := b.TempDir()
	if _, err := store.WriteBuild(data, storetest.Result(storetest.Articles(3, 1000))); err != nil {
		b.Fatal(err)
	}
	r, err := store.OpenCurrent(data)
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { _ = r.Close() })
	return r
}
