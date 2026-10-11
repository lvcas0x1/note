package logx

import (
	"io"
	"log/slog"
	"time"
)

// Pino numeric levels (https://getpino.io): debug 20, info 30, warn 40, error 50.
const (
	pinoDebug = 20
	pinoInfo  = 30
	pinoWarn  = 40
	pinoError = 50
)

// New returns a JSON logger whose records use Pino-compatible keys:
// numeric "level", "time" in epoch milliseconds, and "msg".
func New(w io.Writer, level slog.Leveler) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: replaceAttr,
	}))
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}
	switch a.Key {
	case slog.TimeKey:
		if t, ok := a.Value.Any().(time.Time); ok {
			return slog.Int64(slog.TimeKey, t.UnixMilli())
		}
	case slog.LevelKey:
		if l, ok := a.Value.Any().(slog.Level); ok {
			return slog.Int(slog.LevelKey, pinoLevel(l))
		}
	}
	return a
}

func pinoLevel(l slog.Level) int {
	switch {
	case l >= slog.LevelError:
		return pinoError
	case l >= slog.LevelWarn:
		return pinoWarn
	case l >= slog.LevelInfo:
		return pinoInfo
	default:
		return pinoDebug
	}
}
