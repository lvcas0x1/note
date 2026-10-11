package content

import (
	"strings"
	"testing"
)

func TestBuildExcludesDrafts(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"sidebar.json":   `[{"title":"Notes","children":[{"title":"Pub","path":"pub.md"},{"title":"Secret","path":"secret.md"}]},{"title":"Drafts only","children":[{"title":"Wip","path":"wip.md"}]}]`,
		"pub.md":         "---\ndate: 2026-10-05\n---\n# Pub\n\n![shared](img/shared.png)\n![own](img/pub.png)\n",
		"secret.md":      "---\ndate: 2026-10-05\ndraft: true\n---\n# Secret plan\n\n![shared](img/shared.png)\n![only](img/secret.png)\n",
		"wip.md":         "---\ndate: 2026-10-06\ndraft: true\n---\n# WIP\n",
		"img/shared.png": "shared", "img/pub.png": "pub", "img/secret.png": "secret",
	})
	res, err := Build(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Articles) != 1 || res.Articles[0].ID != "pub" {
		t.Fatalf("articles = %+v", res.Articles)
	}
	// Sidebar: the draft node is omitted and the draft-only category is pruned.
	if len(res.Sidebar) != 1 || len(res.Sidebar[0].Children) != 1 || res.Sidebar[0].Children[0].Path != "pub.md" {
		t.Fatalf("sidebar = %+v", res.Sidebar)
	}
	warned := strings.Join(res.Warnings, "\n")
	for _, want := range []string{`"Secret" omitted`, `"Wip" omitted`} {
		if !strings.Contains(warned, want) {
			t.Errorf("warnings %q missing %q", warned, want)
		}
	}
	// Assets: the shared image is published, the draft-only image is not.
	var names []string
	for _, a := range res.Assets {
		names = append(names, a.Name)
	}
	if strings.Join(names, ",") == "" || strings.Contains(strings.Join(names, ","), "secret.png") {
		t.Fatalf("assets = %v", names)
	}
	if !containsAll(names, "shared.png", "pub.png") {
		t.Fatalf("assets = %v, want shared.png and pub.png", names)
	}
}

func TestDraftAssetsAreNotResolved(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"sidebar.json": `[{"title":"A","path":"a.md"}]`,
		"a.md":         "---\ndate: 2026-01-01\n---\nok\n",
		"d.md":         "---\ndate: 2026-01-01\ndraft: true\n---\n![x](missing.png)\n",
	})
	if _, err := Build(root, Options{}); err != nil {
		t.Fatalf("a draft's missing image must not fail the build: %v", err)
	}
}

func containsAll(have []string, want ...string) bool {
	set := map[string]bool{}
	for _, h := range have {
		set[h] = true
	}
	for _, w := range want {
		if !set[w] {
			return false
		}
	}
	return true
}
