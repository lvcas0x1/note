package httpx

import (
	"bytes"
	"crypto/rand"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"notes/internal/content"
	"notes/internal/logx"
	"notes/internal/store"
)

type fixedBuild struct{ r *store.Reader }

func (f fixedBuild) Reader() (*store.Reader, error) { return f.r, nil }

func assetBuild(t *testing.T, size int) (*store.Reader, []byte, string) {
	t.Helper()
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "demo.mp4")
	if err := os.WriteFile(file, data, 0o600); err != nil {
		t.Fatal(err)
	}
	url := "/assets/00112233445566ff/demo.mp4"
	res := &content.Result{Assets: []*content.Asset{{Key: "v/demo.mp4", URL: url, Name: "demo.mp4", Hash: "00112233445566ff00",
		Size: int64(size), Kind: "mp4", ContentType: "video/mp4", LocalPath: file}}}
	dataDir := t.TempDir()
	if _, err := store.WriteBuild(dataDir, res); err != nil {
		t.Fatal(err)
	}
	r, err := store.OpenCurrent(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = r.Close() })
	return r, data, url
}

func startAssetServer(t *testing.T, r *store.Reader) string {
	t.Helper()
	e, _ := NewEngine(logx.New(io.Discard, slog.LevelError))
	e.GET("/assets/*path", AssetHandler(fixedBuild{r}))
	e.HEAD("/assets/*path", AssetHandler(fixedBuild{r}))
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := NewServer("", e)
	if srv.WriteTimeout != 0 {
		t.Fatal("server must not set a write timeout")
	}
	go func() { _ = srv.Serve(ln) }()
	t.Cleanup(func() { _ = srv.Close() })
	return "http://" + ln.Addr().String()
}

func TestVideoRangeRequest(t *testing.T) {
	r, data, url := assetBuild(t, 3<<20)
	base := startAssetServer(t, r)
	req, _ := http.NewRequest(http.MethodGet, base+url, nil)
	req.Header.Set("Range", "bytes=1000000-1999999")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusPartialContent {
		t.Fatalf("status = %d, want 206", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "video/mp4" || !bytes.Equal(body, data[1000000:2000000]) {
		t.Fatalf("content type %q, %d bytes", resp.Header.Get("Content-Type"), len(body))
	}
	if cc := resp.Header.Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
		t.Errorf("cache-control = %q", cc)
	}
}

func TestUnlistedAssetIs404(t *testing.T) {
	r, _, _ := assetBuild(t, 10)
	base := startAssetServer(t, r)
	for _, p := range []string{"/assets/00112233445566ff/other.mp4", "/assets/../content.db", "/assets/%2e%2e/content.db"} {
		resp, err := http.Get(base + p)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("%s: status %d, want 404", p, resp.StatusCode)
		}
	}
}

func TestSlowReaderCompletes(t *testing.T) {
	if testing.Short() {
		t.Skip("slow download")
	}
	r, data, url := assetBuild(t, 16<<20)
	base := startAssetServer(t, r)
	resp, err := http.Get(base + url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	var got bytes.Buffer
	buf := make([]byte, 256<<10)
	start := time.Now()
	for {
		n, err := resp.Body.Read(buf)
		got.Write(buf[:n])
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("download broke after %d bytes: %v", got.Len(), err)
		}
		time.Sleep(40 * time.Millisecond) // a slow client
	}
	if !bytes.Equal(got.Bytes(), data) {
		t.Fatalf("got %d bytes, want %d", got.Len(), len(data))
	}
	t.Logf("slow download of %d MB took %s", len(data)>>20, time.Since(start).Round(time.Millisecond))
}
