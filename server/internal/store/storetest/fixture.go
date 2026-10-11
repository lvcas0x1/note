// Package storetest generates deterministic build fixtures for tests.
package storetest

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"notes/internal/content"
)

var words = strings.Fields(`kafka cluster bucket lambda network gateway policy role token schema index
query cache queue stream topic broker deploy pipeline container image volume secret certificate
golang python typescript react vite tailwind storybook sqlite fts trigram calendar sidebar`)

// Articles returns n published articles with reproducible titles, dates, and bodies.
// Dates cluster on a few days so calendar counts are greater than one.
func Articles(seed uint64, n int) []*content.Article {
	r := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15)) //nolint:gosec // deterministic test data, not security
	days := []string{"2026-10-01", "2026-10-05", "2026-10-11", "2026-09-30", "2025-03-14"}
	out := make([]*content.Article, 0, n)
	for i := range n {
		title := fmt.Sprintf("%s %s %03d", strings.ToUpper(words[r.IntN(len(words))][:1])+words[r.IntN(len(words))][1:], words[r.IntN(len(words))], i)
		var body []string
		for range 80 + r.IntN(120) {
			body = append(body, words[r.IntN(len(words))])
		}
		out = append(out, &content.Article{
			ID: fmt.Sprintf("gen/a%03d", i), Path: fmt.Sprintf("gen/a%03d.md", i), Title: title,
			Date: days[r.IntN(len(days))], Summary: strings.Join(body[:12], " "), Tags: []string{words[r.IntN(len(words))]},
			HTML: "<p>" + strings.Join(body, " ") + "</p>", Text: title + " " + strings.Join(body, " "),
		})
	}
	return out
}

// Result wraps articles in a build result with a sidebar that lists the first few.
func Result(articles []*content.Article) *content.Result {
	var children []content.Node
	for _, a := range articles[:min(3, len(articles))] {
		children = append(children, content.Node{Title: a.Title, Path: a.Path})
	}
	return &content.Result{
		Sidebar:  []content.Node{{Title: "Generated", Children: []content.Node{{Title: "Nested", Children: children}}}},
		Articles: articles,
	}
}
