---
type: docs
linktitle: CLI Reference
---

# CLI Reference

Complete specification of the `xts` command-line interface.

## Synopsis

```
xts [command] [flags]
```

If no command is given, `run` is assumed.

## Commands

| Command | Description |
|---------|-------------|
| `run` | Read layout and data files, produce PDF. This is the default. |
| `clean` | Remove auxiliary files (`xts-protocol.xml`, `xts-aux.xml`, etc.) |
| `compare <dir>` | Recursively compare generated PDFs against `reference.pdf` files |
| `doc` | Open the documentation website in the default browser |
| `list-fonts` | Print `@font-face` CSS rules for all fonts found in the search path |
| `new [dir]` | Create a starter project with `data.xml` and `layout.xml` |
| `version` | Print the version number and exit |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-c`, `--config` | string | `xts.cfg` | Configuration file to read |
| `--data` | string | `data.xml` | Data file name |
| `--dummy` | boolean | `false` | Ignore data file, use `<data/>` |
| `--dumpoutput` | string | | Write XML dump of PDF structure to this file |
| `--extradir` | string | | Additional directory for file lookups (recursive) |
| `--filter` | string | | Lua script to run before publishing |
| `--jobname` | string | `xts` | Output file name (without `.pdf`) |
| `--layout` | string | `layout.xml` | Layout file name |
| `--loglevel` | string | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `--runs` | integer | `1` | Number of publishing runs |
| `--quiet` | boolean | `false` | Suppress all console output |
| `--suppressinfo` | boolean | `false` | Produce reproducible PDF (no timestamps) |
| `--systemfonts` | boolean | `false` | Include system-installed fonts in search |
| `--trace` | string | | Comma-separated traces: `grid`, `gridallocation` |
| `--verbose` | boolean | `false` | Extra debug output |

## Examples

```bash
# Default: read data.xml + layout.xml, produce xts.pdf
xts

# Custom files and output name
xts --data products.xml --layout catalog.xml --jobname catalog

# Multiple runs (for cross-references, page counts)
xts --runs 2

# Quick test without data file
xts --dummy

# Debug grid placement
xts --trace grid,gridallocation

# Create reference PDF for QA
xts --suppressinfo --jobname reference

# Run QA suite
xts compare qa/

# Generate font CSS
xts list-fonts

# Scaffold a new project
xts new myproject
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| Non-zero | Error (check console output or protocol file) |

## Output files

| File | Description |
|------|-------------|
| `<jobname>.pdf` | The generated PDF |
| `xts-protocol.xml` | Processing protocol with messages, warnings, errors |
| `xts-aux.xml` | Auxiliary data (marks, page numbers) for subsequent runs |
