package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// privateQueryParams are removed from logged query strings so that readers'
// search terms are never written to logs.
var privateQueryParams = []string{"q"}

// RequestLogger logs one record per request. It replaces gin.Logger().
func RequestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.LogAttrs(c.Request.Context(), slog.LevelInfo, "request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", redactQuery(c.Request.URL.RawQuery)),
			slog.Int("status", c.Writer.Status()),
			slog.Int64("durationMs", time.Since(start).Milliseconds()),
			slog.String("ip", c.ClientIP()),
		)
	}
}

func redactQuery(raw string) string {
	if raw == "" {
		return ""
	}
	values, err := url.ParseQuery(raw)
	if err != nil {
		return ""
	}
	for _, k := range privateQueryParams {
		values.Del(k)
	}
	return values.Encode()
}

// Problem is an RFC 9457 problem details body.
type Problem struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

// WriteProblem writes an application/problem+json response and aborts the chain.
func WriteProblem(c *gin.Context, status int, detail string) {
	body, _ := json.Marshal(Problem{Status: status, Title: http.StatusText(status), Detail: detail})
	c.Data(status, "application/problem+json", body)
	c.Abort()
}

// Recovery turns a handler panic into a 500 problem+json response. It replaces gin.Recovery().
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.ErrorContext(c.Request.Context(), "panic", slog.Any("panic", r), slog.String("path", c.Request.URL.Path))
				if !c.Writer.Written() {
					WriteProblem(c, http.StatusInternalServerError, "unexpected server error")
				} else {
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
