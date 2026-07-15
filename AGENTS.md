# AGENTS.md

## Project Overview

JsonForge is a human-friendly JSON CLI toolkit written in Go. It provides intuitive commands for manipulating, querying, validating, diffing, and converting JSON data.

## Architecture

- **cmd/jsonforge/**: CLI entry point using Cobra framework
- **pkg/query/**: Dot-notation path engine with Get/Set/Delete/Wildcard operations
- **pkg/filter/**: Expression parser — recursive descent with &&/||/not, string methods (contains/startsWith/endsWith/matches)
- **pkg/differ/**: Semantic JSON diff engine with change tracking
- **pkg/stats/**: JSON statistics analysis (types, depth, keys)
- **pkg/flatten/**: Flatten/unflatten between nested and dot-notation
- **pkg/validator/**: JSON Schema validation (type, required, min/max)
- **pkg/converter/**: Format conversion (YAML, TOML, CSV, HTML)
- **pkg/output/**: Output formatters (JSON, JSONL, Text)

## Development Commands

```bash
# Build
go build -o jsonforge ./cmd/jsonforge/

# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/query/...
go test ./pkg/filter/...
go test ./pkg/differ/...

# Test with verbose output
go test -v ./...

# Run vet
go vet ./...
```

## Key Design Decisions

1. **Dot-notation paths**: Intuitive path syntax with wildcard support (`*`, `[*]`, `..`)
2. **Semantic diffing**: Compares JSON structure, not text lines
3. **Query command**: Combines path, filter, sort, limit, and format in one command
4. **No external dependencies for core**: Only Cobra for CLI framework
5. **Single binary**: Cross-platform, no runtime dependencies

## Adding New Commands

1. Create a `newXxxCmd()` function in a separate file in `cmd/jsonforge/`
2. Add to `root.AddCommand()` list in `main.go`
3. Add tests for the command logic

## Testing Strategy

- Unit tests for each package (query, differ, stats, flatten, validator, converter)
- Integration tests via CLI commands
- Edge cases: empty input, nested structures, arrays, null values
