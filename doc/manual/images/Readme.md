# Sources of the manual images

Every subdirectory holds the XTS layout that produces one image of the
manual. `rake docimages` (or `rake docimages[name]` for a single one)
renders each `layout.xml` with the freshly built `bin/xts` and converts
the first page to `content/manual/img/<name>.png` with pdftoppm at
200 dpi. The PNG is committed, the intermediate files are ignored.

Rules for an image directory:

- The name of the directory is the name of the PNG.
- The `<PageFormat>` is the size of the image, so there is no cropping.
  180mm wide gives about 1400 pixels, which fits the width of the manual.
- `data.xml` is optional. Without it the layout runs with `--dummy`.
- Fonts come from the XTS defaults (`serif`, `sans`, `monospace`) or from
  files inside the directory, never from the system.

## Annotations

Guide lines, arrows and labels are not drawn by hand afterwards. They
live in an SVG file next to the layout that covers the whole page. Placed
as the first object it lies under the content (shaded areas), placed
last it lies on top (arrows through the text):

```xml
<PlaceObject column="0mm" row="0mm">
  <Image href="annotations.svg" width="180mm" height="58mm"/>
</PlaceObject>
```

The SVG uses a `viewBox` in PostScript points with the origin at the top
left of the page, so the coordinates in the file are page coordinates.
Positions that depend on the typeset text (baselines, box edges) are
read from the dump: run `xts --dumpoutput dump.xml` in the directory,
look for the `<hlist attr-origin="line">` elements and subtract their
`y` from the page height. Write the measured values and where they come
from into a comment in the SVG, so the next change to the layout can
repeat the measurement.

Text in the SVG needs a `font-family` that XTS knows, on the `<text>`
element or on an enclosing `<g>` or the `<svg>` root. A text whose font
is unknown is dropped with a warning in the protocol.
