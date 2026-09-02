---
weight: 50
type: docs
linktitle: Quality Assurance
---

# Quality Assurance and PDF Comparison

XTS has built-in support for regression testing: compare a current PDF against a known-good reference to catch unwanted visual changes.

## How it works

1. Create your layout and data files
2. Generate a reference PDF: `xts --jobname reference`
3. Clean up temp files: `xts --jobname reference clean`
4. Later, run `xts compare <directory>` to re-generate and compare

## Setting up a test case

```
example/
├── data.xml
├── layout.xml
└── reference.pdf
```

Run the comparison:

```
$ xts compare example/
Finished in 1s
```

If only the finish line appears, everything matches. If there's a difference:

```
$ xts compare example/
---------------------------
Finished with comparison in
/path/to/example
Comparison failed. Bad pages are: [0]
Max delta is 2162.76
```

XTS generates difference images, one set per page (`-00`, `-01`, ...):

```
example/
├── data.xml
├── layout.xml
├── pagediff-00.png   ← highlighted differences
├── xts.pdf           ← current output
├── reference.pdf     ← known-good reference
├── reference-00.png  ← reference as bitmap
└── source-00.png     ← current output as bitmap
```

## Running a test suite

Organize test cases in a directory tree:

```
qa/
├── test-tables/
│   ├── data.xml
│   ├── layout.xml
│   └── reference.pdf
├── test-images/
│   ├── data.xml
│   ├── layout.xml
│   └── reference.pdf
└── test-fonts/
    ├── data.xml
    ├── layout.xml
    └── reference.pdf
```

```
$ xts compare qa/
Finished in 4s
```

XTS recursively finds all directories containing a `reference.pdf` and runs each test. A directory with a `layout.xml` but no `reference.pdf` is reported with a warning and skipped.

## Faster comparisons

Use `--suppressinfo` to create reproducible PDFs without timestamps:

```
xts --suppressinfo --jobname reference
xts --jobname reference clean
```

If the checksum matches, the visual comparison is skipped entirely.

## HTML report

When `xts compare` finds differences, an HTML report `compare-report.html` is created in the current directory. Open it in a browser for a visual overview of the failed tests. With `--verbose`, the report is always written and shows all pages of all tests.

## Prerequisites

The visual comparison requires [ImageMagick](https://imagemagick.org/) and [Ghostscript](https://www.ghostscript.com/) (for converting the reference PDF to bitmaps) to be installed.
