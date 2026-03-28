---
type: docs
linktitle: Defaults
---

# XTS defaults

XTS defines some default settings that can be changed in the layout file.
These defaults concern the colors, fonts and page size and margins.

## Fonts

These font faces are predefined:

~~~css
@font-face {
    font-family: "monospace";
    src: url("CamingoCode Regular.ttf");
}
@font-face {
    font-family: "monospace";
    src: url("CamingoCode Bold.ttf");
    font-weight: bold;
}
@font-face {
    font-family: "monospace";
    src: url("CamingoCode BoldItalic.ttf");
    font-weight: bold;
    font-style: italic;
}
@font-face {
    font-family: "monospace";
    src: url("CamingoCode Italic.ttf");
    font-style: italic;
}

@font-face {
    font-family: "serif";
    src: url("CrimsonPro-Regular.ttf");
}
@font-face {
    font-family: "serif";
    src: url("CrimsonPro-Bold.ttf");
    font-weight: bold;
}
@font-face {
    font-family: "serif";
    src: url("CrimsonPro-BoldItalic.ttf");
    font-weight: bold;
    font-style: italic;
}
@font-face {
    font-family: "serif";
    src: url("CrimsonPro-Italic.ttf");
    font-style: italic;
}

@font-face {
    font-family: "sans";
    src: url("texgyreheros-regular.otf");
}
@font-face {
    font-family: "sans";
    src: url("texgyreheros-bold.otf");
    font-weight: bold;
}
@font-face {
    font-family: "sans";
    src: url("texgyreheros-bolditalic.otf");
    font-weight: bold;
    font-style: italic;
}
@font-face {
    font-family: "sans";
    src: url("texgyreheros-italic.otf");
    font-style: italic;
}
~~~

The pre-installed fonts can be accessed as `local()` fonts:

~~~
 src: local("CamingoCode Regular")
 src: local("CamingoCode Bold")
 src: local("CamingoCode BoldItalic")
 src: local("CamingoCode Italic")
 src: local("CrimsonPro Regular")
 src: local("CrimsonPro Bold")
 src: local("CrimsonPro BoldItalic")
 src: local("CrimsonPro Italic")
 src: local("TeXGyreHeros Regular")
 src: local("TeXGyreHeros Bold")
 src: local("TeXGyreHeros BoldItalic")
 src: local("TeXGyreHeros Italic")
~~~


## CSS defaults

The CSS defaults are:

~~~css
html            { font-size: 10pt; tab-size: 4; font-family: sans; }
li              { display: list-item; padding-left: 0; }
head            { display: none }
table           { display: table }
tr              { display: table-row }
thead           { display: table-header-group }
tbody           { display: table-row-group }
tfoot           { display: table-footer-group }
td, th          { display: table-cell }
caption         { display: table-caption }
th              { font-weight: bold; text-align: center }
caption         { text-align: center }
body            { margin: 0pt; line-height: 1.2; hyphens: auto; font-weight: normal; }
h1              { font-size: 2em; margin:  .67em 0 }
h2              { font-size: 1.5em; margin: .75em 0 }
h3              { font-size: 1.17em; margin: .83em 0 }
h4, p,
blockquote, ul,
fieldset, form,
ol, dl, dir,
h5              { font-size: 1em; margin: 1.5em 0; text-align: left; }
h6              { font-size: .75em; margin: 1.67em 0 }
h1, h2, h3, h4,
h5, h6, b,
strong          { font-weight: bold }
blockquote      { margin-left: 40px; margin-right: 40px }
i, cite, em,
var, address    { font-style: italic }
pre, tt, code,
kbd, samp       { font-family: monospace; -bag-font-expansion: 0%;}
pre             { white-space: pre; margin: 1em 0px; }
button, textarea,
input, select   { display: inline-block }
big             { font-size: 1.17em }
small, sub, sup { font-size: .83em }
sub             { vertical-align: sub }
sup             { vertical-align: super }
table           { border-spacing: 2pt; }
thead, tbody,
tfoot           { vertical-align: middle }
td, th, tr      { vertical-align: inherit }
s, strike, del  { text-decoration: line-through }
hr              { border: 1px inset }
ol, ul, dir, dd { padding-left: 20pt }
ol              { list-style-type: decimal }
ul              { list-style-type: disc }
ol ul, ul ol,
ul ul, ol ol    { margin-top: 0; margin-bottom: 0 }
u, ins          { text-decoration: underline }
center          { text-align: center }
~~~



## Page size and margin

The page size defaults to A4 (210mm × 297mm).

The master page for all pages is defined as follows:

~~~xml
<DefineMasterpage name="default page" test="true()" margin="1cm" />
~~~

The page grid is set to 10mm × 10mm.

## Colors

The known CSS colors are defined in the RGB color space. The colors 'black' and 'white' are defined in the grayscale color space. See also the command [`<DefineColor>`](/reference/commands/definecolor), there the predefined colors are listed.

The special colors HKS 1-97 and many Pantone colors are already defined with their CMYK values.

