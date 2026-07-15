# JsonForge

A human-friendly JSON CLI toolkit. Manipulate, query, validate, diff, and convert JSON with ease.

## Why JsonForge?

jq is powerful but notoriously hard to use. JsonForge provides a simpler, more intuitive interface for common JSON operations with clear error messages and a consistent command structure. Combines the best features of jq-style querying with human-friendly dot-notation paths.

## Features

- **pretty** / **minify** - Format or compact JSON with configurable indentation
- **get** / **set** / **del** - Dot-notation path operations (e.g., `users.0.name`, `users[*].age`)
- **query** - Full JSON query engine with path, filter, sort, limit, and output format
- **diff** - Semantic diff of two JSON files with change tracking
- **merge** - Deep merge multiple JSON files (last wins for conflicts)
- **stats** - Analyze JSON structure (types, depth, keys)
- **info** - Show JSON structure summary (type, depth, key counts)
- **flatten** / **unflatten** - Convert between nested and dot-notation
- **validate** - Validate against JSON Schema
- **convert** - Export to YAML, TOML, CSV, HTML

## Installation

```bash
go install github.com/EdgarOrtegaRamirez/jsonforge/cmd/jsonforge@latest
```

## Quick Start

```bash
# Pretty-print JSON
echo '{"name":"John","age":30}' | jsonforge pretty -

# Get a nested value
jsonforge get address.city data.json

# Query with path, filter, sort, and limit
jsonforge query --path users --filter 'age > 25' --sort age --limit 10 data.json

# Diff two files
jsonforge diff old.json new.json

# Convert to YAML
jsonforge convert -f yaml data.json

# Validate against schema
jsonforge validate schema.json data.json

# Flatten nested JSON
echo '{"a":{"b":1}}' | jsonforge flatten -

# Show structure info
jsonforge info data.json
```

## Query Command

The **query** command is the powerhouse — combining path-based extraction, filtering, sorting, limiting, and format control:

```bash
# Extract all names from users array
jsonforge query --path 'users[*].name' data.json

# Filter array items by condition
jsonforge query --path users --filter 'age > 25' data.json

# Combined: path + filter + sort + limit + format
jsonforge query --path users --filter 'age > 25' --sort age --sort-desc --limit 5 --format jsonl data.json

# Filter with string methods
jsonforge query --path users --filter 'name contains "li"' data.json

# Chained path: get nested array, then extract field
jsonforge query --path 'users[*].address.city' data.json
```

### Query Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--path` | `-p` | JSONPath-like query (e.g., `users[*].name`) |
| `--filter` | `-f` | Filter expression (e.g., `age > 25`, `name contains "li"`) |
| `--sort` | `-s` | Sort by field name |
| `--sort-desc` | `-S` | Sort descending |
| `--limit` | `-l` | Limit number of results |
| `--format` | `-F` | Output format: json, jsonl, text (default: json) |
| `--compact` | `-c` | Compact output (no pretty print) |

### Path Syntax

JsonForge uses intuitive dot-notation paths with bracket and wildcard support:

```bash
# Simple key
jsonforge query --path name data.json

# Nested key
jsonforge query --path address.city data.json

# Array index
jsonforge query --path tags.0 data.json

# Last array element
jsonforge query --path tags.last data.json

# Wildcard — all items in array
jsonforge query --path 'users[*].name' data.json
jsonforge query --path 'users.*.name' data.json

# Recursive descent — find anywhere in tree
jsonforge query --path '..email' data.json

# Bracket notation — keys with special characters
jsonforge query --path '[key-with-dash]' data.json

# Combined wildcard + nested
jsonforge query --path 'users[*].address.city' data.json
```

### Filter Expression Language

Supports comparison operators, logical operators, and string methods:

```bash
# Comparisons: ==, !=, >, >=, <, <=
jsonforge query --filter 'age > 25' data.json

# Logical operators: and, or, &&, ||, not
jsonforge query --filter 'age > 18 and age < 65' data.json
jsonforge query --filter 'name startsWith "A" or name startsWith "B"' data.json
jsonforge query --filter 'not (status == "deleted")' data.json

# String methods: contains, startsWith, endsWith, matches
jsonforge query --filter 'name contains "li"' data.json
jsonforge query --filter 'email matches "@example\.com$"' data.json

# Grouping with parentheses
jsonforge query --filter '(age > 30 or role == "admin") and active == true' data.json
```

## Commands

| Command | Description |
|---------|-------------|
| `pretty [file]` | Pretty-print JSON with indentation |
| `minify [file]` | Remove whitespace from JSON |
| `get <path> [file]` | Get value by dot-notation path |
| `set <path> <value> [file]` | Set value by dot-notation path |
| `del <path> [file]` | Delete value by path |
| `query [file]` | Query JSON with path, filter, sort, limit |
| `diff <old> <new>` | Semantic JSON diff |
| `merge <files...>` | Deep merge JSON files |
| `stats [file]` | Show JSON statistics |
| `info [file]` | Show JSON structure summary |
| `flatten [file]` | Flatten to dot-notation keys |
| `unflatten [file]` | Restore nested structure |
| `validate <schema> <data>` | Validate against JSON Schema |
| `convert [file]` | Convert to YAML/TOML/CSV/HTML |
| `version` | Print version |

## Architecture

```
cmd/jsonforge/      CLI entry point
pkg/query/          Dot-notation path engine (Get, Set, Delete, Sort)
pkg/filter/         Expression parser (recursive descent, supports &&/||/not/methods)
pkg/differ/         Semantic JSON diff engine with change tracking
pkg/stats/          JSON statistics analysis
pkg/info/           Structure summary analysis
pkg/flatten/        Flatten/unflatten conversion
pkg/validator/      JSON Schema validation
pkg/converter/      Format conversion (YAML, TOML, CSV, HTML)
pkg/output/         Output formatters (JSON, JSONL, Text)
```

## Development

```bash
# Build
go build -o jsonforge ./cmd/jsonforge/

# Run all tests
go test ./...

# Test with verbose output
go test -v ./...

# Run vet
go vet ./...
```

## License

MIT
