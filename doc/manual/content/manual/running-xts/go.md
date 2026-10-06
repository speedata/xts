---
weight: 70
type: docs
linktitle: Running XTS from Go
---

# Running XTS from Go

The `xts` command is a thin layer around the Go package `github.com/speedata/xts/core`. A Go program can call it directly, for example to make PDFs in a web service or to register its own line models.

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/speedata/xts/core"
)

func main() {
	layout, err := os.Open("layout.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer layout.Close()
	data, err := os.Open("data.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	cfg := &core.XTSConfig{
		Layoutfile:  layout,
		Datafile:    data,
		FindFile:    core.FindFile,
		Jobname:     "report",
		OutFilename: "report.pdf",
		Variables:   map[string]any{"title": "Monthly report"},
	}
	if err := core.RunXTS(cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.Info.Pages, "pages,", cfg.Info.FileSize, "bytes")
}
```

`RunXTS` typesets one run and writes the PDF to `OutFilename`. Afterwards `cfg.Info` holds the number of pages and the size of the file.

## The configuration

| Field | Meaning | Command line |
|---|---|---|
| `Layoutfile`, `Datafile` | The layout and the data as readers | `--layout`, `--data` |
| `FindFile` | Turns a file name from the layout (images, fonts, XML) into a path. `core.FindFile` searches the directories added with `core.AddDir`, then the working directory, and downloads `http` and `https` URLs | `--extradir` |
| `Jobname` | Base name of the files XTS writes besides the PDF, such as `report-aux.xml` | `--jobname` |
| `OutFilename` | The PDF file | |
| `Variables` | Variables the layout can read, here `$title` | `--var` |
| `Mode` | The modes `sd:mode()` tests for | `--mode` |
| `SuppressInfo` | A reproducible PDF, with a fixed date and no random IDs | `--suppressinfo` |
| `Tracing` | `grid` and `gridallocation` draw the grid and its allocation | `--trace` |
| `DumpFile` | Receives the XML dump of the output | `--dumpoutput` |
| `Pdfua`, `Pdfa`, `Pdfx` | The PDF standards to claim, empty for none | `--pdfua`, `--pdfa`, `--pdfx` |
| `LineModels` | Line models by name, see below | |
| `Breakers` | Line breakers by name, see below | |

The command line does a few things around `RunXTS` that a Go program does itself when it needs them: it reads `xts.cfg`, adds the directories, runs the Lua filter, and calls `RunXTS` again for `--runs`, so that cross-references written to the aux file in one run are resolved in the next. Messages go to the default logger of `log/slog`; the protocol file `xts-protocol.xml` is written by the command line.

## Line models

`LineModels` maps a name to an `htmlbag.LineModelFunc`, a function that gets the font size, line height and language of a paragraph and returns a `node.LineModel` from boxes and glue. A paragraph selects it with `-bag-leading-model` in CSS:

```go
cfg.LineModels = map[string]htmlbag.LineModelFunc{
	"word": wordLineModel,
}
```

```css
p { -bag-leading-model: word; }
```

The `ascent-override`, `descent-override` and `line-gap-override` descriptors of `@font-face` are only read by such a line model. The built-in models `half` (CSS line boxes, the default) and `trailing` (TeX style) do not use them.

## Breakers

`Breakers` maps a name to an `htmlbag.BreakerFunc`, a function that gets the name, font size and language of a paragraph (`htmlbag.BreakerStyles`) and returns a `node.Breaker` from boxes and glue. A paragraph selects it with `-bag-line-breaker` in CSS, and the breaker then chooses where the paragraph breaks among its legal breakpoints in place of Knuth-Plass; the lines are measured and set as before:

```go
cfg.Breakers = map[string]htmlbag.BreakerFunc{
	"greedy": func(htmlbag.BreakerStyles) node.Breaker { return firstFit{} },
}
```

```css
p { -bag-line-breaker: greedy; }
```

Here `firstFit` is a type whose `Breaks` method fills each line with as many words as fit. A function that returns `nil` keeps Knuth-Plass for that paragraph. `auto`, the default, is Knuth-Plass; the name is reserved and cannot be registered.
