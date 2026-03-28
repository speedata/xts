---
weight: 10
type: docs
linktitle: Fonts
---

# Using Fonts

XTS supports TrueType (`.ttf`), OpenType (`.otf`), and Type 1 (`.pfb`/`.afm`) fonts. Loading them is done through CSS `@font-face` rules -- the same syntax you'd use on the web.

## Loading font families

A font family groups the regular, bold, italic, and bold-italic variants under one name. Define them in a `<Stylesheet>` block:

```xml
<Stylesheet>
  @font-face {
      font-family: "Minion Pro";
      src: url("MinionPro-Regular.otf");
  }
  @font-face {
      font-family: "Minion Pro";
      src: url("MinionPro-Bold.otf");
      font-weight: bold;
  }
  @font-face {
      font-family: "Minion Pro";
      src: url("MinionPro-It.otf");
      font-style: italic;
  }
  @font-face {
      font-family: "Minion Pro";
      src: url("MinionPro-BoldIt.otf");
      font-weight: bold;
      font-style: italic;
  }
</Stylesheet>
```

You only need to define the variants you actually use. If you never use bold italic, skip it.

![font size and leading](/manual/img/14-fontsize-leading.png)
<figcaption>Font size and line height (leading).</figcaption>

## Selecting fonts

Once loaded, use the font family name in your CSS rules:

```xml
<Stylesheet>
  body {
    font-family: serif;
  }
  .preface {
    font-family: sans;
  }
</Stylesheet>

<Paragraph>
    <Span class="preface"><Value>Preface</Value></Span>
    <Value> more text</Value>
</Paragraph>
```

## Default fonts

XTS ships with three built-in font families that you can use without loading anything:

| Name | Font | Style |
|------|------|-------|
| `sans` | TeXGyreHeros | Helvetica clone |
| `serif` | CrimsonPro | Elegant book font |
| `monospace` | CamingoCode | Coding font |

If you don't specify a font, XTS uses `sans` (set in the CSS defaults as the `html` element's `font-family`).

## Where to put font files

Font files follow the same rules as all other resources in XTS -- see [File Organization](../../running-xts/file-organization). In short:

- Put them in the project directory (or a subdirectory)
- Or use `--extradir` to add a shared font directory
- Or use `--systemfonts` to access system-installed fonts

## The list-fonts shortcut

To save typing, let XTS scan your fonts and generate the `@font-face` rules for you:

```
$ xts list-fonts
@font-face { font-family: "CamingoCode"; src: url("CamingoCode-Bold.ttf"); font-weight: bold; }
@font-face { font-family: "CamingoCode"; src: url("CamingoCode-Regular.ttf"); }
...
```

Copy the output into your `<Stylesheet>` block.

## OpenType features

OpenType fonts often include optional features like old-style figures, small caps, or fraction rendering. Control them with `font-feature-settings`:

```xml
<Stylesheet>
    .regular {
        font-feature-settings: "lnum", "tnum";
    }
    .smcp {
        font-feature-settings: "smcp";
    }
</Stylesheet>

<Record element="data">
    <PlaceObject>
        <Textblock>
            <Paragraph class="regular">
                <Value>Tabular figures: 1234567890</Value>
            </Paragraph>
            <Paragraph class="smcp">
                <Value>Small caps: 1234567890</Value>
            </Paragraph>
        </Textblock>
    </PlaceObject>
</Record>
```

![opentype features](/manual/img/osfsmcp.png)
<figcaption>Tabular figures (top) have equal width for column alignment. Small caps (bottom) are properly designed, not just shrunken capitals.</figcaption>

You can also set features inline:

```xml
<Paragraph style="font-feature-settings: 'frac';">
    <Value>Use 1/4 cup of milk.</Value>
</Paragraph>
```

![fraction feature](/manual/img/frac-feature-hb.png)
<figcaption>The <code>frac</code> feature turns "1/4" into a proper fraction glyph.</figcaption>

The full list of OpenType feature tags is in the [OpenType spec](https://docs.microsoft.com/en-us/typography/opentype/spec/featurelist). XTS enables the same default features as [HarfBuzz](https://harfbuzz.github.io/shaping-opentype-features.html), minus `liga`.
