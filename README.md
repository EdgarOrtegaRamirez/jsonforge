# JsonForge

A human-friendly JSON CLI toolkit. Manipulate, query, validate, diff, and convert JSON with ease.

## Why JsonForge?

jq is powerful but notoriously hard to use. JsonForge provides a simpler, more intuitive interface for common JSON operations with clear error messages and a consistent command structure.

## Features

- **pretty** / **minify** - Format or compact JSON
- **get** / **set** / **del** - Dot-notation path operations (e.g., `users.0.name`)
- **diff** - Semantic diff of two JSON files
- **merge** - Deep merge multiple JSON files
- **stats** - Analyze JSON structure (types, depth, keys)
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

# Diff two files
jsonforge diff old.json new.json

# Convert to YAML
jsonforge convert -f yaml data.json

# Validate against schema
jsonforge validate schema.json data.json

# Flatten nested JSON
echo '{"a":{"b":1}}' | jsonforge flatten -
```

## Commands

| Command | Description |
|---------|-------------|
| `pretty [file]` | Pretty-print JSON with indentation |
| `minify [file]` | Remove whitespace from JSON |
| `get <path> [file]` | Get value by dot-notation path |
| `set <path> <value> [file]` | Set value by dot-notation path |
| `del <path> [file]` | Delete value by path |
| `diff <old> <new>` | Semantic JSON diff |
| `merge <files...>` | Deep merge JSON files |
| `stats [file]` | Show JSON statistics |
| `flatten [file]` | Flatten to dot-notation keys |
| `unflatten [file]` | Restore nested structure |
| `validate <schema> <data>` | Validate against JSON Schema |
| `convert [file]` | Convert to YAML/TOML/CSV/HTML |
| `version` | Print version |

## Path Syntax

JsonForge uses intuitive dot-notation paths:

```bash
# Simple key
jsonforge get name data.json

# Nested key
jsonforge get address.city data.json

# Array index
jsonforge get tags.0 data.json

# Last array element
jsonforge get tags.last data.json

# Wildcard (all items)
jsonforge get users.*.name data.json

# Recursive descent
jsonforge get a..name data.json

# Bracket notation
jsonforge get '["key-with-dash"]' data.json
```

## Architecture

```
cmd/jsonforge/      CLI entry point
pkg/query/          Dot-notation path engine (Get, Set, Delete)
pkg/differ/         Semantic JSON diff engine
pkg/stats/          JSON statistics analysis
pkg/flatten/        Flatten/unflatten conversion
pkg/validator/      JSON Schema validation
pkg/converter/      Format conversion (YAML, TOML, CSV, HTML)
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
