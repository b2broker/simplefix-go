# fixgen Code Generator

## Overview

CLI tool that reads FIX XML schemas and generates strongly-typed Go code: message structs, builder interfaces, field tag constants, and enum values.

## Usage

```bash
fixgen -o=./fix44 -s=./source/fix44.xml -t=./source/types.xml
```

| Flag | Purpose |
|------|---------|
| `-o` | Output directory for generated Go files |
| `-s` | FIX XML schema path (e.g., `source/fix44.xml`) |
| `-t` | Types mapping XML path (maps FIX types → Go types) |

Source: [cmd/fixgen/main.go](../../cmd/fixgen/main.go)

## Input Files

### Schema XML (`fix44.xml`)

Defines the FIX dialect:
- **Messages** — message types with their fields and groups
- **Fields** — tag numbers, names, FIX types
- **Components** — reusable field groups (Header, Trailer, etc.)
- **Groups** — repeating group definitions

### Types XML (`types.xml`)

Maps FIX protocol types to Go types:
- STRING → string
- INT → int
- FLOAT → float64
- etc.

## Generated Output

| Output | Content |
|--------|---------|
| Message structs | Typed message with all fields as properties |
| Builder interfaces | HeaderBuilder, LogonBuilder, etc. |
| Field constants | `FieldMsgType = "35"`, `FieldMsgSeqNum = "34"`, etc. |
| Enum constants | `EnumEncryptMethodNoneother`, `EnumSessionRejectReasonInvalidtagnumber`, etc. |

## Implementation

Generator pipeline in [generator/generator.go](../../generator/generator.go):
1. Parse XML with `utils.ParseXML()`
2. Build data structures: `Doc`, `Field`, `Component`, `Message`
3. Generate Go code using `text/template` templates
4. Auto-format with `gofmt`

## Example Generated Code

See [tests/fix44/](../../tests/fix44/) — generated from `source/fix44.xml` for integration testing.

## Warning

Generated files are overwritten on each run. Never edit files in the output directory manually.
