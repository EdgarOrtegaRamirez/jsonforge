package main

import (
	"fmt"

	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/output"
	"github.com/spf13/cobra"
)

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [file]",
		Short: "Show JSON structure summary",
		Long: `Display a summary of JSON structure: types, depth, key counts, and array lengths.

Examples:
  jsonforge info data.json
  cat data.json | jsonforge info -
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var data interface{}
			var err error

			if len(args) == 0 || args[0] == "-" {
				data, err = readJSON("-")
			} else {
				data, err = readJSON(args[0])
			}
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}

			w := output.NewInfoWriter()
			if err := w.Write(data, cmd.OutOrStdout()); err != nil {
				return fmt.Errorf("info: %w", err)
			}
			return nil
		},
	}
}
