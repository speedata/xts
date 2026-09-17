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

## Defining colors in CSS

The same named colors can be declared in a stylesheet with the `@-bag-color` rule, which keeps color definitions next to the rules that use them. The descriptors mirror the attributes of `<DefineColor>`:

```xml
<StyleSheet>
    @-bag-color gold  { value: #FFC72C; }
    @-bag-color brand { model: cmyk; c: 0; m: 80; y: 90; k: 10; }
    @-bag-color spot  { model: spotcolor; colorname: "PANTONE 300 C"; c: 100; m: 44; y: 0; k: 0; }

    h1 { color: brand; }
    .badge { background-color: gold; }
</StyleSheet>
```

Without a `model` the `value` is any color CSS accepts, including a name defined earlier. The models `cmyk`, `rgb` and `gray` take their components from 0 to 100 (percentages work too), `RGB` and `GRAY` from 0 to 255. A spot color takes the ink name in `colorname` and optional CMYK components for the fallback; without an ink name the color name is used. A color defined this way is available everywhere a color from `<DefineColor>` is, in CSS and in attributes such as `background-color`.

## Pre-defined colors

XTS comes with all standard CSS named colors plus:

- **HKS** (86 colors) -- spot colors with CMYK alternate values
- **Many Pantone colors** -- spot colors with CMYK alternate values

The colors `black` and `white` are in the grayscale color space. All other CSS colors are in RGB.

## See also

- [DefineColor reference](/reference/commands/definecolor)
- [Defaults](/reference/defaults)
