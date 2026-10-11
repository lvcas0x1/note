package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

// SummaryMaxRunes is the length a summary taken from the first paragraph is cut to.
const SummaryMaxRunes = 200

// DateLayout is the article date format (ISO 8601 calendar date).
const DateLayout = "2006-01-02"

var knownFrontMatterKeys = map[string]bool{"title": true, "date": true, "summary": true, "tags": true, "draft": true}

// Article is one compiled Markdown file.
type Article struct {
	ID       string   // content-relative path without ".md"
	Path     string   // content-relative path with ".md"
	Title    string   //
	Date     string   // YYYY-MM-DD
	Summary  string   //
	Tags     []string //
	Draft    bool     //
	HTML     string   // compiled, sanitized HTML body
	Text     string   // plain text of the body, for full-text search
	Assets   []AssetRef
	Warnings []string
}

// AssetRef is a file an article references, resolved during compilation.
type AssetRef struct {
	Kind   string // image, drawio, excalidraw, openapi, mp4
	Source string // as written in the Markdown
}

// ArticleID returns the article ID for a content-relative Markdown path.
func ArticleID(rel string) string { return strings.TrimSuffix(path.Clean(rel), ".md") }

func newMarkdown(extra ...goldmark.Option) goldmark.Markdown {
	opts := []goldmark.Option{
		goldmark.WithExtensions(extension.GFM, extension.Footnote, &frontmatter.Extender{}),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		// Raw HTML is kept here and sanitized after rendering.
		goldmark.WithRendererOptions(html.WithUnsafe()),
	}
	return goldmark.New(append(opts, extra...)...)
}

// parsed is the intermediate result shared by the compile steps.
type parsed struct {
	doc    ast.Node
	source []byte
	ctx    parser.Context
}

func parseMarkdown(md goldmark.Markdown, source []byte) parsed {
	ctx := parser.NewContext()
	doc := md.Parser().Parse(text.NewReader(source), parser.WithContext(ctx))
	return parsed{doc: doc, source: source, ctx: ctx}
}

// readMetadata validates front matter and fills in title, date, and summary.
func readMetadata(root, rel, sidebarTitle string, p parsed, a *Article) error {
	fm := map[string]any{}
	if data := frontmatterData(p); data != nil {
		if err := data.Decode(&fm); err != nil {
			return fmt.Errorf("%s: front matter: %w", rel, err)
		}
	}
	normalizeFrontMatter(fm)
	if err := validateFrontMatter(fm); err != nil {
		return fmt.Errorf("%s: front matter: %w", rel, err)
	}
	var unknown []string
	for k := range fm {
		if !knownFrontMatterKeys[k] {
			unknown = append(unknown, k)
		}
	}
	sort.Strings(unknown)
	for _, k := range unknown {
		a.Warnings = append(a.Warnings, fmt.Sprintf("%s: unknown front matter key %q ignored", rel, k))
	}

	a.Title, _ = fm["title"].(string)
	a.Title = strings.TrimSpace(a.Title)
	a.Date, _ = fm["date"].(string)
	a.Summary, _ = fm["summary"].(string)
	a.Draft, _ = fm["draft"].(bool)
	if tags, ok := fm["tags"].([]any); ok {
		for _, t := range tags {
			a.Tags = append(a.Tags, fmt.Sprint(t))
		}
	}

	if a.Title == "" {
		a.Title = firstHeading(p, 1)
	}
	if a.Title == "" {
		a.Title = strings.TrimSpace(sidebarTitle)
	}
	if a.Title == "" {
		a.Title = strings.TrimSuffix(path.Base(rel), ".md")
	}
	if a.Date == "" {
		info, err := os.Stat(joinRel(root, rel))
		if err != nil {
			return err
		}
		a.Date = info.ModTime().In(time.Local).Format(DateLayout)
		a.Warnings = append(a.Warnings, fmt.Sprintf("%s: no date in front matter; using file modification date %s", rel, a.Date))
	}
	if a.Summary == "" {
		a.Summary = truncateRunes(firstParagraphText(p), SummaryMaxRunes)
	}
	return nil
}

// normalizeFrontMatter converts YAML-decoded values into their JSON forms,
// so `date: 2026-10-11` (a YAML timestamp) validates as a date string.
func normalizeFrontMatter(fm map[string]any) {
	for k, v := range fm {
		if t, ok := v.(time.Time); ok {
			if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0 {
				fm[k] = t.Format(DateLayout)
			} else {
				fm[k] = t.Format(time.RFC3339)
			}
		}
	}
}

func validateFrontMatter(fm map[string]any) error {
	_, sch, err := schemas()
	if err != nil {
		return err
	}
	raw, err := json.Marshal(fm)
	if err != nil {
		return err
	}
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return err
	}
	if err := sch.Validate(inst); err != nil {
		return fmt.Errorf("invalid value at %s (got %s)", describeValidation(err), badValues(fm, err))
	}
	if t, ok := fm["title"].(string); ok && strings.TrimSpace(t) == "" {
		return fmt.Errorf("title is blank")
	}
	return nil
}

func badValues(fm map[string]any, err error) string {
	ve, ok := err.(*jsonschema.ValidationError) //nolint:errorlint // concrete validator type
	if !ok {
		return "?"
	}
	var vals []string
	var walk func(e *jsonschema.ValidationError)
	walk = func(e *jsonschema.ValidationError) {
		if len(e.Causes) == 0 && len(e.InstanceLocation) > 0 {
			vals = append(vals, fmt.Sprintf("%v", fm[e.InstanceLocation[0]]))
		}
		for _, c := range e.Causes {
			walk(c)
		}
	}
	walk(ve)
	if len(vals) == 0 {
		return "?"
	}
	return fmt.Sprintf("%q", vals[0])
}

func firstHeading(p parsed, level int) string {
	var title string
	_ = ast.Walk(p.doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if h, ok := n.(*ast.Heading); ok && entering && h.Level == level {
			title = strings.TrimSpace(plainText(h, p.source))
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return title
}

func firstParagraphText(p parsed) string {
	var out string
	_ = ast.Walk(p.doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if para, ok := n.(*ast.Paragraph); ok && entering {
			if s := strings.TrimSpace(plainText(para, p.source)); s != "" {
				out = s
				return ast.WalkStop, nil
			}
		}
		return ast.WalkContinue, nil
	})
	return out
}

// plainText returns the visible text of n with soft line breaks as spaces.
func plainText(n ast.Node, source []byte) string {
	var b strings.Builder
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch t := c.(type) {
		case *ast.Text:
			b.Write(t.Segment.Value(source))
			if t.SoftLineBreak() || t.HardLineBreak() {
				b.WriteByte(' ')
			}
		case *ast.String:
			b.Write(t.Value)
		case *ast.FencedCodeBlock, *ast.CodeBlock:
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			lines := c.Lines()
			for i := 0; i < lines.Len(); i++ {
				seg := lines.At(i)
				b.Write(seg.Value(source))
			}
			b.WriteByte(' ')
			return ast.WalkSkipChildren, nil
		case *ast.Paragraph, *ast.Heading, *ast.ListItem, *extast.TableCell:
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
		}
		return ast.WalkContinue, nil
	})
	return strings.Join(strings.Fields(b.String()), " ")
}

func truncateRunes(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	r := []rune(s)
	return strings.TrimSpace(string(r[:n]))
}

func joinRel(root, rel string) string {
	return root + string(os.PathSeparator) + strings.ReplaceAll(rel, "/", string(os.PathSeparator))
}

// dropLeadingTitle removes a level-1 heading at the start of the body when it
// repeats the article title, because the article view already shows the title.
func dropLeadingTitle(p parsed, title string) {
	first := p.doc.FirstChild()
	for first != nil && first.Kind() != ast.KindHeading {
		if first.Kind() == ast.KindParagraph || first.Kind() == ast.KindList || first.Kind() == ast.KindFencedCodeBlock {
			return // the body starts with content, not with the title
		}
		first = first.NextSibling()
	}
	h, ok := first.(*ast.Heading)
	if ok && h.Level == 1 && strings.TrimSpace(plainText(h, p.source)) == title {
		p.doc.RemoveChild(p.doc, h)
	}
}

func frontmatterData(p parsed) *frontmatter.Data { return frontmatter.Get(p.ctx) }
