package main

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"notes/internal/api"
)

func newOpenAPICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI 3.1 document of the note API",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			doc, err := openAPIDocument()
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(doc)
			return err
		},
	}
}

func openAPIDocument() ([]byte, error) {
	doc, err := json.MarshalIndent(api.NewSpecAPI().OpenAPI(), "", "  ")
	if err != nil {
		return nil, err
	}
	return append(doc, '\n'), nil
}
