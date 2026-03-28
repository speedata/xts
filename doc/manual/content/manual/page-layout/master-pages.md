---
weight: 10
type: docs
linktitle: Master Pages
---

# Master Pages

A master page defines the blueprint for a page: its margins, grid, and positioning areas. XTS picks the right master page for each new page based on a `test` condition.

## Defining a master page

```xml
<DefineMasterpage name="default" test="true()" margin="1cm"/>
```

This creates a master page called "default" that matches all pages (the test is always true) with 1cm margins on all sides.

## Conditional master pages

You can have different layouts for different pages:

```xml
<!-- First page: extra top margin for a header -->
<DefineMasterpage name="first" test="sd:current-page() = 1" margin="1cm 1cm 3cm 1cm"/>

<!-- Even pages: wider left margin for binding -->
<DefineMasterpage name="even" test="sd:even(sd:current-page())" margin="1cm 2cm 1cm 1cm"/>

<!-- Everything else -->
<DefineMasterpage name="default" test="true()" margin="1cm"/>
```

XTS evaluates the tests **in order** and uses the first match. Put your most specific conditions first and the catch-all `true()` last.

## Master pages with areas

Combine master pages with positioning areas for complex layouts:

```xml
<SetGrid nx="12" ny="20"/>

<DefineMasterpage name="twoColumn" test="true()" margin="1.5cm">
    <PositioningArea name="header">
        <PositioningFrame width="12" height="2" row="1" column="1"/>
    </PositioningArea>
    <PositioningArea name="left">
        <PositioningFrame width="5" height="16" row="4" column="1"/>
    </PositioningArea>
    <PositioningArea name="right">
        <PositioningFrame width="5" height="16" row="4" column="7"/>
    </PositioningArea>
</DefineMasterpage>
```

## Page format

Set the page size with `<Pageformat>`:

```xml
<Pageformat width="210mm" height="297mm"/>  <!-- A4 -->
<Pageformat width="8.5in" height="11in"/>   <!-- US Letter -->
<Pageformat width="15cm" height="20cm"/>    <!-- Custom -->
```

The default is A4 (210mm x 297mm).

## See also

- [DefineMasterpage reference](/reference/commands/definemasterpage)
- [Pageformat reference](/reference/commands/pageformat)
- [Positioning Areas](../../core-concepts/positioning-areas)
