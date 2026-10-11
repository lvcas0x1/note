package content

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMigrateDocsifySidebarGolden(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"ai/bedrock.md":                        "# Bedrock",
		"ai/claude.md":                         "# Claude",
		"aws/sessionmanager/sessionmanager.md": "# SM",
		"azure/README.md":                      "# Azure",
		"azure/cert/ai900.md":                  "# AI 900",
		"tools/a.md":                           "# A",
		"tools/b.md":                           "# B",
		"_index.md":                            "# Usage",
	})
	source := []byte(`<!-- docsify/_sidebar.md -->

- [AI](ai/)
  - [Bedrock](ai/bedrock.md)
  - [Claude](ai/claude "tooltip")
- [AWS](aws/)
  - [Session Manager](aws/sessionmanager/)
- [Azure](azure/)
  - Certifications
    - [AI 900](/azure/cert/ai900.md)
- [Tools](tools/)
- [External](https://example.com)

---

- [Usage](_index.md)
`)
	nodes, warnings, err := MigrateDocsifySidebar(source, root)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := json.MarshalIndent(nodes, "", "  ")
	golden := filepath.Join("testdata", "migrate", "sidebar.golden.json")
	if *update {
		_ = os.MkdirAll(filepath.Dir(golden), 0o750)
		_ = os.WriteFile(golden, append(got, '\n'), 0o600)
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("missing golden (run with -update): %v", err)
	}
	if string(append(got, '\n')) != string(want) {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
	if len(warnings) != 2 {
		t.Errorf("warnings = %q, want the folder listing and the external link", warnings)
	}

	// The output must pass sidebar validation.
	raw, _ := json.Marshal(nodes)
	writeFiles(t, root, map[string]string{SidebarFile: string(raw)})
	if _, err := LoadSidebar(root); err != nil {
		t.Fatalf("migrated sidebar does not validate: %v", err)
	}
}

// TestMigrateRealDocsify runs on a temporary copy of the user's docsify notes when present.
func TestMigrateRealDocsify(t *testing.T) {
	home, _ := os.UserHomeDir()
	src := filepath.Join(home, "Documents", "github", "docsify")
	if _, err := os.Stat(filepath.Join(src, "_sidebar.md")); err != nil {
		t.Skip("docsify notes not present")
	}
	root := filepath.Join(t.TempDir(), "docsify")
	if out, err := exec.Command("cp", "-R", src, root).CombinedOutput(); err != nil {
		t.Fatalf("copy: %v %s", err, out)
	}
	source, err := os.ReadFile(filepath.Join(root, "_sidebar.md"))
	if err != nil {
		t.Fatal(err)
	}
	nodes, warnings, err := MigrateDocsifySidebar(source, root)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(nodes)
	writeFiles(t, root, map[string]string{SidebarFile: string(raw)})
	if _, err := LoadSidebar(root); err != nil {
		t.Fatalf("migrated docsify sidebar does not validate: %v", err)
	}
	var sm *Node
	var find func([]Node)
	find = func(ns []Node) {
		for i := range ns {
			if ns[i].Title == "Session Manager" {
				sm = &ns[i]
			}
			find(ns[i].Children)
		}
	}
	find(nodes)
	if sm == nil || sm.Path != "aws/sessionmanager/sessionmanager.md" {
		t.Errorf("Session Manager node = %+v", sm)
	}
	t.Logf("%d top-level nodes, warnings: %q", len(nodes), warnings)
}
