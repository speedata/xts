---
type: docs
linktitle: Units
---

# Measurement Units

XTS recognizes the following units for length values.

| Unit | Name | Size |
|------|------|------|
| `pt` | Point | 1/72 inch (0.3528 mm) |
| `mm` | Millimeter | 1 mm |
| `cm` | Centimeter | 10 mm |
| `in` | Inch | 25.4 mm |
| `pc` | Pica | 12 pt |
| `pp` | PostScript point | 1/72 inch (same as `pt`) |
| `dd` | Didot point | 0.376 mm |
| `cc` | Cicero | 12 dd |
| `sp` | Scaled point | 1/65536 pt (internal unit) |
| `em` | Em | Relative to current font size |

## Usage in attributes

When an attribute expects a length, include the unit:

```xml
<Pageformat width="210mm" height="297mm"/>
<SetGrid height="12pt" width="5mm"/>
<Image href="photo.jpg" width="8cm"/>
```

## Grid cells vs. lengths

In some attributes (like `row`, `column`, `width` on `<PlaceObject>` and `<Image>`), a plain number means **grid cells**, while a number with a unit means an **absolute length**:

```xml
<!-- 5 grid cells wide -->
<Image href="photo.jpg" width="5"/>

<!-- 5 centimeters wide -->
<Image href="photo.jpg" width="5cm"/>
```

## Unit conversion in XPath

Use `sd:to-unit()` to convert between units:

```xml
<Value select="sd:to-unit('12pt', 'pt', 'mm')"/>
```
