package httpx

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"notes/internal/store"
)

// Builds gives handlers the build that is currently served.
type Builds interface {
	Reader() (*store.Reader, error)
}

// AssetHandler serves /assets/<hash>/<name> from the current build. Only URLs
// listed in the build's asset table are served, so no other file is reachable.
// http.ServeContent handles Range and If-Range, so videos can seek.
func AssetHandler(builds Builds) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := builds.Reader()
		if err != nil {
			WriteProblem(c, http.StatusServiceUnavailable, "no content has been built yet")
			return
		}
		a, file, err := r.Asset(c.Request.Context(), c.Request.URL.Path)
		if errors.Is(err, store.ErrNotFound) {
			WriteProblem(c, http.StatusNotFound, "asset not found")
			return
		}
		if err != nil {
			WriteProblem(c, http.StatusInternalServerError, "cannot read the asset")
			return
		}
		f, err := os.Open(file)
		if err != nil {
			WriteProblem(c, http.StatusNotFound, "asset not found")
			return
		}
		defer func() { _ = f.Close() }()
		h := c.Writer.Header()
		h.Set("Content-Type", a.ContentType)
		// Asset URLs contain the content hash, so they never change.
		h.Set("Cache-Control", "public, max-age=31536000, immutable")
		h.Set("ETag", `"`+a.Hash+`"`)
		http.ServeContent(c.Writer, c.Request, a.Name, time.Time{}, f)
	}
}

// NewServer returns an http.Server for handler. It sets header and idle
// timeouts but no write timeout, so long video transfers are never cut off;
// slow-client protection for large downloads is left to the CDN or proxy.
func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}
