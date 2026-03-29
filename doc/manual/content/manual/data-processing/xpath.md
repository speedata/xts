---
weight: 40
type: docs
linktitle: XPath
---

# XPath and Layout Functions

XTS uses XPath as its expression language. XPath appears in `select` attributes, `test` conditions, and inside curly braces `{...}` in some contexts.

If you're new to XPath, the [W3Schools XPath tutorial](https://www.w3schools.com/xml/xpath_intro.asp) is a good starting point. This page focuses on the XTS-specific extensions.

## Where XPath is used

```xml
<!-- select: evaluate and use the result -->
<Value select="@price"/>
<Value select="concat(@name, ' (', @sku, ')')"/>
<ForAll select="article[@active = 'yes']"/>

<!-- test: evaluate to true/false -->
<Case test="$count > 10"/>
<DefineMasterPage test="sd:current-page() > 1"/>

<!-- Curly braces in attributes -->
<Column width="{sd:grid-width(3)}"/>
<SetVariable variable="{ concat('item', $i) }"/>
```

## XTS layout functions

XTS extends XPath with functions in the `sd:` namespace. These let you query the layout state at runtime.

### Page and position

| Function | Returns |
|----------|---------|
| `sd:current-page()` | Current page number |
| `sd:current-row('area')` | Current cursor row (optional area name) |
| `sd:page-number('markname')` | Page number where a mark was placed |
| `sd:last-page-number()` | Number of the last page |
| `sd:total-pages('selector')` | Total page count |

### Grid dimensions

| Function | Returns |
|----------|---------|
| `sd:number-of-columns('area')` | Number of grid columns (optional area name) |
| `sd:number-of-rows('area')` | Number of grid rows (optional area name) |
| `sd:grid-width(columns, 'unit')` | Width of N columns in the given unit |
| `sd:grid-height(rows, 'unit')` | Height of N rows in the given unit |

### Slates and images

| Function | Returns |
|----------|---------|
| `sd:slate-width('name', 'unit')` | Width of a named slate |
| `sd:slate-height('name', 'unit')` | Height of a named slate |
| `sd:image-width('file', page, 'box', 'unit')` | Image width |
| `sd:image-height('file', page, 'box', 'unit')` | Image height |
| `sd:aspect-ratio('file', page, 'box')` | Image aspect ratio (width/height) |

### Variables

| Function | Returns |
|----------|---------|
| `sd:variable('name')` | Value of a variable (useful for dynamic names) |
| `sd:attribute('name')` | Value of an attribute |

### Text and formatting

| Function | Returns |
|----------|---------|
| `sd:dummy-text(count)` | Lorem ipsum text (optional paragraph count) |
| `sd:markdown('text')` | Convert Markdown to HTML |
| `sd:roman-numeral(number)` | Roman numeral string |
| `sd:format-number(number, 'format', 'locale')` | Formatted number string |

### Utility

| Function | Returns |
|----------|---------|
| `sd:even(number)` | True if number is even |
| `sd:odd(number)` | True if number is odd |
| `sd:file-exists('filename')` | True if file exists |
| `sd:file-contents('filename')` | File contents as string |
| `sd:to-unit(value, 'from', 'to')` | Unit conversion |

### Cryptographic

| Function | Returns |
|----------|---------|
| `sd:md5('string')` | MD5 hash |
| `sd:sha1('string')` | SHA-1 hash |
| `sd:sha256('string')` | SHA-256 hash |
| `sd:sha512('string')` | SHA-512 hash |
| `sd:decode-base64('string')` | Base64 decoded string |
| `sd:decode-html('string')` | HTML entity decoded string |

## XPath in HTML content

Inside `<HTML>` elements, XPath expressions are only evaluated when `expand-text="yes"` is set:

```xml
<HTML expand-text="yes">
    <p>Page {sd:current-page()} of {sd:last-page-number()}</p>
</HTML>
```

## Full reference

For the complete function signatures with all optional parameters, see the [XPath Functions Reference](/reference/xpath-functions).
