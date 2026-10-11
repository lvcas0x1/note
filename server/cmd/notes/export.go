package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"notes/internal/export"
	"notes/internal/webui"
)

func newExportCmd() *cobra.Command {
	var o export.Options
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Write a static copy of the current build for static hosts such as GitHub Pages",
		Long: "Writes the site to --out and videos to <out>-videos. Upload the videos folder to the\n" +
			"location given by --video-base-url yourself; push the site folder to the Pages repository.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			o.Web = webui.FS()
			s, err := export.Run(cmd.Context(), o)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "exported %d articles and %d assets to %s (%.1f MB)\n", s.Articles, s.Assets, o.Out, float64(s.SiteBytes)/(1<<20))
			if s.Videos > 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%d videos are in %s; upload that folder so its files are served under %s\n",
					s.Videos, export.VideosDir(o.Out), o.VideoBaseURL)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&o.DataDir, "data", envOr("NOTES_DATA", "data"), "data folder with the builds (env NOTES_DATA)")
	cmd.Flags().StringVar(&o.Out, "out", envOr("NOTES_OUT", "static"), "output folder for the static site; must be empty or missing (env NOTES_OUT)")
	cmd.Flags().StringVar(&o.BasePath, "base-path", "/", "URL path the site is served under, for example /notes/")
	cmd.Flags().StringVar(&o.VideoBaseURL, "video-base-url", "", "absolute URL where the videos folder will be uploaded")
	return cmd
}
