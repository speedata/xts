---
weight: 10
type: docs
linktitle: Command Line
---

# Running XTS on the Command Line

```
xts <command> <parameters>
```

## Commands

| Command | Description |
|---------|-------------|
| `run` | Load layout and data, create PDF (default -- you can omit it) |
| `clean` | Remove auxiliary and protocol files |
| `compare` | Compare PDFs for [quality assurance](../quality-assurance) |
| `doc` | Open the documentation in a web browser |
| `list-fonts` | List available fonts with `@font-face` rules |
| `new` | Create a starter project in the given directory |
| `version` | Print version number |

## Parameters

| Flag | Description |
|------|-------------|
| `-c, --config=NAME` | Config file name (default: `xts.cfg`) |
| `--data=NAME` | Data file (default: `data.xml`) |
| `--dummy` | Use `<data/>` as input instead of a data file |
| `--dumpoutput=FILE` | Write complete XML dump of generated PDF |
| `--extradir=DIR` | Additional directory for [file search](../file-organization) |
| `--filter=NAME` | Run a [Lua script](../../advanced/lua-filter) before publishing |
| `--jobname=NAME` | Output PDF name without extension (default: `xts`) |
| `--layout=NAME` | Layout file (default: `layout.xml`) |
| `--loglevel=LVL` | Console log level: `debug`, `info`, `warn`, `error` |
| `--runs=N` | Run XTS N times (for cross-references) |
| `--quiet` | No console output |
| `--suppressinfo` | Create reproducible PDF (no timestamps/random IDs) |
| `--systemfonts` | Include system-installed fonts |
| `--trace=NAMES` | Enable traces: `grid`, `gridallocation` |
| `--verbose` | Extra logging output |

## Examples

```bash
# Create a PDF with defaults
xts

# Use custom file names
xts --data products.xml --layout catalog.xml --jobname catalog

# Multiple runs for cross-references
xts --runs 2

# Quick test without a data file
xts --dummy

# Show the grid for debugging
xts --trace grid
```

For the full parameter reference, see [CLI Reference](/reference/cli).
