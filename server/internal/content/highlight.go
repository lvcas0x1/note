package content

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Code highlighting colors: Catppuccin Latte hues, darkened in OKLab until each
// reaches about 5:1 on the code background, Latte base #eff1f5. The official
// Latte accents are pastel and mostly below 4.5:1 on light backgrounds. Text,
// subtext1, and base are the official Latte values.
const (
	codeBackground = "#eff1f5" // base
	codeText       = "#4c4f69" // text
	codeMuted      = "#5c5f77" // subtext1
	deepMauve      = "#8535eb"
	deepBlue       = "#135cea"
	deepYellow     = "#9f5200"
	deepGreen      = "#087800"
	deepPeach      = "#c52a00"
	deepSapphire   = "#007186"
	deepPink       = "#a93a8f"
	deepRed        = "#ce0436"
)

// codeStyle is the highlighting theme whose CSS the web app ships.
var codeStyle = chroma.MustNewStyle("catppuccin-latte-deep", chroma.StyleEntries{
	chroma.Background:            codeText + " bg:" + codeBackground,
	chroma.Error:                 deepRed,
	chroma.Keyword:               deepMauve,
	chroma.KeywordConstant:       deepPeach,
	chroma.KeywordType:           deepYellow,
	chroma.Name:                  codeText,
	chroma.NameAttribute:         deepBlue,
	chroma.NameBuiltin:           deepSapphire,
	chroma.NameBuiltinPseudo:     deepSapphire,
	chroma.NameClass:             deepYellow,
	chroma.NameConstant:          deepPeach,
	chroma.NameDecorator:         deepPink,
	chroma.NameException:         deepYellow,
	chroma.NameFunction:          deepBlue,
	chroma.NameLabel:             deepSapphire,
	chroma.NameNamespace:         deepYellow,
	chroma.NameTag:               deepMauve,
	chroma.NameVariable:          codeText,
	chroma.Literal:               deepPeach,
	chroma.LiteralDate:           deepPeach,
	chroma.LiteralString:         deepGreen,
	chroma.LiteralStringEscape:   deepPink,
	chroma.LiteralStringRegex:    deepPink,
	chroma.LiteralStringInterpol: deepPink,
	chroma.LiteralNumber:         deepPeach,
	chroma.Operator:              codeMuted,
	chroma.OperatorWord:          deepMauve,
	chroma.Punctuation:           codeMuted,
	chroma.Comment:               "italic " + codeMuted,
	chroma.CommentPreproc:        "noitalic " + deepPink,
	chroma.GenericDeleted:        deepRed,
	chroma.GenericInserted:       deepGreen,
	chroma.GenericEmph:           "italic",
	chroma.GenericStrong:         "bold",
	chroma.GenericHeading:        "bold " + deepBlue,
	chroma.GenericSubheading:     "bold " + deepMauve,
	chroma.GenericPrompt:         codeMuted,
	chroma.GenericOutput:         codeText,
	chroma.GenericError:          deepRed,
})

// lexerAliases maps fence languages that chroma does not know to a lexer.
var lexerAliases = map[string]string{"curl": "bash"}

var chromaFormatter = chromahtml.New(chromahtml.WithClasses(true), chromahtml.PreventSurroundingPre(true))

// codeRenderer highlights fenced and indented code blocks at build time with
// class-based chroma output, keeping a language-<lang> class on <code>.
type codeRenderer struct{}

func (codeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, renderCode)
	reg.Register(ast.KindCodeBlock, renderCode)
}

func renderCode(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	var lang string
	if f, ok := n.(*ast.FencedCodeBlock); ok {
		lang = strings.ToLower(string(f.Language(source)))
	}
	var code bytes.Buffer
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		code.Write(seg.Value(source))
	}

	_, _ = w.WriteString(`<pre class="chroma"><code`)
	if lang != "" {
		_, _ = w.WriteString(` class="language-`)
		_, _ = w.Write(util.EscapeHTML([]byte(lang)))
		_, _ = w.WriteString(`"`)
	}
	_, _ = w.WriteString(">")
	if !highlight(w, lang, code.String()) {
		_, _ = w.Write(util.EscapeHTML(code.Bytes()))
	}
	_, _ = w.WriteString("</code></pre>\n")
	return ast.WalkSkipChildren, nil
}

func highlight(w util.BufWriter, lang, code string) bool {
	if lang == "" {
		return false
	}
	lexer := lexers.Get(lang)
	if lexer == nil && lexerAliases[lang] != "" {
		lexer = lexers.Get(lexerAliases[lang])
	}
	if lexer == nil {
		return false
	}
	it, err := chroma.Coalesce(lexer).Tokenise(nil, code)
	if err != nil {
		return false
	}
	var buf bytes.Buffer
	if err := chromaFormatter.Format(&buf, codeStyle, it); err != nil {
		return false
	}
	_, _ = w.Write(buf.Bytes())
	return true
}

func codeRendererOption() renderer.Option {
	return renderer.WithNodeRenderers(util.Prioritized(codeRenderer{}, 200))
}

// ChromaCSS returns the stylesheet for the class-based highlighting output.
func ChromaCSS() (string, error) {
	var buf bytes.Buffer
	err := chromaFormatter.WriteCSS(&buf, codeStyle)
	return buf.String(), err
}
