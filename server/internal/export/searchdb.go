package export

import (
	"path/filepath"

	"notes/internal/store"
)

// SearchDBFile is the search index the static site downloads on first search.
const SearchDBFile = "search.db"

// writeSearchDB copies the build database and strips everything search does
// not need, so the browser runs the same FTS5 index and query as the server.
func (w *writer) writeSearchDB() error {
	dst := filepath.Join(w.o.Out, "data", SearchDBFile)
	if err := copyFile(filepath.Join(w.r.Dir(), store.DBFile), dst); err != nil {
		return err
	}
	return store.StripForSearch(dst)
}
