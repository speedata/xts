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

<!-- Internal link (to a mark) -->
<Mark select="'chapter1'"/>
<!-- ... later ... -->
<A link="chapter1"><Value>See Chapter 1</Value></A>
```

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
    duplex="DuplexFlipLongEdge"
    printscaling="None"
    showhyperlinks="false"/>
```

## See also

- [Bookmark reference](/reference/commands/bookmark)
- [A reference](/reference/commands/a)
- [PDFOptions reference](/reference/commands/pdfoptions)
- [Mark reference](/reference/commands/mark)
