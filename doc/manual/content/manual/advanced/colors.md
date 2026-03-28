---
weight: 40
type: docs
linktitle: Colors
---

# Colors

## Using colors

Colors can be used anywhere a `color` or `backgroundcolor` attribute is accepted:

```xml
<Box width="4" height="2" backgroundcolor="limegreen"/>
<Paragraph style="color: #336699;">...</Paragraph>
<Circle radiusx="2" backgroundcolor="rgb(255, 128, 0)"/>
```

## CSS color values

XTS supports the standard CSS color formats:

- Named colors: `red`, `darkblue`, `limegreen`, ...
- Hex: `#ff0000`, `#369`
- RGB: `rgb(255, 0, 0)`

## Defining custom colors

Create named colors with `<DefineColor>`:

```xml
<DefineColor name="brandblue" value="#1a73e8"/>
<DefineColor name="brandgray" value="rgb(100, 100, 100)"/>
```

Then use them by name:

```xml
<Box width="4" height="2" backgroundcolor="brandblue"/>
```

## Pre-defined colors

XTS comes with all standard CSS named colors plus:

- **HKS 1-97** -- defined in CMYK
- **Many Pantone colors** -- defined in CMYK

The colors `black` and `white` are in the grayscale color space. All other CSS colors are in RGB.

## See also

- [DefineColor reference](/reference/commands/definecolor)
- [Defaults](/reference/defaults)
