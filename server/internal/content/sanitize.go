package content

import (
	"regexp"
	"sync"

	"github.com/microcosm-cc/bluemonday"
)

var (
	policyOnce sync.Once
	policy     *bluemonday.Policy
)

// sanitizer extends bluemonday's UGC policy with what compiled articles need:
// embed placeholders, chroma and language classes, heading and footnote IDs,
// task-list checkboxes, and inline data: images.
func sanitizer() *bluemonday.Policy {
	policyOnce.Do(func() {
		p := bluemonday.UGCPolicy()
		p.AllowAttrs("data-embed").Matching(regexp.MustCompile(`^(drawio|excalidraw|openapi|mp4|mermaid|echarts)$`)).OnElements("div", "pre")
		p.AllowAttrs("data-src").Matching(regexp.MustCompile(`^/assets/[0-9a-f]{16}/[^"<>]+$`)).OnElements("div")
		p.AllowAttrs("class").Matching(regexp.MustCompile(`^(chroma|language-[\w+#.-]+|[a-z][a-z0-9]{0,3}|line|cl|footnotes?|footnote-ref|footnote-backref|task-list-item)$`)).OnElements("pre", "code", "span", "a", "sup", "div", "li", "ul", "hr")
		p.AllowAttrs("id").Matching(regexp.MustCompile(`^[\w\p{L}-]+$|^fn(ref)?:\d+$`)).OnElements("h1", "h2", "h3", "h4", "h5", "h6", "li", "sup")
		p.AllowAttrs("role").Matching(regexp.MustCompile(`^doc-(noteref|endnotes|backlink)$`)).OnElements("a", "div", "sup")
		p.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
		p.AllowAttrs("checked", "disabled").OnElements("input")
		p.AllowStyles("color").OnElements("span")
		p.AllowDataURIImages()
		policy = p
	})
	return policy
}

// sanitizeHTML removes scripts, event handlers, and anything else outside the policy.
func sanitizeHTML(html string) string { return sanitizer().Sanitize(html) }
