---
weight: 85
type: docs
linktitle: Accessible PDF
---

# Accessible PDF documents (PDF/UA)

PDF/UA (ISO 14289, "Universal Accessibility") is the standard for accessible,
*tagged* PDF. A PDF/UA document carries a logical structure tree (headings,
paragraphs, lists, figures, formulas, …) that assistive technology such as
screen readers uses to convey the document in a meaningful reading order.

XTS can produce PDF/UA-1 and PDF/UA-2 output. This chapter explains how to turn
it on, what your layout has to provide, and how to verify the result.

## Turning on conformance

The output conformance is configured on the [command line](../running-xts/command-line)
or in the [configuration file](../running-xts/configuration).

```bash
# PDF/UA-2 (based on PDF 2.0)
xts --pdfua 2
```

```toml title="xts.cfg"
pdfua = "2"
```

There are three orthogonal conformance axes; each is set independently:

| Axis    | Values               | Standard               |
| ------- | -------------------- | ---------------------- |
| `pdfua` | `none`, `1`, `2`     | PDF/UA (accessibility) |
| `pdfa`  | `none`, `3b`         | PDF/A (archiving)      |
| `pdfx`  | `none`, `X-3`, `X-4` | PDF/X (printing)       |

Axes from different families can be combined as long as they agree on the base
PDF version, for example an accessible and archivable PDF:

```bash
xts --pdfa 3b --pdfua 1
```

{{< callout type="info" >}}
**Version coupling.** PDF/UA-2 requires PDF 2.0, while PDF/A-3 and PDF/X-3/4 are
PDF 1.x. XTS rejects contradictory combinations such as `--pdfua 2 --pdfa 3b`
with a clear error. Pair PDF/UA-1 with PDF/A-3 (both PDF 1.x), or use PDF/UA-2 on
its own.
{{< /callout >}}

## What an accessible document needs

Turning on `pdfua` makes XTS emit the structure tree, marked content, language
and metadata automatically. A few things, however, have to come from your
layout and data.

### A document title

PDF/UA requires a title that is shown instead of the file name. Set it with
[`<PDFOptions>`](../advanced/pdf-options):

```xml
<PDFOptions title="Annual Report 2026"/>
```

{{< callout type="info" >}}
If no title is set by the time the metadata is written, XTS falls back to a
placeholder title and prints a warning. The document stays conformant, but you
should always set a real, meaningful title.
{{< /callout >}}

### A document language

Declare the natural language of the document with [`<Options>`](/reference/commands/options):

```xml
<Options mainlanguage="en"/>
```

### Logical structure: headings, lists, figures, formulas

The structure tree is built from the HTML-like content you place with
[`<HTML>`](/reference/commands/html). Use real semantic elements — headings
`<h1>`–`<h6>`, paragraphs `<p>`, lists `<ul>`/`<ol>`, tables — and they become
the corresponding structure elements.

```xml
<HTML>
  <h1>Quarterly figures</h1>
  <p>Revenue rose in all regions.</p>
  <ul>
    <li>Europe: +4%</li>
    <li>Asia: +7%</li>
  </ul>
</HTML>
```

**Images** must carry alternative text so they are tagged as a `Figure` with
`/Alt`:

```xml
<HTML>
  <img src="chart.pdf" alt="Bar chart of revenue per region"/>
</HTML>
```

**Formulas** are written as [MathML](https://developer.mozilla.org/en-US/docs/Web/MathML)
inside `<HTML>`. They are tagged as a `Formula` structure element with an
alternative text and the original MathML attached as an associated file. Inline
and block (`display="block"`) formulas are both supported. A MathML formula
needs a font with an OpenType `MATH` table, declared in CSS:

```xml
<StyleSheet>
  @font-face { font-family: "Latin Modern Math"; src: url("latinmodern-math.otf"); }
  math { font-family: "Latin Modern Math"; }
</StyleSheet>
```

```xml
<HTML>
  <p>The Pythagorean theorem is
    <math><msup><mi>a</mi><mn>2</mn></msup><mo>+</mo><msup><mi>b</mi><mn>2</mn></msup><mo>=</mo><msup><mi>c</mi><mn>2</mn></msup></math>.</p>
</HTML>
```

## A complete minimal example

```toml title="xts.cfg"
pdfua = "2"
```

```xml title="layout.xml"
<Layout xmlns="urn:speedata.de/2021/xts/en">
  <Options mainlanguage="en"/>
  <PDFOptions title="Accessible report"/>

  <Record match="data">
    <PlaceObject>
      <TextBlock>
        <HTML>
          <h1>Accessible report</h1>
          <p>This document is tagged for accessibility.</p>
        </HTML>
      </TextBlock>
    </PlaceObject>
  </Record>
</Layout>
```

```bash
xts
```

## Verifying accessibility

Generating a PDF/UA document does not guarantee it is *meaningfully* accessible —
alt texts have to be sensible, the reading order has to make sense, and so on.
Always verify the result.

The easiest way is the online checker at
[pdfuacheck.speedata.de](https://pdfuacheck.speedata.de): upload your PDF and it
reports, per checkpoint, what passes and what fails, including the document's
structure tree.

Other established tools:

- **veraPDF** — the reference open-source PDF/UA and PDF/A validator.
- **PAC (PDF Accessibility Checker)** — a free Windows tool from the PDF/UA
  Foundation with a visual structure inspector.

A typical failure and its cause:

| Report                          | Cause                                           |
| ------------------------------- | ----------------------------------------------- |
| No `dc:title` / title missing   | No `title` on `<PDFOptions>`                    |
| Document language not set       | No `mainlanguage` on `<Options>`                |
| Figure without alternative text | `<img>` without an `alt` attribute              |
| Untagged content                | Content placed outside the tagged `<HTML>` flow |

## See also

- [PDF Options](../advanced/pdf-options)
- [Command line](../running-xts/command-line)
- [Configuration](../running-xts/configuration)
- [PDFOptions reference](/reference/commands/pdfoptions)
- [Options reference](/reference/commands/options)
