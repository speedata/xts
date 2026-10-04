# Flow



Pours a sequence of blocks into the frames of a positioning area, page after page. Paragraphs break between lines, with the CSS properties `widows` and `orphans`, and tables between rows, their [TableHead](../tablehead) and [TableFoot](../tablefoot) repeated. `break-after: avoid` keeps a block with the next one, `break-before` and `break-after` with `column` go on in the next frame and with `page`, `left` or `right` on the next page.

The flow starts at the area's current row and fills the free rows of each frame: rows that hold objects placed earlier are passed over. The blocks are not rounded to the grid; the rows they reach into are allocated. A free band does not start on the next whole row but at the exact bottom of the objects that end in the row above it, so where the text starts does not depend on the row height. A [PlaceObject](../placeobject) after the flow starts on the next whole row. A [Flow](../flow) that follows another directly, with nothing placed in between, goes on exactly where it ended, its last margin collapsing with the first block's, and so does every frame of its area that is free from there on the same page.

Every child with an `id` keeps it, also when it breaks across frames and pages: each part shows up in `--dumpoutput` with its page and position.

A [TextBlock](../textblock), an [Image](../image) or a [Box](../box) cannot be part of a flow yet and is left out with a warning. A flow is not possible in a [Slate](../slate).

The children are read before the flow is laid out, so a command that acts on the page at once, [PlaceObject](../placeobject), [ClearPage](../clearpage), [NextFrame](../nextframe), [NextRow](../nextrow) or [Mark](../mark), is an error inside a flow. A forced break goes through CSS, such as `break-before: page`. [CallTemplate](../calltemplate) works as everywhere.



##  Child elements

[ForAll](../forall), [HTML](../html), [Ol](../ol), [Paragraph](../paragraph), [Switch](../switch), [Table](../table), [Ul](../ul)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Template](../template), [Until](../until), [While](../while)


## Attributes



`area` (text, optional)
:   Name of the positioning area to fill. Defaults to the page's area.




`bottom` (text, optional)
:   Name of a variable that is set to where the flow ends, in points below the top of the page's grid, as a number. With frames side by side, that is the bottom of the longest one on the last page.




## Example

```xml
<Record match="data">
  <Flow area="text">
    <Paragraph id="intro" style="break-after: avoid"><Value>Introduction</Value></Paragraph>
    <ForAll select="paragraph">
      <Paragraph><Value select="."/></Paragraph>
    </ForAll>
    <Table id="prices">…</Table>
  </Flow>
</Record>
```

Pours a heading, the paragraphs and a table into the frames of the area `text`, starting new pages as needed.








## See also

- Commands: [PlaceObject](../placeobject), [PositioningArea](../positioningarea), [Table](../table), [Paragraph](../paragraph)
- Manual: [Flowing Text](/manual/core-concepts/flow), [Multi-Page Content](/manual/page-layout/multi-page), [Working with Tables: Tables across pages](/manual/tables/tables#tables-across-pages)

