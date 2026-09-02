---
weight: 90
type: docs
linktitle: XPath in XTS
aliases:
  - /manual/data-processing/xpath/
---

# XPath and Layout Functions

XPath is a query language for selecting nodes and computing values in XML documents. It lets you navigate the XML tree, filter elements by conditions, and combine values with expressions — similar to how file paths navigate a directory tree, but much more powerful.

XTS uses XPath as its expression language. XPath appears in `select` attributes, `test` conditions, and inside curly braces `{...}` in some contexts. The XPath implementation is provided by the [goxpath package](https://github.com/speedata/goxpath) — see its documentation for the full list of supported XPath functions and syntax details.

If you're new to XPath, start with [XPath basics](/programming/xpath-basics) for a
practical primer (paths, predicates, operators, functions). This page focuses on
*where* XPath is used in XTS and the XTS-specific `sd:` extension functions.

## Where XPath is used

```xml
<!-- select: evaluate and use the result -->
<Value select="@price"/>
<Value select="concat(@name, ' (', @sku, ')')"/>
<ForAll select="article[@active = 'yes']"/>

<!-- test: evaluate to true/false -->
<Case test="$count > 10"/>
<DefineMasterPage name="followingpages" margin="1cm" test="sd:current-page() > 1"/>

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
| `sd:total-pages('filename')` | Page count of the referenced PDF file |

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
| `sd:dummy-text()` | One paragraph of Lorem ipsum text |
| `sd:markdown('text')` | Convert Markdown to HTML |
| `sd:roman-numeral(number)` | Roman numeral string |
| `sd:format-number(number, 'thousands-sep', 'decimal-sep')` | Formatted number string |

### Utility

| Function | Returns |
|----------|---------|
| `sd:even(number)` | True if number is even |
| `sd:odd(number)` | True if number is odd |
| `sd:file-exists('filename')` | True if file exists |
| `sd:file-contents(bytes)` | Write bytes to a temporary file, return its name |
| `sd:to-unit(value, 'unit', precision?)` | Unit conversion |

### Cryptographic

| Function | Returns |
|----------|---------|
| `sd:md5('string')` | MD5 hash |
| `sd:sha1('string')` | SHA-1 hash |
| `sd:sha256('string')` | SHA-256 hash |
| `sd:sha512('string')` | SHA-512 hash |
| `sd:decode-base64('string')` | Base64 decoded string |
| `sd:decode-html('string')` | Parse a string containing HTML markup into nodes |

## XPath in HTML content

Inside `<HTML>` elements, the `{...}` text templates in the body are only evaluated when `expand-text="yes"` is set (the `select` attribute and embedded XTS commands are always evaluated):

```xml
<HTML expand-text="yes">
    <p>Page {sd:current-page()} of {sd:last-page-number()}</p>
</HTML>
```

## Maps and arrays

The XPath engine is version 3.1, so it also supports structured data: arrays
(`[1, 2, 3]`), maps (`map { 'a': 1 }`), and the `?` lookup operator. These have
their own chapter: [Maps and arrays](/programming/maps-and-arrays).

## Full reference

For the complete function signatures with all optional parameters, see the [XPath Functions Reference](/reference/xpath-functions).
