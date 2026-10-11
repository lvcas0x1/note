// Package webui embeds the built web app (apps/web, built into static/).
package webui

import (
	"embed"
	"io/fs"
)

//go:embed all:static
var static embed.FS

// FS returns the built web app files.
func FS() fs.FS {
	sub, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
