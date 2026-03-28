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

The `href` (or `file`) attribute points to the image file. Width can be in absolute units or grid cells.

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
<Image href="photo.jpg" minwidth="3cm" maxwidth="10cm"/>
```

Use `stretch="max"` to fill the available space while maintaining the aspect ratio.

## Multi-page PDFs

When including a PDF file, you can select a specific page:

```xml
<Image href="document.pdf" page="3" width="10cm"/>
```

You can also choose which PDF box to use for sizing: `mediabox`, `cropbox`, `bleedbox`, `trimbox`, or `artbox`:

```xml
<Image href="document.pdf" visiblebox="cropbox" width="10cm"/>
```

## Image dimensions in XPath

Query image dimensions for dynamic layouts:

```xml
<!-- Get the aspect ratio -->
<Value select="sd:aspect-ratio('photo.jpg')"/>

<!-- Get width/height in a specific unit -->
<Value select="sd:image-width('photo.jpg', 1, 'cropbox', 'cm')"/>
<Value select="sd:image-height('photo.jpg', 1, 'cropbox', 'cm')"/>
```

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

## See also

- [Image reference](/reference/commands/image)
