package content

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
)

// resolveAssets rewrites image and file-embed sources to published asset URLs.
func resolveAssets(rel string, p parsed, resolver *Resolver, a *Article) error {
	var firstErr error
	_ = ast.Walk(p.doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch t := n.(type) {
		case *ast.Image:
			src := string(t.Destination)
			url, err := resolver.Resolve(rel, "image", src)
			if err != nil {
				firstErr = fmt.Errorf("%s: image %q: %w", rel, src, err)
				return ast.WalkStop, nil
			}
			t.Destination = []byte(url)
			a.Assets = append(a.Assets, AssetRef{Kind: "image", Source: src})
		case *Embed:
			if t.Src == "" {
				return ast.WalkContinue, nil
			}
			src := t.Src
			url, err := resolver.Resolve(rel, t.EmbedType, src)
			if err != nil {
				firstErr = fmt.Errorf("%s: %s embed %q: %w", rel, t.EmbedType, src, err)
				return ast.WalkStop, nil
			}
			t.Src = url
			a.Assets = append(a.Assets, AssetRef{Kind: t.EmbedType, Source: src})
		}
		return ast.WalkContinue, nil
	})
	return firstErr
}
