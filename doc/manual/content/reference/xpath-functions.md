---
type: docs
linktitle: XPath Functions
---

# XPath Functions Reference

XTS extends XPath with functions in the `urn:speedata.de/2021/xtsfunctions/en` namespace, typically bound to the `sd:` prefix. All functions are called as `sd:functionname(...)`.

## Page and position

`sd:current-page()`
:   Returns the current page number.

`sd:current-row(areaname?)`
:   Returns the current cursor row. Optional area name.

`sd:page-number(markname)`
:   Returns the page number where the given mark was placed.

`sd:last-page-number()`
:   Returns the last page number of the document.

`sd:total-pages(selector)`
:   Returns the total number of pages.

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
:   Returns the width of a named slate. Optional unit.

`sd:slate-height(slatename, unit?)`
:   Returns the height of a named slate. Optional unit.

## Images

`sd:image-width(filename, page?, box?, unit?)`
:   Returns the width of an image. Optional page number (for PDF), box type (`'cropbox'`, `'mediabox'`, `'bleedbox'`, `'trimbox'`, `'artbox'`), and unit.

`sd:image-height(filename, page?, box?, unit?)`
:   Returns the height of an image. Same optional parameters as `image-width`.

`sd:aspect-ratio(filename, page?, box?)`
:   Returns the aspect ratio (width / height) of an image.

## Variables and attributes

`sd:variable(name)`
:   Returns the value of the variable with the given name. Useful for dynamic variable names: `sd:variable(('prefix', $i))` concatenates the arguments into a variable name.

`sd:attribute(name)`
:   Returns the value of the named attribute.

## Text and formatting

`sd:dummy-text(count?)`
:   Returns Lorem Ipsum text. Optional paragraph count.

`sd:markdown(text)`
:   Converts Markdown text to HTML.

`sd:roman-numeral(number)`
:   Converts a number to a Roman numeral string (e.g. `4` → `"IV"`).

`sd:format-number(number, format, locale)`
:   Formats a number according to the given format string and locale.

## Math and logic

`sd:even(number)`
:   Returns `true()` if the number is even.

`sd:odd(number)`
:   Returns `true()` if the number is odd.

`sd:mode(number)`
:   Returns the most frequent value.

## Unit conversion

`sd:to-unit(value, fromunit?, tounit?)`
:   Converts a value between units. Example: `sd:to-unit('12pt', 'pt', 'mm')`.

## File operations

`sd:file-exists(filename)`
:   Returns `true()` if the file exists.

`sd:file-contents(filename)`
:   Returns the file contents as a string.

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
