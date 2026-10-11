package httpx

import (
	"log/slog"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"

	"notes/internal/store"
)

// Current serves the build that dataDir/current points to and switches to a
// new build when the link changes, without a restart.
type Current struct {
	dataDir string
	log     *slog.Logger
	reader  atomic.Pointer[store.Reader]
	mu      sync.Mutex // serializes reloads
	grace   time.Duration
	watcher *fsnotify.Watcher
	done    chan struct{}
}

// NewCurrent opens the current build if there is one and starts watching dataDir.
func NewCurrent(dataDir string, log *slog.Logger) (*Current, error) {
	c := &Current{dataDir: dataDir, log: log, grace: 30 * time.Second, done: make(chan struct{})}
	c.Reload()
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := w.Add(dataDir); err != nil {
		_ = w.Close()
		return nil, err
	}
	c.watcher = w
	go c.watch()
	return c, nil
}

// Reader returns the served build or store.ErrNoBuild.
func (c *Current) Reader() (*store.Reader, error) {
	if r := c.reader.Load(); r != nil {
		return r, nil
	}
	return nil, store.ErrNoBuild
}

// Reload opens the build that current points to, if it changed.
func (c *Current) Reload() {
	c.mu.Lock()
	defer c.mu.Unlock()
	dir, err := store.CurrentDir(c.dataDir)
	if err != nil {
		return
	}
	if old := c.reader.Load(); old != nil && old.Dir() == dir {
		return
	}
	r, err := store.OpenBuild(dir)
	if err != nil {
		c.log.Error("open build", "dir", dir, "err", err)
		return
	}
	if err := r.Ping(); err != nil {
		c.log.Error("open build", "dir", dir, "err", err)
		_ = r.Close()
		return
	}
	old := c.reader.Swap(r)
	c.log.Info("serving build", "build", filepath.Base(dir))
	if old != nil {
		// Requests may still use the old build; close it after a grace period.
		time.AfterFunc(c.grace, func() { _ = old.Close() })
	}
}

// watch reloads on any change in the data folder. Renaming a new symlink over
// current does not always produce an event for "current" itself (kqueue on
// macOS reports only the temporary name), so every event triggers a cheap
// check, and a slow poll covers missed events.
func (c *Current) watch() {
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for {
		select {
		case _, ok := <-c.watcher.Events:
			if !ok {
				return
			}
			c.Reload()
		case <-tick.C:
			c.Reload()
		case <-c.done:
			return
		case err, ok := <-c.watcher.Errors:
			if !ok {
				return
			}
			c.log.Error("watch data folder", "err", err)
		}
	}
}

// Close stops watching and closes the served build.
func (c *Current) Close() error {
	close(c.done)
	if c.watcher != nil {
		_ = c.watcher.Close()
	}
	if r := c.reader.Swap(nil); r != nil {
		return r.Close()
	}
	return nil
}
