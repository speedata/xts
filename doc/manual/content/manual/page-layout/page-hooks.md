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
<DefineMasterpage name="default" test="true()" margin="1cm">
    <AtPageCreation>
        <PlaceObject row="1" column="1" allocate="no">
            <Image href="letterhead.pdf" width="210mm" height="297mm"/>
        </PlaceObject>
    </AtPageCreation>
</DefineMasterpage>
```

## AtPageShipout

Commands inside `<AtPageShipout>` run when a page is finalized and written to the PDF. This is the place for:

- Page numbers (you now know the final page number)
- Running headers with chapter titles
- Content that depends on what's on the page

```xml
<DefineMasterpage name="default" test="true()" margin="1cm">
    <AtPageShipout>
        <PlaceObject column="{sd:number-of-columns()}" row="{sd:number-of-rows()}"
            hreference="right" halign="right">
            <Textblock>
                <Paragraph>
                    <Value select="sd:current-page()"/>
                </Paragraph>
            </Textblock>
        </PlaceObject>
    </AtPageShipout>
</DefineMasterpage>
```

## Combining both

A typical setup uses `<AtPageCreation>` for the background and `<AtPageShipout>` for page-dependent content:

```xml
<DefineMasterpage name="standard" test="true()" margin="2cm 1cm 1cm 1cm">
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
            <Textblock>
                <Paragraph style="text-align: center;">
                    <Value select="sd:current-page()"/>
                </Paragraph>
            </Textblock>
        </PlaceObject>
    </AtPageShipout>
</DefineMasterpage>
```

## See also

- [AtPageCreation reference](/reference/commands/atpagecreation)
- [AtPageShipout reference](/reference/commands/atpageshipout)
