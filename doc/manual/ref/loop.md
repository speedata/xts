# Loop



Repeat the contents of this element several times.



##  Child elements

[ClearPage](../clearpage), [Column](../column), [ForAll](../forall), [HTML](../html), [Li](../li), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Paragraph](../paragraph), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Switch](../switch), [Td](../td), [Tr](../tr), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Columns](../columns), [Contents](../contents), [ForAll](../forall), [Function](../function), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [SaveXML](../savexml), [Table](../table), [Td](../td), [Tr](../tr), [Until](../until), [While](../while)


## Attributes



`select` ([XPath expressions](/manual/data-processing/xpath))
:   The number of loops. Must be a number or castable as a number.




`variable` (text, optional)
:   If given, store the current loop value in this variable. If omitted, the loop value is stored in the variable `_loopcounter`.




## Example

```xml
<PlaceObject>
  <Table>
    <Loop select="5" variable="i">
      <Tr>
        <Td><Paragraph><Value select="$i"/></Paragraph></Td>
      </Tr>
    </Loop>
  </Table>
</PlaceObject>

```





