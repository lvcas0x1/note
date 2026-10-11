package logx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"time"
)

func TestPinoCompatibleKeys(t *testing.T) {
	var buf bytes.Buffer
	before := time.Now().UnixMilli()
	New(&buf, slog.LevelDebug).Warn("disk low", "free", 3)

	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("output is not JSON: %v (%s)", err, buf.String())
	}
	if rec["level"] != float64(pinoWarn) {
		t.Errorf("level = %v, want %d", rec["level"], pinoWarn)
	}
	if rec["msg"] != "disk low" {
		t.Errorf("msg = %v", rec["msg"])
	}
	ts, ok := rec["time"].(float64)
	if !ok || int64(ts) < before || int64(ts) > time.Now().UnixMilli() {
		t.Errorf("time = %v, want epoch ms", rec["time"])
	}
	if rec["free"] != float64(3) {
		t.Errorf("free = %v", rec["free"])
	}
}

func TestPinoLevels(t *testing.T) {
	cases := map[slog.Level]int{
		slog.LevelDebug: pinoDebug,
		slog.LevelInfo:  pinoInfo,
		slog.LevelWarn:  pinoWarn,
		slog.LevelError: pinoError,
	}
	for in, want := range cases {
		if got := pinoLevel(in); got != want {
			t.Errorf("pinoLevel(%v) = %d, want %d", in, got, want)
		}
	}
}
