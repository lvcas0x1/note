package content

import (
	"bytes"
	"embed"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

// schemaFS holds the JSON Schemas generated from packages/schemas.
//
//go:embed schemas/*.json
var schemaFS embed.FS

var (
	schemasOnce sync.Once
	schemasErr  error
	sidebarSch  *jsonschema.Schema
	frontSch    *jsonschema.Schema
)

func schemas() (sidebar, frontMatter *jsonschema.Schema, err error) {
	schemasOnce.Do(func() {
		c := jsonschema.NewCompiler()
		c.AssertFormat()
		for _, name := range []string{"sidebar.schema.json", "frontmatter.schema.json"} {
			raw, err := schemaFS.ReadFile("schemas/" + name)
			if err != nil {
				schemasErr = err
				return
			}
			doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
			if err != nil {
				schemasErr = fmt.Errorf("%s: %w", name, err)
				return
			}
			if err := c.AddResource(name, doc); err != nil {
				schemasErr = err
				return
			}
		}
		if sidebarSch, schemasErr = c.Compile("sidebar.schema.json"); schemasErr != nil {
			return
		}
		frontSch, schemasErr = c.Compile("frontmatter.schema.json")
	})
	return sidebarSch, frontSch, schemasErr
}

// describeValidation turns a schema validation error into "<pointer>: <reason>" lines.
// It keeps the deepest failing leaves, because those name the offending node
// even when the failure surfaces through an anyOf.
func describeValidation(err error) string {
	ve, ok := err.(*jsonschema.ValidationError) //nolint:errorlint // the validator returns this concrete type
	if !ok {
		return err.Error()
	}
	var leaves []*jsonschema.ValidationError
	var walk func(e *jsonschema.ValidationError)
	walk = func(e *jsonschema.ValidationError) {
		if len(e.Causes) == 0 {
			leaves = append(leaves, e)
			return
		}
		for _, c := range e.Causes {
			walk(c)
		}
	}
	walk(ve)
	deepest := 0
	for _, l := range leaves {
		deepest = max(deepest, len(l.InstanceLocation))
	}
	seen := map[string]bool{}
	var lines []string
	for _, l := range leaves {
		if len(l.InstanceLocation) != deepest {
			continue
		}
		line := "/" + strings.Join(l.InstanceLocation, "/") + ": " + l.ErrorKind.LocalizedString(printer)
		if !seen[line] {
			seen[line] = true
			lines = append(lines, line)
		}
	}
	sort.Strings(lines)
	return strings.Join(lines, "; ")
}
