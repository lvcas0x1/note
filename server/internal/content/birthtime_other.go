//go:build !darwin && !windows

package content

import (
	"io/fs"
	"time"
)

// birthTime reports false: the standard library does not expose creation
// times on this OS, so callers fall back to the modification time.
func birthTime(fs.FileInfo) (time.Time, bool) { return time.Time{}, false }
