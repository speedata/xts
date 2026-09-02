---
weight: 40
type: docs
linktitle: Colors
---

# Colors

## Using colors

Colors can be used anywhere a `color` or `background-color` attribute is accepted:

```xml
<Box width="4" height="2" background-color="limegreen"/>
<Paragraph style="color: #336699;">...</Paragraph>
<Circle radius-x="2" background-color="rgb(255, 128, 0)"/>
```

## CSS color values

XTS supports the standard CSS color formats:

- Named colors: `red`, `darkblue`, `limegreen`, ...
- Hex: `#ff0000`, `#369`
- RGB: `rgb(255, 0, 0)`, `rgba(255, 0, 0, 0.5)`
- CMYK: `cmyk(0%, 20%, 100%, 5%)` or `device-cmyk(0, 0.2, 1, 0.05)`

## Defining custom colors

Create named colors with `<DefineColor>`:

```xml
<DefineColor name="brandblue" value="#1a73e8"/>
<DefineColor name="brandgray" value="rgb(100, 100, 100)"/>
```

Then use them by name:

```xml
<Box width="4" height="2" background-color="brandblue"/>
```

## Pre-defined colors

XTS comes with all standard CSS named colors plus:

- **HKS** (86 colors) -- spot colors with CMYK alternate values
- **Many Pantone colors** -- spot colors with CMYK alternate values

The colors `black` and `white` are in the grayscale color space. All other CSS colors are in RGB.

## See also

- [DefineColor reference](/reference/commands/definecolor)
- [Defaults](/reference/defaults)
