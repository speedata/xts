---
weight: 30
type: docs
linktitle: Multi-Page Content
---

# Multi-Page Content

## Tables across pages

Tables automatically break across pages. If you define a `<Tablehead>`, it repeats on each new page:

```xml
<Table>
    <Tablehead>
        <Tr><Td><Paragraph><Value>Header</Value></Paragraph></Td></Tr>
    </Tablehead>
    <!-- Hundreds of rows... they'll flow to new pages automatically -->
</Table>
```

Individual table cells are *never* split -- each cell is rendered as a single box. If a row doesn't fit on the current page, it moves to the next one.

## Page breaks

Force a new page with `<ClearPage>`:

```xml
<Record element="catalog">
    <ForAll select="category">
        <ProcessNode select="."/>
        <ClearPage/>
    </ForAll>
</Record>
```

## Frame switching

If an area has multiple frames, use `<NextFrame>` to jump to the next one:

```xml
<NextFrame area="text"/>
```

If there's no next frame, a page break is inserted and content continues in the first frame of the area on the new page.

## Positioning frames for multi-column layouts

Define multiple frames within a single area to create flowing multi-column layouts:

```xml
<DefineMasterpage name="threeColumn" test="true()" margin="1cm">
    <PositioningArea name="text">
        <PositioningFrame width="5" height="20" row="1" column="1"/>
        <PositioningFrame width="5" height="20" row="1" column="7"/>
        <PositioningFrame width="5" height="20" row="1" column="13"/>
    </PositioningArea>
</DefineMasterpage>
```

Content placed in the "text" area fills the first column, then flows to the second, then the third. When all three are full, a new page is created and the cycle repeats.

## See also

- [ClearPage reference](/reference/commands/clearpage)
- [NextFrame reference](/reference/commands/nextframe)
- [PositioningArea reference](/reference/commands/positioningarea)
