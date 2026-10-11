package content

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Result is a validated, compiled content folder with drafts removed.
type Result struct {
	Sidebar  []Node
	Articles []*Article // published articles, sorted by ID
	Assets   []*Asset   // assets used by published articles
	Warnings []string
}

// Options configures a build.
type Options struct {
	Fetcher *Fetcher // downloads remote embeds; nil forbids remote references
}

// Build compiles every Markdown file under root.
func Build(root string, opts Options) (*Result, error) {
	sidebar, err := LoadSidebar(root)
	if err != nil {
		return nil, err
	}
	titles := map[string]string{}
	collectTitles(sidebar, titles)

	files, err := markdownFiles(root)
	if err != nil {
		return nil, err
	}
	resolver := NewResolver(root, opts.Fetcher)
	res := &Result{}
	published := map[string]bool{}
	drafts := map[string]bool{}
	for _, rel := range files {
		a, err := compileArticle(root, rel, titles[rel], resolver)
		if err != nil {
			return nil, err
		}
		res.Warnings = append(res.Warnings, a.Warnings...)
		if a.Draft {
			drafts[rel] = true
			continue
		}
		published[rel] = true
		res.Articles = append(res.Articles, a)
	}
	sort.Slice(res.Articles, func(i, j int) bool { return res.Articles[i].ID < res.Articles[j].ID })
	res.Sidebar = pruneDrafts(sidebar, drafts, &res.Warnings)
	res.Assets = resolver.Published(published)
	return res, nil
}

func collectTitles(nodes []Node, out map[string]string) {
	for _, n := range nodes {
		if n.Path != "" {
			out[n.Path] = n.Title
		}
		collectTitles(n.Children, out)
	}
}

// markdownFiles lists content-relative .md paths, skipping hidden files and folders.
func markdownFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p != root && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() || !strings.EqualFold(filepath.Ext(d.Name()), ".md") {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	sort.Strings(files)
	return files, err
}

// pruneDrafts drops sidebar nodes that point to drafts and categories left empty.
func pruneDrafts(nodes []Node, drafts map[string]bool, warnings *[]string) []Node {
	var out []Node
	for _, n := range nodes {
		if n.Path != "" && drafts[n.Path] {
			*warnings = append(*warnings, fmt.Sprintf("sidebar.json: node %q omitted because %s is a draft", n.Title, n.Path))
			n.Path = ""
		}
		n.Children = pruneDrafts(n.Children, drafts, warnings)
		if n.Path == "" && len(n.Children) == 0 {
			continue
		}
		out = append(out, n)
	}
	return out
}
