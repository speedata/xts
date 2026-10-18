---
weight: 10
type: docs
linktitle: PDF Options
---

# PDF Options

## Bookmarks

Create a navigable table of contents in the PDF viewer:

```xml
<Bookmark select="'Chapter 1'" level="1"/>
<Bookmark select="'Section 1.1'" level="2"/>
<Bookmark select="'Section 1.2'" level="2"/>
<Bookmark select="'Chapter 2'" level="1"/>
```

Use `open="yes"` to expand a bookmark node by default.

## Links

Create clickable links with `<A>`:

```xml
<!-- External link -->
<Paragraph>
    <A href="https://example.com"><Value>Visit our website</Value></A>
</Paragraph>

<!-- Internal link (to a mark that is a PDF target) -->
<Action><Mark select="'chapter1'" pdftarget="yes"/></Action>
<!-- ... later ... -->
<Paragraph>
    <A link="chapter1"><Value>See Chapter 1</Value></A>
</Paragraph>
```

`link` points to a named destination in the PDF, and only a `<Mark>` with `pdftarget="yes"` makes one. Without it the mark only records its page for `sd:page-number()`, and the link has no target to jump to.

### Page destinations

Every page automatically gets a named destination (`page-1`, `page-2`, ...). You can link to a specific page using the `page` attribute:

```xml
<A page="2"><Value>Go to page 2</Value></A>
```

This also works in HTML mode (htmlbag) using `<a href="#page-2">Go to page 2</a>`.

## PDF metadata

Set document properties with `<PDFOptions>`:

```xml
<PDFOptions
    author="ACME Corp"
    title="Product Catalog 2026"
    subject="Complete product listing"
    creator="XTS"/>
```

## Viewer preferences

Control how the PDF opens:

```xml
<PDFOptions
    displaymode="fullscreen"
    duplex="duplexfliplongedge"
    printscaling="none"
    showhyperlinks="no"/>
```

## Conformance (PDF/UA, PDF/A, PDF/X)

The output conformance is **not** set with `<PDFOptions>`. It has to be known before the document is built, so it is configured on the [command line](../../running-xts/command-line) or in the [configuration file](../../running-xts/configuration) instead, using the orthogonal axes `pdfua`, `pdfa` and `pdfx`:

```bash
xts --pdfua 2
```

```toml title="xts.cfg"
pdfua = "2"
```

For PDF/UA you should also set a meaningful document `title` (see above); without it XTS emits a placeholder title and a warning. See the [Accessible PDF documents](../../accessibility) chapter for the full workflow.

## See also

- [Accessible PDF documents](../../accessibility)
- [Bookmark reference](/reference/commands/bookmark)
- [A reference](/reference/commands/a)
- [PDFOptions reference](/reference/commands/pdfoptions)
- [Mark reference](/reference/commands/mark)
