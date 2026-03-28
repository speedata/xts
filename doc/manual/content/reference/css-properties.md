---
type: docs
linktitle: CSS Properties
---

# CSS Properties Reference

CSS properties supported by XTS for styling layout elements and HTML content.

## Text

| Property | Values | Example |
|----------|--------|---------|
| `font-family` | Font family name | `font-family: serif;` |
| `font-size` | Length or em | `font-size: 12pt;` |
| `font-weight` | `normal`, `bold` | `font-weight: bold;` |
| `font-style` | `normal`, `italic` | `font-style: italic;` |
| `font-feature-settings` | OpenType feature tags | `font-feature-settings: "smcp";` |
| `color` | Color value | `color: #333;` |
| `text-align` | `left`, `right`, `center`, `justify` | `text-align: justify;` |
| `text-indent` | Length | `text-indent: 1em;` |
| `text-decoration` | `none`, `underline`, `line-through` | `text-decoration: underline;` |
| `line-height` | Number or length | `line-height: 1.4;` |
| `white-space` | `normal`, `pre` | `white-space: pre;` |
| `hyphens` | `auto`, `none` | `hyphens: auto;` |
| `vertical-align` | `top`, `middle`, `bottom`, `sub`, `super` | `vertical-align: top;` |

## Box model

| Property | Values | Example |
|----------|--------|---------|
| `margin` | Length (shorthand or per side) | `margin: 10pt;` |
| `margin-top` | Length | `margin-top: 12pt;` |
| `margin-right` | Length | `margin-right: 12pt;` |
| `margin-bottom` | Length | `margin-bottom: 12pt;` |
| `margin-left` | Length | `margin-left: 12pt;` |
| `padding` | Length (shorthand or per side) | `padding: 5pt 10pt;` |
| `padding-top` | Length | `padding-top: 5pt;` |
| `padding-right` | Length | `padding-right: 10pt;` |
| `padding-bottom` | Length | `padding-bottom: 5pt;` |
| `padding-left` | Length | `padding-left: 10pt;` |

## Borders

| Property | Values | Example |
|----------|--------|---------|
| `border` | Width style color | `border: 1pt solid black;` |
| `border-top` | Width style color | `border-top: 2pt solid red;` |
| `border-right` | Width style color | `border-right: 1pt dashed gray;` |
| `border-bottom` | Width style color | `border-bottom: 0.5pt solid #ccc;` |
| `border-left` | Width style color | `border-left: 1pt solid blue;` |
| `border-radius` | Length | `border-radius: 3pt;` |
| `border-spacing` | Length | `border-spacing: 2pt;` |

Border styles: `solid`, `dashed`.

## Background

| Property | Values | Example |
|----------|--------|---------|
| `background-color` | Color value | `background-color: #ffffcc;` |

## Display and list

| Property | Values | Example |
|----------|--------|---------|
| `display` | `block`, `inline`, `table`, `list-item`, `none`, etc. | `display: none;` |
| `list-style-type` | `disc`, `decimal`, `none` | `list-style-type: decimal;` |

## Table

| Property | Values | Example |
|----------|--------|---------|
| `width` | Length or percentage | `width: 100%;` |

## Selectors

XTS supports standard CSS selectors:

| Selector | Example |
|----------|---------|
| Element | `p { ... }` |
| Class | `.highlight { ... }` |
| ID | `#main { ... }` |
| Descendant | `table td { ... }` |
| Child | `table > tr { ... }` |
| Pseudo-class | `tr:nth-child(even) { ... }` |
| Multiple | `h1, h2, h3 { ... }` |

## @-rules

| Rule | Description |
|------|-------------|
| `@font-face` | Define a font family with `font-family`, `src`, `font-weight`, `font-style` |
