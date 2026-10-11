package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"notes/internal/content"
)

func sampleResult(t *testing.T) *content.Result {
	t.Helper()
	img := filepath.Join(t.TempDir(), "a.png")
	if err := os.WriteFile(img, []byte("png"), 0o600); err != nil {
		t.Fatal(err)
	}
	return &content.Result{
		Sidebar: []content.Node{{Title: "Cloud", Children: []content.Node{{Title: "S3", Path: "aws/s3.md"}}}},
		Articles: []*content.Article{{ID: "aws/s3", Path: "aws/s3.md", Title: "S3", Date: "2026-10-05",
			Summary: "Buckets", Tags: []string{"aws"}, HTML: "<p>x</p>", Text: "S3 buckets"}},
		Assets: []*content.Asset{{Key: "img/a.png", URL: "/assets/0123456789abcdef/a.png", Name: "a.png",
			Hash: "0123456789abcdef0123", Size: 3, Kind: "image", ContentType: "image/png", LocalPath: img}},
	}
}

func TestWriteBuildSwapsCurrent(t *testing.T) {
	data := t.TempDir()
	first, err := WriteBuild(data, sampleResult(t))
	if err != nil {
		t.Fatal(err)
	}
	cur, err := CurrentDir(data)
	if err != nil || cur != first {
		t.Fatalf("current = %s err = %v, want %s", cur, err, first)
	}
	if _, err := os.Stat(filepath.Join(first, AssetsDir, "0123456789abcdef", "a.png")); err != nil {
		t.Fatalf("asset not copied: %v", err)
	}
	second, err := WriteBuild(data, sampleResult(t))
	if err != nil {
		t.Fatal(err)
	}
	if cur, _ := CurrentDir(data); cur != second {
		t.Fatalf("current = %s, want %s", cur, second)
	}
}

func TestFailedBuildLeavesCurrent(t *testing.T) {
	data := t.TempDir()
	good, err := WriteBuild(data, sampleResult(t))
	if err != nil {
		t.Fatal(err)
	}
	bad := sampleResult(t)
	bad.Assets[0].LocalPath = filepath.Join(t.TempDir(), "gone.png")
	if _, err := WriteBuild(data, bad); err == nil {
		t.Fatal("expected failure")
	}
	if cur, _ := CurrentDir(data); cur != good {
		t.Fatalf("current moved to %s after a failed build", cur)
	}
	entries, _ := os.ReadDir(filepath.Join(data, BuildsDir))
	if len(entries) != 1 {
		t.Fatalf("partial build left behind: %d builds", len(entries))
	}
}

func TestKeepsOnlyTheNewestBuild(t *testing.T) {
	data := t.TempDir()
	var last string
	for range 5 {
		dir, err := WriteBuild(data, sampleResult(t))
		if err != nil {
			t.Fatal(err)
		}
		last = dir
	}
	entries, _ := os.ReadDir(filepath.Join(data, BuildsDir))
	if len(entries) != KeepBuilds {
		t.Fatalf("%d builds kept, want %d", len(entries), KeepBuilds)
	}
	if cur, _ := CurrentDir(data); cur != last {
		t.Fatalf("current = %s, want newest %s", cur, last)
	}
}

func TestNoBuild(t *testing.T) {
	if _, err := CurrentDir(t.TempDir()); !errors.Is(err, ErrNoBuild) {
		t.Fatalf("err = %v, want ErrNoBuild", err)
	}
}
