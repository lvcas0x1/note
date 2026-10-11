package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"notes/internal/content"
)

func newMigrateSidebarCmd() *cobra.Command {
	var root, out string
	cmd := &cobra.Command{
		Use:   "migrate-sidebar <path/to/_sidebar.md>",
		Short: "Convert a docsify _sidebar.md into sidebar.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			if root == "" {
				root = filepath.Dir(args[0])
			}
			if out == "" {
				out = filepath.Join(root, content.SidebarFile)
			}
			nodes, warnings, err := content.MigrateDocsifySidebar(source, root)
			if err != nil {
				return err
			}
			for _, w := range warnings {
				cmd.PrintErrln("warning:", w)
			}
			raw, err := json.MarshalIndent(nodes, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(out, append(raw, '\n'), 0o600); err != nil {
				return err
			}
			if _, err := content.LoadSidebar(filepath.Dir(out)); err != nil && filepath.Dir(out) == root {
				return fmt.Errorf("wrote %s, but it does not validate: %w", out, err)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", out)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", "", "content folder that sidebar links are relative to (default: the folder of _sidebar.md)")
	cmd.Flags().StringVar(&out, "out", "", "output file (default: <root>/sidebar.json)")
	return cmd
}
