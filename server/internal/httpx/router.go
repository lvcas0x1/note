package httpx

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// Mount adds the routes that are not part of the API: assets, health, and the
// web app with history fallback. Only GET and HEAD are routed.
func Mount(e *gin.Engine, builds Builds, web fs.FS) {
	e.GET("/config.js", ServerConfig)
	e.HEAD("/config.js", ServerConfig)
	e.GET("/assets/*path", AssetHandler(builds))
	e.HEAD("/assets/*path", AssetHandler(builds))
	e.GET("/healthz", Health(builds))
	e.HEAD("/healthz", Health(builds))
	spa := SPA(web)
	e.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			WriteProblem(c, http.StatusMethodNotAllowed, "only GET and HEAD are supported")
			return
		}
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api/") || strings.HasPrefix(p, "/assets/") {
			WriteProblem(c, http.StatusNotFound, "not found")
			return
		}
		spa(c)
	})
}

// Health reports 200 when the current build is readable and 503 otherwise.
func Health(builds Builds) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		r, err := builds.Reader()
		if err == nil {
			err = r.Ping()
		}
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "detail": "no readable build"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// SPA serves files from web and falls back to index.html for app routes.
func SPA(web fs.FS) gin.HandlerFunc {
	files := http.FS(web)
	return func(c *gin.Context) {
		p := path.Clean("/" + c.Request.URL.Path)
		if p != "/" {
			if f, err := web.Open(strings.TrimPrefix(p, "/")); err == nil {
				info, statErr := f.Stat()
				_ = f.Close()
				if statErr == nil && !info.IsDir() {
					if strings.HasPrefix(p, "/static/") || strings.HasPrefix(p, "/vendor/") {
						c.Header("Cache-Control", "public, max-age=31536000, immutable")
					}
					http.FileServer(files).ServeHTTP(c.Writer, c.Request)
					return
				}
			}
		}
		index, err := fs.ReadFile(web, "index.html")
		if err != nil {
			WriteProblem(c, http.StatusServiceUnavailable, "the web app has not been built")
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	}
}

// serverConfigScript tells the web app it talks to this server's API.
var serverConfigScript = []byte(`window.__NOTES_CONFIG__ = {"mode":"server","basePath":"/"};` + "\n")

// ServerConfig serves config.js for the server version.
func ServerConfig(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/javascript; charset=utf-8", serverConfigScript)
}
