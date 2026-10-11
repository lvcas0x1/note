//go:build windows

package content

import (
	"io/fs"
	"syscall"
	"time"
)

// birthTime returns the file's creation time where the OS records it.
func birthTime(info fs.FileInfo) (time.Time, bool) {
	d, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return time.Time{}, false
	}
	return time.Unix(0, d.CreationTime.Nanoseconds()), true
}
