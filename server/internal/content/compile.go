package content

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
)

// CompileArticle parses and compiles one Markdown file under root on its own.
// Remote references are not allowed. Use Build to compile a whole content folder.
func CompileArticle(root, rel, sidebarTitle string) (*Article, error) {
	return compileArticle(root, rel, sidebarTitle, NewResolver(root, nil))
}

func compileArticle(root, rel, sidebarTitle string, resolver *Resolver) (*Article, error) {
	source, err := os.ReadFile(joinRel(root, rel))
	if err != nil {
		return nil, err
	}
	md := newMarkdown(goldmark.WithRendererOptions(embedRendererOption(), codeRendererOption()))
	p := parseMarkdown(md, source)
	a := &Article{ID: ArticleID(rel), Path: rel}
	if err := readMetadata(root, rel, sidebarTitle, p, a); err != nil {
		return nil, err
	}
	if a.Draft {
		// Drafts never leave the build, so their body and assets are not processed.
		return a, nil
	}
	a.Text = plainText(p.doc, source)
	dropLeadingTitle(p, a.Title)
	if err := transformEmbeds(rel, p); err != nil {
		return nil, err
	}
	if err := resolveAssets(rel, p, resolver, a); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, source, p.doc); err != nil {
		return nil, fmt.Errorf("%s: render: %w", rel, err)
	}
	a.HTML = sanitizeHTML(buf.String())
	return a, nil
}
