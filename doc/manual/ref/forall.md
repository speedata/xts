# ForAll



Executes the given commands for all elements in the data XML file that match the contents of the attribute `select`.



##  Child elements

[ClearPage](../clearpage), [Column](../column), [ForAll](../forall), [HTML](../html), [Li](../li), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Paragraph](../paragraph), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Switch](../switch), [Td](../td), [Tr](../tr), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Columns](../columns), [Contents](../contents), [DefineMasterPage](../definemasterpage), [ForAll](../forall), [Function](../function), [Loop](../loop), [Ol](../ol), [Otherwise](../otherwise), [Record](../record), [SaveXML](../savexml), [Table](../table), [TableHead](../tablehead), [Td](../td), [TextBlock](../textblock), [Tr](../tr), [Ul](../ul), [Until](../until), [While](../while)


## Attributes



`select` ([XPath expressions](/manual/data-processing/xpath))
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





