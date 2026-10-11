package content

import (
	"strings"
	"testing"
)

func TestScriptsAndHandlersAreStripped(t *testing.T) {
	a, err := compileOne(t, map[string]string{"x.md": "---\ndate: 2026-01-01\n---\n# T\n\n<script>alert(1)</script>\n\n<a href=\"#\" onclick=\"steal()\">x</a>\n\n<img src=\"data:image/png;base64,iVBORw0KGgo=\" onerror=\"boom()\">\n\n<span style=\"color:red\">Red</span>\n"}, "x.md", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"<script", "alert(1)", "onclick", "onerror", "steal()"} {
		if strings.Contains(a.HTML, bad) {
			t.Errorf("HTML still contains %q:\n%s", bad, a.HTML)
		}
	}
	if !strings.Contains(a.HTML, `style="color: red"`) && !strings.Contains(a.HTML, `style="color:red"`) {
		t.Errorf("span color was removed:\n%s", a.HTML)
	}
}

func TestGFMAndHighlighting(t *testing.T) {
	md := "---\ndate: 2026-01-01\n---\n# Title\n\n## Sub\n\n| a | b |\n|---|---|\n| 1 | 2 |\n\n- [x] done\n- [ ] todo\n\nNote.[^1]\n\n~~old~~ https://example.com\n\n```go\nfunc main() {}\n```\n\n[^1]: The footnote.\n"
	a, err := compileOne(t, map[string]string{"g.md": md}, "g.md", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"<table>", `<input checked="" disabled="" type="checkbox"`, `type="checkbox"`, `<del>old</del>`,
		`href="https://example.com"`, `href="#fn:1"`, `id="fn:1"`, `<code class="language-go">`, `<span class="kd">func</span>`,
		`<h2 id="sub">`,
	} {
		if !strings.Contains(a.HTML, want) {
			t.Errorf("HTML missing %q:\n%s", want, a.HTML)
		}
	}
}

func TestEmbedPlaceholdersSurviveSanitizing(t *testing.T) {
	a, err := compileOne(t, map[string]string{
		"e.md":     "---\ndate: 2026-01-01\n---\n![drawio](d.drawio)\n\n```mermaid\nflowchart LR; A-->B\n```\n",
		"d.drawio": "<mxfile/>",
	}, "e.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(a.HTML, `<div data-embed="drawio" data-src="/assets/`) || !strings.Contains(a.HTML, `<pre data-embed="mermaid">flowchart LR; A--&gt;B`) {
		t.Fatalf("placeholders lost:\n%s", a.HTML)
	}
}

func TestChromaCSS(t *testing.T) {
	css, err := ChromaCSS()
	if err != nil || !strings.Contains(css, ".chroma") {
		t.Fatalf("css = %.80q err = %v", css, err)
	}
}

// chroma classes with digits (s1, s2, c1, …) must survive sanitizing.
func TestChromaClassesWithDigitsSurvive(t *testing.T) {
	a, err := compileOne(t, map[string]string{"p.md": "---\ndate: 2026-01-01\n---\n```python\nprint('Hi', \"x\")  # note\nn = 0x1F\n```\n"}, "p.md", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`<span class="s1">`, `<span class="s2">`, `<span class="c1">`, `<span class="mh">`} {
		if !strings.Contains(a.HTML, want) {
			t.Errorf("HTML lost %s:\n%s", want, a.HTML)
		}
	}
	if strings.Contains(a.HTML, "<span>") {
		t.Errorf("a highlighted token lost its class:\n%s", a.HTML)
	}
}
