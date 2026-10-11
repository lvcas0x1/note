package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestOpenAPIDocumentIsValid(t *testing.T) {
	raw, err := openAPIDocument()
	if err != nil {
		t.Fatal(err)
	}
	doc, err := openapi3.NewLoader().LoadFromData(raw)
	if err != nil {
		t.Fatal(err)
	}
	if doc.OpenAPI != "3.1.0" {
		t.Errorf("openapi = %q, want 3.1.0", doc.OpenAPI)
	}
	if err := doc.Validate(context.Background()); err != nil {
		t.Fatalf("invalid OpenAPI document: %v", err)
	}
	for _, p := range []string{"/api/v1/sidebar", "/api/v1/articles/{id}", "/api/v1/calendar", "/api/v1/days/{date}", "/api/v1/search"} {
		if doc.Paths.Find(p) == nil {
			t.Errorf("path %s missing", p)
		}
	}
	if doc.Components != nil && len(doc.Components.SecuritySchemes) > 0 {
		t.Error("the API must declare no security scheme")
	}
}

func TestCommittedOpenAPIIsCurrent(t *testing.T) {
	committed, err := os.ReadFile(filepath.Join("..", "..", "..", "openapi", "openapi.json"))
	if err != nil {
		t.Skip("openapi/openapi.json not generated yet")
	}
	raw, _ := openAPIDocument()
	if string(raw) != string(committed) {
		t.Fatal("openapi/openapi.json is out of date; run `npm run generate:openapi`")
	}
}
