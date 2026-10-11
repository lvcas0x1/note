// Command notes builds, serves, and exports the note site.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "notes",
		Short:         "Build, serve, and export the note site",
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	root.AddCommand(
		newServeCmd(),
		newBuildCmd(),
		newMigrateSidebarCmd(),
		newStampDatesCmd(),
		newOpenAPICmd(),
		newExportCmd(),
	)
	return root
}
