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
- HSL: `hsl(120 100% 25%)`, `hsla(120, 100%, 25%, 0.5)`
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

## Overprinting

By default a color knocks out what is below it: the plates of the other colors stay empty in that area. Small black text on a colored background then depends on the exact fit of the plates, and a slight misregister shows a white edge along the letters. An overprinting color prints on top instead. Set `overprint` on the color, in `<DefineColor>` or in `@-bag-color`:

```xml
<DefineColor name="ink" model="cmyk" c="0" m="0" y="0" k="100" overprint="yes"/>
<StyleSheet>
    @-bag-color pms300 { model: spotcolor; colorname: "PANTONE 300 C"; overprint: yes; }
</StyleSheet>
```

The color overprints wherever its name is used, in text, backgrounds and borders. Three things to keep in mind:

- It works with CMYK and spot colors. An RGB or gray color knocks out, so `black` does not overprint, a CMYK color with `k="100"` does. xts warns about this.
- A CMYK color prints only on the plates of its nonzero components. White therefore leaves no mark at all, which xts warns about, too.
- A PDF viewer shows overprinting only with overprint preview switched on. To check the plates, separate them, for example with Ghostscript: `gs -sDEVICE=tiffsep -o plates.tif document.pdf`.

## Pre-defined colors

XTS comes with all standard CSS named colors plus:

- **HKS** (86 colors) -- spot colors with CMYK alternate values
- **Many Pantone colors** -- spot colors with CMYK alternate values

The colors `black` and `white` are in the grayscale color space. All other CSS colors are in RGB.

## See also

- [DefineColor reference](/reference/commands/definecolor)
- [Defaults](/reference/defaults)
