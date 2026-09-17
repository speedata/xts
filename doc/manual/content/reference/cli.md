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
| `new [dir]` | Create a starter project with `data.xml` and `layout.xml` |
| `help` | Show help and exit |
| `version` | Print the version number and exit |
| `watch` | Run once, then re-run on every input file change |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-c`, `--config` | string | `xts.cfg` | Configuration file to read |
| `--data` | string | `data.xml` | Data file name |
| `--dummy` | boolean | `false` | Ignore data file, use `<data/>` |
| `--dumpoutput` | string | | Write XML dump of PDF structure to this file |
| `--extradir` | string | | Additional directory for file lookups (recursive); can be given multiple times |
| `--filter` | string | | Lua script to run before publishing |
| `--jobname` | string | `xts` | Output file name (without `.pdf`) |
| `--layout` | string | `layout.xml` | Layout file name |
| `--loglevel` | string | `info` | Log level: `trace`, `debug`, `info`, `notice`, `warn`, `error` |
| `--mode` | string | | Set one or more modes (comma separated list) |
| `--pdfa` | string | | Claim PDF/A archival conformance: `none` or `3b` |
| `--pdfua` | string | | Claim PDF/UA accessibility conformance: `none`, `1`, or `2` |
| `--pdfx` | string | | Claim PDF/X print conformance: `none`, `X-3`, or `X-4` |
| `--runs` | integer | `1` | Number of publishing runs |
| `--quiet` | boolean | `false` | Suppress output on STDOUT |
| `--suppressinfo` | boolean | `false` | Produce reproducible PDF (no timestamps) |
| `--systemfonts` | boolean | `false` | Include system-installed fonts in search |
| `--trace` | string | | Comma-separated traces: `grid`, `gridallocation` |
| `-v`, `--var=VALUE` | string | | Set a variable for the publishing run (`name=value`); can be given multiple times |
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

# Scaffold a new project
xts new myproject

# Re-run automatically on every change to the input files
xts watch
```

## Watch mode

`xts watch` runs the publishing process once and then watches the current
directory, the directories of the layout and data files and all extra
directories (`--extradir`) for changes, including their subdirectories
(hidden directories such as `.git` are skipped). Every change to an input
file starts a new publishing run, saving a file with unchanged content
included, so saving the layout again is a simple way to force a run. Errors
in a run do not stop the watcher. Press `ctrl-c` to quit.

Files written by xts itself (the PDF, protocol and auxiliary files), hidden
files and editor backup files do not trigger a run. A file that is written
while a run is in progress (for example a data file rewritten by a Lua
filter) starts another run only if its content differs from the previous
run.

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
