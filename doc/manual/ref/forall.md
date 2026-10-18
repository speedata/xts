# ForAll



Executes the given commands for all elements in the data XML file that match the contents of the attribute `select`.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Bookmark](../bookmark), [Br](../br), [CallTemplate](../calltemplate), [ClearPage](../clearpage), [Column](../column), [CopyOf](../copyof), [Flow](../flow), [ForAll](../forall), [HTML](../html), [I](../i), [Li](../li), [LoadXML](../loadxml), [Loop](../loop), [Mark](../mark), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Ol](../ol), [Paragraph](../paragraph), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Span](../span), [Switch](../switch), [Tab](../tab), [Table](../table), [Td](../td), [Tr](../tr), [U](../u), [Ul](../ul), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[A](../a), [AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [B](../b), [Case](../case), [Columns](../columns), [Contents](../contents), [DefineMasterPage](../definemasterpage), [Flow](../flow), [ForAll](../forall), [Function](../function), [I](../i), [Li](../li), [Loop](../loop), [Ol](../ol), [Otherwise](../otherwise), [Paragraph](../paragraph), [Record](../record), [SaveXML](../savexml), [SetVariable](../setvariable), [Span](../span), [Table](../table), [TableFoot](../tablefoot), [TableHead](../tablehead), [Td](../td), [Template](../template), [TextBlock](../textblock), [Tr](../tr), [U](../u), [Ul](../ul), [Until](../until), [While](../while)


## Attributes



`select` ([XPath expressions](/programming/xpath))
:   Selects the child elements from the data XML




## Example

```xml
<Record match="data">
  <PlaceObject>
    <Table>
      <ForAll select="entry">
        <Tr><Td><Paragraph><Value select="string(.)"/></Paragraph></Td></Tr>
      </ForAll>
    </Table>
  </PlaceObject>
</Record>
```

Creates a table row for all elements `entry` in the data element `data`. The data XML should look similar to this:


```xml
<data>
  <entry>a</entry>
  <entry>b</entry>
  <entry>c</entry>
</data>
```






## See also

- Commands: [ProcessNode](../processnode), [Loop](../loop), [While](../while)
- Manual: [Control flow](/programming/control-flow), [Records and dispatch: ProcessNode versus ForAll](/programming/records-and-dispatch#processnode-versus-forall)

