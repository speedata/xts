---
weight: 10
type: docs
linktitle: Images
---

# Images

Including images in XTS is straightforward. The `<Image>` command supports **PDF**, **JPEG**, **PNG**, and **SVG** formats.

## Basic usage

```xml
<PlaceObject>
    <Image href="_samplea.pdf" width="5cm"/>
</PlaceObject>
```

The `href` attribute points to the image file. Width can be in absolute units or grid cells.

## Sizing

You can control dimensions in several ways:

```xml
<!-- Fixed width, height calculated from aspect ratio -->
<Image href="photo.jpg" width="8cm"/>

<!-- Fixed height -->
<Image href="photo.jpg" height="5cm"/>

<!-- Both width and height (may distort) -->
<Image href="photo.jpg" width="8cm" height="5cm"/>

<!-- Width in grid cells -->
<Image href="photo.jpg" width="4"/>

<!-- Minimum / maximum constraints -->
<Image href="photo.jpg" min-width="3cm" max-width="10cm"/>
```

Use `stretch="yes"` together with `max-width` and `max-height` to enlarge the image up to those limits while maintaining the aspect ratio.

## Multi-page PDFs

When including a PDF file, you can select a specific page:

```xml
<Image href="document.pdf" page="3" width="10cm"/>
```

You can also choose which PDF box determines the visible part of the page (`mediabox`, `cropbox`, `bleedbox`, `trimbox`, or `artbox`); the natural size for the width/height computation always comes from the media box:

```xml
<Image href="document.pdf" visiblebox="cropbox" width="10cm"/>
```

## Image dimensions in XPath

Sometimes the layout has to react to the image instead of the other way round: a landscape photo gets the full width, a portrait photo sits next to its caption, or a box is drawn exactly as tall as the picture. Three XPath functions read the dimensions of an image file before it is placed:

```xml
<!-- Width divided by height, 1.5 for a 3:2 photo -->
<Value select="sd:aspect-ratio('photo.jpg')"/>

<!-- Width and height in a unit of your choice -->
<Value select="sd:image-width('photo.jpg', 'cm')"/>
<Value select="sd:image-height('photo.jpg', 'cm')"/>

<!-- Without a unit: the number of grid cells the image would occupy -->
<Value select="sd:image-width('photo.jpg')"/>

<!-- PDF: page 2, measured by its crop box -->
<Value select="sd:image-height('document.pdf', 2, 'cropbox', 'mm')"/>
```

The first argument is the file name. The remaining arguments are optional and can be given in any order: a page number and a box name for PDF files, and a unit such as `'mm'`, `'cm'`, `'pt'` or `'in'`. Without a unit the result is in grid cells, which is what `width` and `height` of `<Image>` expect, so the value can go straight back into the layout:

```xml
<Switch>
  <Case test="sd:aspect-ratio(@src) &gt; 1">
    <!-- landscape: full width of the frame -->
    <PlaceObject><Image href="{@src}" width="{sd:number-of-columns()}"/></PlaceObject>
  </Case>
  <Otherwise>
    <!-- portrait: half the width, caption goes to the right -->
    <PlaceObject><Image href="{@src}" width="{sd:number-of-columns() idiv 2}"/></PlaceObject>
  </Otherwise>
</Switch>
```

See the [XPath functions](/reference/xpath-functions) reference for the exact signatures.

## File locations

Images follow the same file lookup rules as everything else in XTS -- see [File Organization](../../running-xts/file-organization). You can use relative paths, absolute paths, URLs, or the `--extradir` lookup:

```xml
<!-- Relative path -->
<Image href="img/photo.jpg" width="5cm"/>

<!-- URL -->
<Image href="https://example.com/photo.jpg" width="5cm"/>

<!-- Just the filename (if found via extradir) -->
<Image href="photo.jpg" width="5cm"/>
```

## Placeholder images

During development you often need images that don't exist yet. Instead of hunting for dummy files, use the built-in `placeholder://` scheme:

```xml
<Image href="placeholder://200x150" />
<Image href="placeholder://400x200" width="10cm"/>
```

The dimensions after `placeholder://` are in PDF points (WxH). XTS generates a simple placeholder graphic on the fly -- gray background, diagonal cross, a colored circle, and the dimensions as text. No files on disk needed.

## See also

- [Image reference](/reference/commands/image)
