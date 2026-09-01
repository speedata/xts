---
weight: 30
type: docs
linktitle: Multi-Page Content
---

# Multi-Page Content

## Tables across pages

A table with a `<TableHead>` breaks across frames and pages. Its rows are laid
out down the frame, and when the next row does not fit, the table continues in
the next frame of the area, or on a new page once the frames are used up. The
header rows are repeated at the top of every continuation:

```xml
<Table>
    <TableHead>
        <Tr><Td><Paragraph><Value>Header</Value></Paragraph></Td></Tr>
    </TableHead>
    <!-- Hundreds of rows... they'll flow to new frames and pages automatically -->
</Table>
```

Two conditions apply. The table needs a `<TableHead>`: without one there is
nothing to repeat and the table is placed as a single object. And it has to be
placed on the grid, which is the default. A table placed at an absolute
position (`<PlaceObject column="2cm" row="5cm">`) is put where it was asked for
and is not broken up.

Individual table cells are *never* split -- each cell is rendered as a single
box. If a row doesn't fit on the current page, it moves to the next one. A row
taller than an empty frame is placed and allowed to overflow.

## Page breaks

Force a new page with `<ClearPage>`:

```xml
<Record match="catalog">
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
<DefineMasterPage name="threeColumn" test="true()" margin="1cm">
    <PositioningArea name="text">
        <PositioningFrame width="5" height="20" row="1" column="1"/>
        <PositioningFrame width="5" height="20" row="1" column="7"/>
        <PositioningFrame width="5" height="20" row="1" column="13"/>
    </PositioningArea>
</DefineMasterPage>
```

Content placed in the "text" area fills the first column, then flows to the second, then the third. When all three are full, a new page is created and the cycle repeats.

## See also

- [ClearPage reference](/reference/commands/clearpage)
- [NextFrame reference](/reference/commands/nextframe)
- [PositioningArea reference](/reference/commands/positioningarea)
