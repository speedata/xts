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
Total run time: 1.62956s
```

No output means everything matches. If there's a difference:

```
$ xts compare example/
/path/to/example
Comparison failed. Bad pages are: [0]
Max delta is 2162.760009765625
```

XTS generates difference images:

```
example/
├── data.xml
├── layout.xml
├── pagediff.png      ← highlighted differences
├── publisher.pdf     ← current output
├── reference.pdf     ← known-good reference
├── reference.png     ← reference as bitmap
└── source.png        ← current output as bitmap
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
Total run time: 4.54s
```

XTS recursively finds all directories containing `layout.xml` and runs each test.

## Faster comparisons

Use `--suppressinfo` to create reproducible PDFs without timestamps:

```
xts --suppressinfo --jobname reference
xts --jobname reference clean
```

If the checksum matches, the visual comparison is skipped entirely.

## HTML report

After running `xts compare`, an HTML report `compare-report.html` is created in the current directory. Open it in a browser for a visual overview of all test results.

## Prerequisites

The visual comparison requires [ImageMagick](https://imagemagick.org/) to be installed.
