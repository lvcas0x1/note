package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"notes/internal/content"
	"notes/internal/store"
)

func newBuildCmd() *cobra.Command {
	var contentDir, dataDir string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Validate and compile the content folder into a new build",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := os.MkdirAll(dataDir, 0o750); err != nil {
				return err
			}
			tmp, err := os.MkdirTemp(dataDir, ".fetch-")
			if err != nil {
				return err
			}
			defer func() { _ = os.RemoveAll(tmp) }()
			res, err := content.Build(contentDir, content.Options{Fetcher: content.NewFetcher(tmp)})
			if err != nil {
				return err
			}
			for _, w := range res.Warnings {
				cmd.PrintErrln("warning:", w)
			}
			dir, err := store.WriteBuild(dataDir, res)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "built %d articles and %d assets into %s\n", len(res.Articles), len(res.Assets), filepath.Base(dir))
			return nil
		},
	}
	cmd.Flags().StringVar(&contentDir, "content", envOr("NOTES_CONTENT", "content"), "content folder (env NOTES_CONTENT)")
	cmd.Flags().StringVar(&dataDir, "data", envOr("NOTES_DATA", "data"), "data folder for builds (env NOTES_DATA)")
	return cmd
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
