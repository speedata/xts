---
type: docs
linktitle: CSS Properties
---

# CSS Properties Reference

CSS properties supported by XTS for styling layout elements and HTML content. Lengths take the [units](../units) of XTS; `em` and `%` are relative to the font size or the containing block as in CSS.

<!-- Keep this page in sync with htmlbag: StylesToStyles in
     inheritablestyles.go (longhands), computedstyles.go (shorthands),
     doPage and doFontFace in css.go (at-rules). -->

## Text

| Property | Values | Example |
|----------|--------|---------|
| `font-family` | Font family name, a comma separated list falls back per glyph | `font-family: "Minion Pro", serif;` |
| `font-size` | Length, `em`, `%`, keywords `xx-small` to `xxx-large`, `smaller`, `larger` | `font-size: 12pt;` |
| `font-weight` | `normal`, `bold`, a number from 100 to 900 | `font-weight: 600;` |
| `font-style` | `normal`, `italic`, `oblique` | `font-style: italic;` |
| `font` | Shorthand, must contain size and family | `font: italic 10pt/12pt serif;` |
| `font-feature-settings` | OpenType feature tags | `font-feature-settings: "smcp", "onum";` |
| `font-variation-settings` | Axes of a variable font | `font-variation-settings: "wght" 650;` |
| `color` | Color value or a defined color name | `color: #333;` |
| `text-align` | `left`, `right`, `center`, `justify`, `start`, `end` | `text-align: justify;` |
| `text-indent` | Length | `text-indent: 1em;` |
| `text-decoration` | Shorthand of the three below | `text-decoration: underline dotted red;` |
| `text-decoration-line` | `none`, `underline`, `overline`, `line-through` | `text-decoration-line: underline;` |
| `text-decoration-style` | `solid`, `double`, `dotted`, `dashed`, `wavy` | `text-decoration-style: wavy;` |
| `text-decoration-color` | Color value | `text-decoration-color: red;` |
| `line-height` | Number or length, see [Fonts](/manual/text-and-styling/fonts#font-size-and-line-height) | `line-height: 1.4;` |
| `letter-spacing` | Length | `letter-spacing: 0.05em;` |
| `white-space` | `normal`, `nowrap`, `pre`, `pre-wrap`, `pre-line` | `white-space: pre;` |
| `hyphens` | `auto`, `manual`, `none` | `hyphens: none;` |
| `hanging-punctuation` | `none`, `allow-end` | `hanging-punctuation: allow-end;` |
| `vertical-align` | `baseline`, `sub`, `super`, `top`, `middle`, `bottom`, a length | `vertical-align: super;` |
| `direction` | `ltr`, `rtl` | `direction: rtl;` |
| `unicode-bidi` | `normal`, `embed`, `isolate`, `bidi-override`, `isolate-override`, `plaintext` | `unicode-bidi: isolate;` |
| `tab-size` | Number of spaces | `tab-size: 4;` |
| `initial-letter` | Number of lines for a drop cap | `initial-letter: 3;` |

## Box model

| Property | Values | Example |
|----------|--------|---------|
| `margin` | Length, shorthand for the four sides | `margin: 10pt;` |
| `margin-top`, `margin-right`, `margin-bottom`, `margin-left` | Length | `margin-top: 12pt;` |
| `padding` | Length, shorthand for the four sides | `padding: 5pt 10pt;` |
| `padding-top`, `padding-right`, `padding-bottom`, `padding-left` | Length | `padding-left: 10pt;` |
| `width` | Length or percentage, on blocks, images and table cells | `width: 100%;` |
| `height` | Length, on blocks and images | `height: 4cm;` |
| `background-color` | Color value, painted on block elements and table cells | `background-color: #ffffcc;` |
| `background` | Shorthand, only the color is read | `background: #ffffcc;` |

## Borders

| Property | Values | Example |
|----------|--------|---------|
| `border` | Width style color, all four sides | `border: 1pt solid black;` |
| `border-top`, `border-right`, `border-bottom`, `border-left` | Width style color | `border-top: 2pt solid red;` |
| `border-width`, `border-style`, `border-color` | One value per side | `border-width: 1pt 0;` |
| `border-*-width`, `border-*-style`, `border-*-color` | Longhands per side | `border-left-width: 3pt;` |
| `border-radius` | Length, one to four values | `border-radius: 3pt;` |
| `border-top-left-radius` and the other three corners | Length | `border-top-left-radius: 3pt;` |
| `border-spacing` | Length, on tables | `border-spacing: 2pt;` |

Border styles: `none`, `solid`, `dashed`, `dotted`, `double`.

## Layout

| Property | Values | Example |
|----------|--------|---------|
| `display` | `block`, `inline`, `list-item`, `table` and the table parts, `none` | `display: none;` |
| `float` | `left`, `right`, `none` | `float: right;` |
| `clear` | `left`, `right`, `both`, `none` | `clear: both;` |
| `position` | `static`, `relative`, `absolute` | `position: absolute;` |
| `top`, `right`, `bottom`, `left` | Length, with `position` | `top: 1cm;` |
| `z-index` | Number, with `position` | `z-index: 1;` |
| `page-break-before`, `page-break-after` | `auto`, `always`, `avoid` | `page-break-before: always;` |
| `page-break-inside` | `auto`, `avoid` | `page-break-inside: avoid;` |

## Lists and generated content

| Property | Values | Example |
|----------|--------|---------|
| `list-style-type` | `disc`, `circle`, `square`, `decimal`, `none` | `list-style-type: decimal;` |
| `list-style` | Shorthand, also takes `inside` or `outside` | `list-style: square inside;` |
| `counter-reset`, `counter-increment` | Counter name and optional number | `counter-reset: section;` |
| `content` | On `::before`, `::after` and `::marker`: strings, `attr()`, `counter()`, `counters()`, `target-counter()`, `target-text()`, `element()` | `content: counter(section) ". ";` |

## XTS specific properties

| Property | Values | Example |
|----------|--------|---------|
| `-bag-font-expansion` | Percentage of allowed glyph stretching, 0% turns it off | `-bag-font-expansion: 0%;` |
| `-bag-italic-correction` | `auto`, `none` | `-bag-italic-correction: none;` |
| `-bag-leading-model` | `half` (CSS line boxes, the default), `trailing` (TeX style) | `-bag-leading-model: trailing;` |
| `-bag-linebreak-tolerance` | Number, the TeX tolerance for line breaking | `-bag-linebreak-tolerance: 500;` |
| `-bag-linebreak-hyphen-penalty` | Number, the TeX hyphen penalty | `-bag-linebreak-hyphen-penalty: 200;` |
| `-bag-bookmark` | `none`, or a level number optionally followed by `open` or `closed`, adds the element to the PDF outline | `-bag-bookmark: 2 closed;` |

## Selectors

XTS supports the CSS selectors of levels 1 to 3 and the most common of level 4:

| Selector | Example |
|----------|---------|
| Element | `p { ... }` |
| Class | `.highlight { ... }` |
| ID | `#main { ... }` |
| Attribute | `td[align="right"] { ... }` |
| Descendant | `table td { ... }` |
| Child | `table > tr { ... }` |
| Sibling | `h1 + p { ... }`, `h1 ~ p { ... }` |
| Pseudo-class | `tr:nth-child(even) { ... }`, `li:first-child { ... }`, `p:not(.intro) { ... }` |
| Pseudo-element | `h1::before { ... }`, `li::marker { ... }` |
| Multiple | `h1, h2, h3 { ... }` |

## @-rules

| Rule | Description |
|------|-------------|
| `@font-face` | Define a font face with `font-family`, `src` (`url()` or `local()`), `font-weight`, `font-style` and `font-feature-settings`, see [Fonts](/manual/text-and-styling/fonts) |
| `@-bag-color` | Define a named color with `value` or a `model` and its components, see [Colors](/manual/advanced/colors#defining-colors-in-css) |
| `@page` | Margins of a master page and the margin boxes `@top-left` to `@bottom-right`, see [Master Pages](/manual/page-layout/master-pages#page-setup-from-css) |
