package content

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite golden files")

// hashRe hides content hashes so golden files stay readable.
var hashRe = regexp.MustCompile(`/assets/[0-9a-f]{16}/`)

func TestEmbedGolden(t *testing.T) {
	assets := map[string]string{
		"diagrams/net.drawio": "<mxfile/>",
		"api/swagger.json":    `{"swagger":"2.0"}`,
		"videos/demo.mp4":     "mp4",
		"sketch.excalidraw":   `{"type":"excalidraw"}`,
		"img/a.png":           "png",
	}
	mds, _ := filepath.Glob("testdata/embeds/*.md")
	for _, md := range mds {
		name := strings.TrimSuffix(filepath.Base(md), ".md")
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			src, err := os.ReadFile(md)
			if err != nil {
				t.Fatal(err)
			}
			files := map[string]string{name + ".md": string(src)}
			for k, v := range assets {
				files[k] = v
			}
			writeFiles(t, root, files)
			a, err := CompileArticle(root, name+".md", "")
			if err != nil {
				t.Fatal(err)
			}
			got := hashRe.ReplaceAllString(a.HTML, "/assets/HASH/")
			golden := strings.TrimSuffix(md, ".md") + ".golden.html"
			if *update {
				if err := os.WriteFile(golden, []byte(got), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("missing golden file (run with -update): %v", err)
			}
			if got != string(want) {
				t.Errorf("HTML mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", name, got, want)
			}
		})
	}
}

func TestEChartsInvalidJSON(t *testing.T) {
	_, err := compileOne(t, map[string]string{"c.md": "---\ndate: 2026-01-01\n---\n# Chart\n\ntext\n\n```echarts\n{\"series\": [\n```\n"}, "c.md", "")
	if err == nil || !strings.Contains(err.Error(), "c.md:8:") || !strings.Contains(err.Error(), "echarts") {
		t.Fatalf("err = %v, want file and line 8", err)
	}
	_, err = compileOne(t, map[string]string{"d.md": "---\ndate: 2026-01-01\n---\n```echarts\n[1, 2]\n```\n"}, "d.md", "")
	if err == nil || !strings.Contains(err.Error(), "JSON object") {
		t.Fatalf("err = %v, want JSON object error", err)
	}
}
