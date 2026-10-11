package httpx

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewEngine returns a Gin engine with the project's logging and recovery,
// raw-path routing (article IDs are URL-encoded and contain %2F), and
// problem+json responses for unknown API routes and wrong methods.
// Forwarding headers (X-Forwarded-For) are trusted only from trustedProxies.
func NewEngine(log *slog.Logger, trustedProxies ...string) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	// Gin trusts every proxy by default; trust none unless configured.
	if len(trustedProxies) == 0 {
		trustedProxies = nil
	}
	if err := e.SetTrustedProxies(trustedProxies); err != nil {
		return nil, err
	}
	e.UseRawPath = true
	e.UnescapePathValues = true
	e.HandleMethodNotAllowed = true
	e.RedirectTrailingSlash = false
	e.Use(Recovery(log), RequestLogger(log))
	e.NoMethod(func(c *gin.Context) {
		WriteProblem(c, http.StatusMethodNotAllowed, "only GET and HEAD are supported")
	})
	return e, nil
}
