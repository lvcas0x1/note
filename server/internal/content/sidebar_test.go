package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoadSidebar(t *testing.T) {
	cases := []struct {
		name    string
		files   map[string]string
		wantErr []string
	}{
		{
			name: "nested categories keep nesting and order",
			files: map[string]string{
				"sidebar.json": `[{"title":"Cloud","children":[{"title":"AWS","children":[{"title":"S3","path":"aws/s3.md"}]}]},{"title":"Go","path":"go.md"}]`,
				"aws/s3.md":    "# S3",
				"go.md":        "# Go",
			},
		},
		{
			name:    "missing sidebar",
			files:   map[string]string{"a.md": "# A"},
			wantErr: []string{"sidebar.json not found"},
		},
		{
			name:    "node without title",
			files:   map[string]string{"sidebar.json": `[{"title":"A","children":[{"path":"a.md"}]}]`, "a.md": ""},
			wantErr: []string{"invalid node at /0/children/0"},
		},
		{
			name:    "node with neither path nor children",
			files:   map[string]string{"sidebar.json": `[{"title":"A","path":"a.md"},{"title":"Lonely"}]`, "a.md": ""},
			wantErr: []string{"invalid node at /1"},
		},
		{
			name:    "blank title",
			files:   map[string]string{"sidebar.json": `[{"title":"  ","path":"a.md"}]`, "a.md": ""},
			wantErr: []string{"/0: title is blank"},
		},
		{
			name:    "broken path",
			files:   map[string]string{"sidebar.json": `[{"title":"Gone","path":"missing.md"}]`},
			wantErr: []string{`node "Gone"`, `"missing.md" does not exist`},
		},
		{
			name:    "path outside content",
			files:   map[string]string{"sidebar.json": `[{"title":"Evil","path":"../evil.md"}]`},
			wantErr: []string{"outside the content folder"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFiles(t, root, tc.files)
			nodes, err := LoadSidebar(root)
			if len(tc.wantErr) == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(nodes) != 2 || nodes[0].Title != "Cloud" || nodes[0].Children[0].Children[0].Path != "aws/s3.md" || nodes[1].Title != "Go" {
					t.Fatalf("nodes = %+v", nodes)
				}
				return
			}
			if err == nil {
				t.Fatal("expected an error")
			}
			for _, want := range tc.wantErr {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not contain %q", err, want)
				}
			}
		})
	}
}

func TestResolveInsideRejectsSymlinkEscape(t *testing.T) {
	outside := t.TempDir()
	writeFiles(t, outside, map[string]string{"secret.md": "x"})
	root := t.TempDir()
	if err := os.Symlink(filepath.Join(outside, "secret.md"), filepath.Join(root, "link.md")); err != nil {
		t.Skip("symlinks unsupported:", err)
	}
	if _, err := ResolveInside(root, "link.md"); err == nil || !strings.Contains(err.Error(), "outside") {
		t.Fatalf("err = %v, want outside-content error", err)
	}
}
