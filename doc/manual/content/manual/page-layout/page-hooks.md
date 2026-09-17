---
weight: 20
type: docs
linktitle: Page Hooks
---

# Page Hooks

XTS provides two hooks that let you run commands at specific points in a page's lifecycle:

## AtPageCreation

Commands inside `<AtPageCreation>` run when a new page is created, *before* any content is placed. This is the right place for:

- Background images or stationery
- Watermarks
- Repeating page headers or footers

```xml
<DefineMasterPage name="default" test="true()" margin="1cm">
    <AtPageCreation>
        <PlaceObject row="1" column="1" allocate="no">
            <Image href="letterhead.pdf" width="210mm" height="297mm"/>
        </PlaceObject>
    </AtPageCreation>
</DefineMasterPage>
```

## AtPageShipout

Commands inside `<AtPageShipout>` run when a page is finalized and written to the PDF. This is the place for:

- Page numbers (you now know the final page number)
- Running headers with chapter titles
- Content that depends on what's on the page

```xml
<DefineMasterPage name="default" test="true()" margin="1cm">
    <AtPageShipout>
        <PlaceObject column="{sd:number-of-columns()}" row="{sd:number-of-rows()}"
            hreference="right" halign="right">
            <TextBlock>
                <Paragraph>
                    <Value select="sd:current-page()"/>
                </Paragraph>
            </TextBlock>
        </PlaceObject>
    </AtPageShipout>
</DefineMasterPage>
```

## Combining both

A typical setup uses `<AtPageCreation>` for the background and `<AtPageShipout>` for page-dependent content:

```xml
<DefineMasterPage name="standard" test="true()" margin="2cm 1cm 1cm 1cm">
    <AtPageCreation>
        <!-- Company logo in the top-right corner -->
        <PlaceObject column="{sd:number-of-columns()}" row="1"
            hreference="right" allocate="no">
            <Image href="logo.pdf" width="3cm"/>
        </PlaceObject>
    </AtPageCreation>
    <AtPageShipout>
        <!-- Page number at the bottom center -->
        <PlaceObject column="1" row="{sd:number-of-rows()}" allocate="no">
            <TextBlock>
                <Paragraph style="text-align: center;">
                    <Value select="sd:current-page()"/>
                </Paragraph>
            </TextBlock>
        </PlaceObject>
    </AtPageShipout>
</DefineMasterPage>
```

## Page margin boxes from CSS

For static headers and footers there is a third option that needs no hook at all: the page margin boxes of a CSS `@page` rule in a [StyleSheet](/reference/commands/stylesheet). The rule `@page name { ... }` belongs to the master page with the same `name`; the generic `@page { ... }` rule is the base for every master page.

```xml
<StyleSheet>
    @page default {
        margin: 2cm;
        @top-left { content: "NORDWERK"; font-weight: bold; vertical-align: bottom; }
        @bottom-right { content: "Page " counter(page); font-size: 9pt; }
    }
</StyleSheet>
<DefineMasterPage name="default" test="true()" />
```

The margin boxes differ from the hooks in three ways:

- They are placed in the page margin, outside the grid, so they never occupy grid cells. A header in `<AtPageCreation>` at row 1 takes that row away from the content.
- They are rendered when the page is written to the PDF, like `<AtPageShipout>`, so `counter(page)` is the final page number.
- Their content is limited to text, `counter(page)` and an image via `url()`. Tables, several paragraphs or values from the data still need a hook.

The `margin` attribute of `<DefineMasterPage>` stays authoritative for the page geometry, because it defines the grid. Without the attribute the margins come from the `@page` rule. If both are given and differ, the attribute wins and XTS issues a warning. The CSS page selectors `:first`, `:left` and `:right` are not evaluated, the `test` attribute selects the master page.

## See also

- [AtPageCreation reference](/reference/commands/atpagecreation)
- [AtPageShipout reference](/reference/commands/atpageshipout)
