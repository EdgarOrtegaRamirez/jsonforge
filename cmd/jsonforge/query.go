package main

import (
	"fmt"

	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/filter"
	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/output"
	"github.com/EdgarOrtegaRamirez/jsonforge/pkg/query"
	"github.com/spf13/cobra"
)

var (
	queryPath     string
	queryFilter   string
	querySort     string
	querySortDesc bool
	queryLimit    int
	queryFormat   string
	queryCompact  bool
)

func newQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query [file]",
		Short: "Query JSON with path, filter, sort, and limit",
		Long: `Extract, filter, sort, and format data from JSON files.

Path syntax supports:
  - "key" — top-level key lookup
  - "key.subkey" — nested key lookup
  - "key[*]" — array of objects at key
  - "key[*].subkey" — extract subkey from array of objects
  - "*key" — find all keys containing 'key'

Filter expressions support:
  - Comparisons: ==, !=, >, >=, <, <=
  - Logic: and, or, &&, ||, not
  - String methods: contains, startsWith, endsWith, matches
  - Parentheses for grouping

Examples:
  jsonforge query --path users data.json
  jsonforge query --path 'users[*].name' data.json
  jsonforge query --path users --filter 'age > 25' data.json
  jsonforge query --path items --sort price --limit 5 data.json
  jsonforge query --path data --format jsonl data.json
  cat data.json | jsonforge query --path items
`,
		Args: cobra.MaximumNArgs(1),
		RunE: runQuery,
	}

	cmd.Flags().StringVarP(&queryPath, "path", "p", "", "JSONPath-like query (e.g., 'users[*].name')")
	cmd.Flags().StringVarP(&queryFilter, "filter", "f", "", "Filter expression (e.g., 'age > 25')")
	cmd.Flags().StringVarP(&querySort, "sort", "s", "", "Sort by field name")
	cmd.Flags().BoolVarP(&querySortDesc, "sort-desc", "S", false, "Sort descending")
	cmd.Flags().IntVarP(&queryLimit, "limit", "l", 0, "Limit number of results")
	cmd.Flags().StringVarP(&queryFormat, "format", "F", "json", "Output format: json, jsonl, text")
	cmd.Flags().BoolVarP(&queryCompact, "compact", "c", false, "Compact output (no pretty print)")

	return cmd
}

func runQuery(cmd *cobra.Command, args []string) error {
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

	// Apply query path
	if queryPath != "" {
		data, err = query.Get(data, queryPath)
		if err != nil {
			return fmt.Errorf("query: %w", err)
		}
	}

	// Apply filter if specified
	if queryFilter != "" {
		f, err := filter.NewParser()
		if err != nil {
			return fmt.Errorf("filter: %w", err)
		}
		data, err = f.Evaluate(data, queryFilter)
		if err != nil {
			return fmt.Errorf("filter: %w", err)
		}
	}

	// Apply sort if specified
	if querySort != "" {
		data, err = sortResults(data, querySort, querySortDesc)
		if err != nil {
			return fmt.Errorf("sort: %w", err)
		}
	}

	// Apply limit
	if queryLimit > 0 {
		data, err = limitResults(data, queryLimit)
		if err != nil {
			return fmt.Errorf("limit: %w", err)
		}
	}

	// Output
	outFormat := output.ParseFormat(queryFormat)

	var w output.Writer
	switch outFormat {
	case output.JSON:
		w = output.NewJSONWriter(queryCompact)
	case output.JSONL:
		w = output.NewJSONLWriter()
	case output.TEXT:
		w = output.NewTextWriter()
	default:
		w = output.NewJSONWriter(queryCompact)
	}

	if err := w.Write(data, cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("output: %w", err)
	}

	return nil
}

// sortResults sorts an array of objects by a field name.
func sortResults(data interface{}, field string, desc bool) (interface{}, error) {
	arr, ok := data.([]interface{})
	if !ok {
		return data, nil // not an array, return as-is
	}

	sorted := make([]interface{}, len(arr))
	copy(sorted, arr)

	query.SortSlice(sorted, field, desc)

	return sorted, nil
}

// limitResults limits the number of items in an array.
func limitResults(data interface{}, n int) (interface{}, error) {
	arr, ok := data.([]interface{})
	if !ok || n <= 0 {
		return data, nil
	}
	if n >= len(arr) {
		return data, nil
	}
	return arr[:n], nil
}
