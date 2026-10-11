package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"notes/internal/content"
)

func newStampDatesCmd() *cobra.Command {
	var contentDir string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "stamp-dates",
		Short: "Write date: YYYY-MM-DD into articles that have no date (creation time, else modification time)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			stamps, err := content.StampDates(contentDir, dryRun)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			for _, s := range stamps {
				_, _ = fmt.Fprintf(out, "%s  %s  (from %s time)\n", s.Date, s.Path, s.Source)
			}
			verb := "stamped"
			if dryRun {
				verb = "would stamp"
			}
			_, _ = fmt.Fprintf(out, "%s %d articles\n", verb, len(stamps))
			return nil
		},
	}
	cmd.Flags().StringVar(&contentDir, "content", envOr("NOTES_CONTENT", "content"), "content folder (env NOTES_CONTENT)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "list the dates without writing files")
	return cmd
}
