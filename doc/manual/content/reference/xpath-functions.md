---
type: docs
linktitle: XPath Functions
---

# XPath Functions Reference

XTS extends XPath with functions in the `urn:speedata.de/2021/xtsfunctions/en` namespace, typically bound to the `sd:` prefix. All functions are called as `sd:functionname(...)`.

In addition to the `sd:` functions listed below, all standard XPath functions (such as `string()`, `number()`, `concat()`, `contains()`, `count()`, etc.) are available. XTS uses the [goxml XPath package](https://doc.speedata.de/goxml/xpath/) — see its documentation for the full list of supported standard XPath functions and operators.

## Page and position

`sd:current-page()`
:   Returns the current page number.

`sd:current-row(areaname?)`
:   Returns the current cursor row. Optional area name.

`sd:page-number(markname)`
:   Returns the page number where the given mark was placed.

`sd:last-page-number()`
:   Returns the last page number of the document.

`sd:total-pages(filename)`
:   Returns the number of pages of the referenced PDF (or image) file. For the page count of the document being created, use `sd:last-page-number()`.

## Grid dimensions

`sd:number-of-columns(areaname?)`
:   Returns the number of grid columns. Optional area name.

`sd:number-of-rows(areaname?)`
:   Returns the number of grid rows. Optional area name.

`sd:grid-width(columns, unit?)`
:   Returns the width of the given number of columns. Optional unit (e.g. `'cm'`, `'mm'`).

`sd:grid-height(rows, unit?)`
:   Returns the height of the given number of rows. Optional unit.

## Slates

`sd:slate-width(slatename, unit?)`
:   Returns the width of a named slate: without `unit` as the number of grid columns, with `unit` (e.g. `'cm'`) as a number in that unit.

`sd:slate-height(slatename, unit?)`
:   Returns the height of a named slate: without `unit` as the number of grid rows, with `unit` as a number in that unit.

## Images

`sd:image-width(filename, page?, box?, unit?)`
:   Returns the width of an image: without `unit` as the number of grid columns, with `unit` (one of `'sp'`, `'mm'`, `'cm'`, `'in'`, `'pt'`, `'px'`, `'pc'`, `'m'`) as a number in that unit. Optional page number (for PDF) and box type (`'cropbox'`, `'mediabox'`, `'bleedbox'`, `'trimbox'`, `'artbox'`).

`sd:image-height(filename, page?, box?, unit?)`
:   Returns the height of an image (grid rows without `unit`). Same optional parameters as `image-width`.

`sd:aspect-ratio(filename, page?, box?)`
:   Returns the aspect ratio (width / height) of an image.

## Variables and attributes

`sd:variable(name)`
:   Returns the value of the variable with the given name. Useful for dynamic variable names: `sd:variable(('prefix', $i))` concatenates the arguments into a variable name.

`sd:attribute(name)`
:   Returns the value of the named attribute.

## Text and formatting

`sd:dummy-text()`
:   Returns one paragraph of Lorem Ipsum text.

`sd:markdown(text)`
:   Converts Markdown text to HTML.

`sd:roman-numeral(number)`
:   Converts a number to a Roman numeral string (e.g. `4` → `"IV"`).

`sd:format-number(number, thousands-separator, decimal-separator)`
:   Formats a number with the given thousands separator and decimal separator. Example: `sd:format-number(12345.6, '.', ',')` returns `12.345,6`.

## Math and logic

`sd:even(number)`
:   Returns `true()` if the number is even.

`sd:odd(number)`
:   Returns `true()` if the number is odd.

`sd:mode(name)`
:   Returns `true()` if the given mode is set (via `--mode` on the command line or `mode` in the configuration file).

## Unit conversion

`sd:to-unit(value, unit, precision?)`
:   Converts a value to the given unit, optionally rounded to `precision` decimal places. Example: `sd:to-unit('12pt', 'mm', 2)` returns 4.23.

## File operations

`sd:file-exists(filename)`
:   Returns `true()` if the file exists.

`sd:file-contents(bytes)`
:   Writes the given byte sequence (for example the result of `sd:decode-base64()`) to a temporary file and returns its file name, so the data can be used where a file name is expected (such as `<Image href>`).

## String processing

`sd:decode-html(string)`
:   Decodes HTML entities in a string.

`sd:decode-base64(string)`
:   Decodes a Base64-encoded string.

## Cryptographic

`sd:md5(string)`
:   Returns the MD5 hash of the string.

`sd:sha1(string)`
:   Returns the SHA-1 hash.

`sd:sha256(string)`
:   Returns the SHA-256 hash.

`sd:sha512(string)`
:   Returns the SHA-512 hash.
