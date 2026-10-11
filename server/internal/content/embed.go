package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Embed kinds written into data-embed.
const (
	EmbedDrawio     = "drawio"
	EmbedExcalidraw = "excalidraw"
	EmbedOpenAPI    = "openapi"
	EmbedMP4        = "mp4"
	EmbedMermaid    = "mermaid"
	EmbedECharts    = "echarts"
)

// imageEmbeds maps image alt text to embed kind.
var imageEmbeds = map[string]string{
	"drawio":     EmbedDrawio,
	"excalidraw": EmbedExcalidraw,
	"openapi":    EmbedOpenAPI,
	"swagger":    EmbedOpenAPI,
	"mp4":        EmbedMP4,
}

// fenceEmbeds maps fenced code languages to embed kind.
var fenceEmbeds = map[string]string{"mermaid": EmbedMermaid, "echarts": EmbedECharts}

// KindEmbed is the AST node kind of an embed placeholder.
var KindEmbed = ast.NewNodeKind("Embed")

// Embed is a placeholder that the web app turns into a diagram, chart, API reference, or video.
type Embed struct {
	ast.BaseBlock
	EmbedType string
	Src       string // source reference for file embeds; resolved to a URL later
	Body      string // inline source for fence embeds
	Inline    bool   // the embed sat inside a paragraph with other content
}

// Kind implements ast.Node.
func (*Embed) Kind() ast.NodeKind { return KindEmbed }

// Dump implements ast.Node.
func (e *Embed) Dump(source []byte, level int) {
	ast.DumpHelper(e, source, level, map[string]string{"EmbedType": e.EmbedType, "Src": e.Src}, nil)
}

// IsRaw reports that the node has no Markdown children.
func (*Embed) IsRaw() bool { return true }

// transformEmbeds replaces embed images and fences with Embed nodes.
func transformEmbeds(rel string, p parsed) error {
	var images []*ast.Image
	var fences []*ast.FencedCodeBlock
	_ = ast.Walk(p.doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch t := n.(type) {
		case *ast.Image:
			if _, ok := imageEmbeds[strings.ToLower(strings.TrimSpace(plainText(t, p.source)))]; ok {
				images = append(images, t)
			}
		case *ast.FencedCodeBlock:
			if _, ok := fenceEmbeds[strings.ToLower(string(t.Language(p.source)))]; ok {
				fences = append(fences, t)
			}
		}
		return ast.WalkContinue, nil
	})

	for _, img := range images {
		kind := imageEmbeds[strings.ToLower(strings.TrimSpace(plainText(img, p.source)))]
		e := &Embed{EmbedType: kind, Src: string(img.Destination)}
		parent := img.Parent()
		if para, ok := parent.(*ast.Paragraph); ok && onlyMeaningfulChild(para, img, p.source) {
			para.Parent().ReplaceChild(para.Parent(), para, e)
			continue
		}
		e.Inline = true
		parent.ReplaceChild(parent, img, e)
	}

	for _, f := range fences {
		kind := fenceEmbeds[strings.ToLower(string(f.Language(p.source)))]
		var body bytes.Buffer
		lines := f.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			body.Write(seg.Value(p.source))
		}
		if kind == EmbedECharts {
			if err := validateEChartsOption(body.Bytes()); err != nil {
				return fmt.Errorf("%s:%d: echarts block: %w", rel, fenceLine(f, p.source), err)
			}
		}
		f.Parent().ReplaceChild(f.Parent(), f, &Embed{EmbedType: kind, Body: body.String()})
	}
	return nil
}

func validateEChartsOption(body []byte) error {
	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if _, ok := v.(map[string]any); !ok {
		return fmt.Errorf("the option must be a JSON object")
	}
	return nil
}

// fenceLine returns the 1-based line of the opening fence.
func fenceLine(f *ast.FencedCodeBlock, source []byte) int {
	var start int
	if f.Info != nil {
		start = f.Info.Segment.Start
	} else if f.Lines().Len() > 0 {
		start = f.Lines().At(0).Start
	}
	return bytes.Count(source[:start], []byte("\n")) + 1
}

func onlyMeaningfulChild(para *ast.Paragraph, img *ast.Image, source []byte) bool {
	for c := para.FirstChild(); c != nil; c = c.NextSibling() {
		if c == img {
			continue
		}
		if t, ok := c.(*ast.Text); ok && strings.TrimSpace(string(t.Segment.Value(source))) == "" {
			continue
		}
		return false
	}
	return true
}

// embedRenderer renders Embed nodes as placeholders.
type embedRenderer struct{}

func (embedRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindEmbed, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		e := n.(*Embed)
		if e.Body != "" || e.EmbedType == EmbedMermaid || e.EmbedType == EmbedECharts {
			_, _ = fmt.Fprintf(w, `<pre data-embed="%s">`, e.EmbedType)
			_, _ = w.Write(util.EscapeHTML([]byte(e.Body)))
			_, _ = w.WriteString("</pre>\n")
			return ast.WalkSkipChildren, nil
		}
		_, _ = fmt.Fprintf(w, `<div data-embed="%s" data-src="`, e.EmbedType)
		_, _ = w.Write(util.EscapeHTML([]byte(e.Src)))
		_, _ = w.WriteString(`"></div>`)
		if !e.Inline {
			_, _ = w.WriteString("\n")
		}
		return ast.WalkSkipChildren, nil
	})
}

func embedRendererOption() renderer.Option {
	return renderer.WithNodeRenderers(util.Prioritized(embedRenderer{}, 100))
}
