package content

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// SidebarFile is the name of the sidebar definition at the content root.
const SidebarFile = "sidebar.json"

// Node is one entry of sidebar.json. A node has a Path, Children, or both.
type Node struct {
	Title    string `json:"title"`
	Path     string `json:"path,omitempty"`
	Children []Node `json:"children,omitempty"`
}

// LoadSidebar reads and validates sidebar.json in root. It checks the schema,
// non-blank titles, and that every path points to an existing file inside root.
func LoadSidebar(root string) ([]Node, error) {
	raw, err := os.ReadFile(filepath.Join(root, SidebarFile))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%s not found in content root %s", SidebarFile, root)
	}
	if err != nil {
		return nil, err
	}
	sch, _, err := schemas()
	if err != nil {
		return nil, fmt.Errorf("load sidebar schema: %w", err)
	}
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("%s: invalid JSON: %w", SidebarFile, err)
	}
	if err := sch.Validate(inst); err != nil {
		return nil, fmt.Errorf("%s: invalid node at %s", SidebarFile, describeValidation(err))
	}
	var nodes []Node
	if err := json.Unmarshal(raw, &nodes); err != nil {
		return nil, fmt.Errorf("%s: %w", SidebarFile, err)
	}
	if err := checkNodes(root, nodes, ""); err != nil {
		return nil, err
	}
	return nodes, nil
}

func checkNodes(root string, nodes []Node, pointer string) error {
	for i := range nodes {
		n := &nodes[i]
		p := fmt.Sprintf("%s/%d", pointer, i)
		n.Title = strings.TrimSpace(n.Title)
		if n.Title == "" {
			return fmt.Errorf("%s: invalid node at %s: title is blank", SidebarFile, p)
		}
		if n.Path != "" {
			if _, err := ResolveInside(root, n.Path); err != nil {
				return fmt.Errorf("%s: node %q at %s: %w", SidebarFile, n.Title, p, err)
			}
		}
		if err := checkNodes(root, n.Children, p+"/children"); err != nil {
			return err
		}
	}
	return nil
}

// ResolveInside returns the absolute path of rel under root after resolving
// symlinks. It fails when the file does not exist or resolves outside root.
func ResolveInside(root, rel string) (string, error) {
	if filepath.IsAbs(rel) {
		return "", fmt.Errorf("path %q is outside the content folder", rel)
	}
	absRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	absRoot, err = filepath.Abs(absRoot)
	if err != nil {
		return "", err
	}
	joined := filepath.Join(absRoot, filepath.FromSlash(rel))
	if !within(absRoot, joined) {
		return "", fmt.Errorf("path %q is outside the content folder", rel)
	}
	resolved, err := filepath.EvalSymlinks(joined)
	if errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("file %q does not exist", rel)
	}
	if err != nil {
		return "", err
	}
	if !within(absRoot, resolved) {
		return "", fmt.Errorf("path %q is outside the content folder", rel)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("path %q is a directory, not a file", rel)
	}
	return resolved, nil
}

func within(root, p string) bool {
	rel, err := filepath.Rel(root, p)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}
