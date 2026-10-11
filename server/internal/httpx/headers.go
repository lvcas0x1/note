package httpx

import "github.com/gin-gonic/gin"

// ContentSecurityPolicy allows only the site's own origin (design D2).
// Inline styles are allowed because Mermaid, ECharts, and swagger-ui inject them.
const ContentSecurityPolicy = "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; " +
	"img-src 'self' data: blob:; media-src 'self'; connect-src 'self'; font-src 'self'; worker-src 'self' blob:; " +
	"object-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'none'"

// SecurityHeaders sets the security headers on every response.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("Content-Security-Policy", ContentSecurityPolicy)
		h.Set("Strict-Transport-Security", "max-age=31536000")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
