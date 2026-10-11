//go:build darwin

package content

import (
	"io/fs"
	"syscall"
	"time"
)

// birthTime returns the file's creation time where the OS records it.
func birthTime(info fs.FileInfo) (time.Time, bool) {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, false
	}
	return time.Unix(st.Birthtimespec.Sec, st.Birthtimespec.Nsec), true
}
