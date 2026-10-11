package content

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func relLuminance(hex string) float64 {
	c := [3]float64{}
	for i := range 3 {
		v, _ := strconv.ParseUint(hex[1+2*i:3+2*i], 16, 8)
		x := float64(v) / 255
		if x <= 0.04045 {
			c[i] = x / 12.92
		} else {
			c[i] = math.Pow((x+0.055)/1.055, 2.4)
		}
	}
	return 0.2126*c[0] + 0.7152*c[1] + 0.0722*c[2]
}

func contrastRatio(a, b string) float64 {
	la, lb := relLuminance(a), relLuminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

// Every code color must reach WCAG AA (4.5:1) on the code background.
func TestCodeColorsContrast(t *testing.T) {
	for _, tt := range codeStyle.Types() {
		e := codeStyle.Get(tt)
		if !e.Colour.IsSet() {
			continue
		}
		fg := e.Colour.String()
		if r := contrastRatio(fg, codeBackground); r < 4.5 {
			t.Errorf("%s %s: contrast %.2f on %s", tt, fg, r, codeBackground)
		}
	}
}

func TestCurlFenceIsHighlightedAsBash(t *testing.T) {
	a, err := compileOne(t, map[string]string{"c.md": "---\ndate: 2026-01-01\n---\n```curl\ncurl -s https://example.com | jq .\n```\n"}, "c.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(a.HTML, `class="language-curl"`) || !strings.Contains(a.HTML, `<span class="`) {
		t.Fatalf("curl block not highlighted: %s", a.HTML)
	}
}
