package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/converter"
	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/differ"
	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/flatten"
	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/query"
	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/stats"
	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/validator"
	"github.com/spf13/cobra"
)

var version = "2.0.0"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "jsonforge",
		Short: "A human-friendly JSON toolkit",
		Long:  "JsonForge - manipulate, query, validate, diff, and convert JSON with ease.",
	}

	root.AddCommand(
		newPrettyCmd(),
		newMinifyCmd(),
		newGetCmd(),
		newSetCmd(),
		newDelCmd(),
		newDiffCmd(),
		newMergeCmd(),
		newStatsCmd(),
		newFlattenCmd(),
		newUnflattenCmd(),
		newValidateCmd(),
		newConvertCmd(),
		newQueryCmd(),
		newInfoCmd(),
		newVersionCmd(),
	)

	return root
}

func readJSON(path string) (interface{}, error) {
	var data []byte
	var err error
	if path == "-" || path == "" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}
	return query.ParseJSON(string(data))
}

func readTwoJSON(args []string) (interface{}, interface{}, error) {
	if len(args) < 2 {
		return nil, nil, fmt.Errorf("requires 2 arguments")
	}
	old, err := readJSON(args[0])
	if err != nil {
		return nil, nil, fmt.Errorf("reading %s: %w", args[0], err)
	}
	new, err := readJSON(args[1])
	if err != nil {
		return nil, nil, fmt.Errorf("reading %s: %w", args[1], err)
	}
	return old, new, nil
}

func marshalJSON(data interface{}, indent int) (string, error) {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	if indent > 0 {
		enc.SetIndent("", strings.Repeat(" ", indent))
	}
	if err := enc.Encode(data); err != nil {
		return "", err
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

func deepMerge(a, b interface{}) interface{} {
	aMap, aOk := a.(map[string]interface{})
	bMap, bOk := b.(map[string]interface{})

	if aOk && bOk {
		result := make(map[string]interface{})
		for k, v := range aMap {
			result[k] = v
		}
		for k, v := range bMap {
			if existing, ok := result[k]; ok {
				result[k] = deepMerge(existing, v)
			} else {
				result[k] = v
			}
		}
		return result
	}

	return b
}

func newPrettyCmd() *cobra.Command {
	var indent int
	cmd := &cobra.Command{
		Use:   "pretty [file]",
		Short: "Pretty-print JSON with indentation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := readJSON(args[0])
			if err != nil {
				return err
			}
			output, _ := marshalJSON(data, indent)
			fmt.Fprintln(cmd.OutOrStdout(), output)
			return nil
		},
	}
	cmd.Flags().IntVarP(&indent, "indent", "i", 2, "indentation spaces")
	return cmd
}

func newMinifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "minify [file]",
		Short: "Minify JSON (remove whitespace)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := readJSON(args[0])
			if err != nil {
				return err
			}
			output, _ := marshalJSON(data, 0)
			fmt.Fprintln(cmd.OutOrStdout(), output)
			return nil
		},
	}
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <path> [file]",
		Short: "Get a value by dot-notation path",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			file := "-"
			if len(args) > 1 {
				file = args[1]
			}
			data, err := readJSON(file)
			if err != nil {
				return err
			}
			val, err := query.Get(data, path)
			if err != nil {
				return err
			}
			output, _ := marshalJSON(val, 2)
			fmt.Fprintln(cmd.OutOrStdout(), output)
			return nil
		},
	}
}

func newSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <path> <value> [file]",
		Short: "Set a value by dot-notation path",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			valueStr := args[1]
			file := "-"
			if len(args) > 2 {
				file = args[2]
			}

			data, err := readJSON(file)
			if err != nil {
				return err
			}

			value, err := query.ParseJSON(valueStr)
			if err != nil {
				value = valueStr
			}

			result, err := query.Set(data, path, value)
			if err != nil {
				return err
			}
			output, _ := marshalJSON(result, 2)
			fmt.Fprintln(cmd.OutOrStdout(), output)
			return nil
		},
	}
}

func newDelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "del <path> [file]",
		Short: "Delete a value by dot-notation path",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			file := "-"
			if len(args) > 1 {
				file = args[1]
			}
			data, err := readJSON(file)
			if err != nil {
				return err
			}
			result, err := query.Delete(data, path)
			if err != nil {
				return err
			}
			output, _ := marshalJSON(result, 2)
			fmt.Fprintln(cmd.OutOrStdout(), output)
			return nil
		},
	}
}

func newDiffCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "diff <old.json> <new.json>",
		Short: "Semantic diff of two JSON files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			old, new, err := readTwoJSON(args)
			if err != nil {
				return err
			}
			result := differ.Diff(old, new)
			switch format {
			case "json":
				fmt.Fprintln(cmd.OutOrStdout(), differ.FormatJSON(result))
			default:
				fmt.Fprint(cmd.OutOrStdout(), differ.FormatText(result))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json)")
	return cmd
}

func newMergeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "merge <file1> [file2] ...",
		Short: "Deep merge multiple JSON files (last wins)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var result interface{}
			for _, arg := range args {
				data, err := readJSON(arg)
				if err != nil {
					return err
				}
				if result == nil {
					result = data
					continue
				}
				result = deepMerge(result, data)
			}
			output, _ := marshalJSON(result, 2)
			fmt.Fprintln(cmd.OutOrStdout(), output)
			return nil
		},
	}
}

func newStatsCmd() *cobra.Command {
	var includePaths bool
	cmd := &cobra.Command{
		Use:   "stats [file]",
		Short: "Show JSON statistics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := readJSON(args[0])
			if err != nil {
				return err
			}
			s := stats.Analyze(data, includePaths)
			fmt.Fprint(cmd.OutOrStdout(), stats.FormatText(s))
			return nil
		},
	}
	cmd.Flags().BoolVar(&includePaths, "paths", false, "include all leaf paths")
	return cmd
}

func newFlattenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "flatten [file]",
		Short: "Flatten nested JSON to dot-notation keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := readJSON(args[0])
			if err != nil {
				return err
			}
			flat := flatten.Flatten(data)
			keys := flatten.Keys(flat)
			for _, k := range keys {
				output, _ := marshalJSON(flat[k], 0)
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", k, output)
			}
			return nil
		},
	}
}

func newUnflattenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unflatten [file]",
		Short: "Unflatten dot-notation keys back to nested JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := readJSON(args[0])
			if err != nil {
				return err
			}
			flat, ok := data.(map[string]interface{})
			if !ok {
				return fmt.Errorf("input must be a JSON object for unflatten")
			}
			result := flatten.Unflatten(flat)
			output, _ := marshalJSON(result, 2)
			fmt.Fprintln(cmd.OutOrStdout(), output)
			return nil
		},
	}
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <schema.json> <data.json>",
		Short: "Validate JSON data against a JSON Schema",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaData, err := readJSON(args[0])
			if err != nil {
				return err
			}
			schema, ok := schemaData.(map[string]interface{})
			if !ok {
				return fmt.Errorf("schema must be a JSON object")
			}

			data, err := readJSON(args[1])
			if err != nil {
				return err
			}

			result := validator.Validate(data, schema)
			fmt.Fprint(cmd.OutOrStdout(), validator.FormatValidation(result))
			if !result.Valid {
				os.Exit(1)
			}
			return nil
		},
	}
}

func newConvertCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "convert [file]",
		Short: "Convert JSON to other formats",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := readJSON(args[0])
			if err != nil {
				return err
			}

			switch format {
			case "yaml":
				fmt.Fprint(cmd.OutOrStdout(), converter.ToYAML(data))
			case "toml":
				fmt.Fprint(cmd.OutOrStdout(), converter.ToTOML(data))
			case "csv":
				csv, err := converter.ToCSV(data)
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), csv)
			case "html":
				html, err := converter.ToHTMLTable(data)
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), html)
			default:
				return fmt.Errorf("unsupported format: %s (use yaml, toml, csv, html)", format)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "output format (yaml, toml, csv, html)")
	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "jsonforge v%s\n", version)
		},
	}
}
