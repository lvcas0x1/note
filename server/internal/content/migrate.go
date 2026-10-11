package content

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
)

// MigrateDocsifySidebar converts a docsify _sidebar.md (a nested Markdown link
// list) into sidebar nodes. root is the docsify content folder that links are
// relative to. Article files are never modified.
func MigrateDocsifySidebar(source []byte, root string) ([]Node, []string, error) {
	p := parseMarkdown(newMarkdown(), source)
	m := migrator{root: root, source: source}
	var nodes []Node
	for c := p.doc.FirstChild(); c != nil; c = c.NextSibling() {
		if list, ok := c.(*ast.List); ok {
			nodes = append(nodes, m.list(list)...)
		}
	}
	if len(nodes) == 0 {
		return nil, m.warnings, fmt.Errorf("no sidebar entries found")
	}
	return nodes, m.warnings, nil
}

type migrator struct {
	root     string
	source   []byte
	warnings []string
}

func (m *migrator) list(list *ast.List) []Node {
	var out []Node
	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		if n, ok := m.item(item); ok {
			out = append(out, n)
		}
	}
	return out
}

func (m *migrator) item(item ast.Node) (Node, bool) {
	var n Node
	var link string
	for c := item.FirstChild(); c != nil; c = c.NextSibling() {
		switch t := c.(type) {
		case *ast.Paragraph, *ast.TextBlock:
			if l := firstLink(t); l != nil {
				link = string(l.Destination)
				n.Title = plainText(l, m.source)
			} else {
				n.Title = plainText(t, m.source)
			}
		case *ast.List:
			n.Children = m.list(t)
		}
	}
	n.Title = strings.TrimSpace(n.Title)
	if n.Title == "" {
		return n, false
	}
	if !m.resolveLink(&n, link) {
		return n, false
	}
	if n.Path == "" && len(n.Children) == 0 {
		m.warnings = append(m.warnings, fmt.Sprintf("%q has no article and no children; skipped", n.Title))
		return n, false
	}
	return n, true
}

func firstLink(n ast.Node) *ast.Link {
	var found *ast.Link
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if l, ok := c.(*ast.Link); ok && entering {
			found = l
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return found
}

// resolveLink sets n.Path from a docsify link. It returns false when the
// entry should be dropped (it links outside the content).
func (m *migrator) resolveLink(n *Node, link string) bool {
	link = strings.TrimSpace(strings.SplitN(link, "#", 2)[0])
	switch {
	case link == "":
		return true
	case strings.Contains(link, "://") || strings.HasPrefix(link, "mailto:"):
		m.warnings = append(m.warnings, fmt.Sprintf("%q links outside the content (%s); skipped", n.Title, link))
		return false
	}
	rel := path.Clean(strings.TrimPrefix(link, "/"))
	if strings.HasSuffix(link, "/") || rel == "." {
		m.folder(n, rel)
		return true
	}
	if !strings.EqualFold(path.Ext(rel), ".md") {
		rel += ".md" // docsify allows links without the extension
	}
	if _, err := ResolveInside(m.root, rel); err != nil {
		m.warnings = append(m.warnings, fmt.Sprintf("%q: %v; kept without a path", n.Title, err))
		return true
	}
	n.Path = rel
	return true
}

// folder applies the folder-link rules: README.md, then a lone .md file when
// the node has no children, otherwise a category whose children are the
// folder's .md files if it has none of its own.
func (m *migrator) folder(n *Node, dir string) {
	readme := path.Join(dir, "README.md")
	if _, err := ResolveInside(m.root, readme); err == nil {
		n.Path = readme
		return
	}
	files := mdFilesIn(filepath.Join(m.root, filepath.FromSlash(dir)))
	if len(n.Children) > 0 {
		return
	}
	switch len(files) {
	case 0:
	case 1:
		n.Path = path.Join(dir, files[0])
	default:
		for _, f := range files {
			n.Children = append(n.Children, Node{Title: strings.TrimSuffix(f, path.Ext(f)), Path: path.Join(dir, f)})
		}
		m.warnings = append(m.warnings, fmt.Sprintf("%q: folder %s/ has several articles; listed them as children", n.Title, dir))
	}
}

func mdFilesIn(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".md") {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}
