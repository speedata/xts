# Loop



Repeat the contents of this element several times.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Bookmark](../bookmark), [Br](../br), [CallTemplate](../calltemplate), [ClearPage](../clearpage), [Column](../column), [CopyOf](../copyof), [Flow](../flow), [ForAll](../forall), [HTML](../html), [I](../i), [Li](../li), [LoadXML](../loadxml), [Loop](../loop), [Mark](../mark), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Ol](../ol), [Paragraph](../paragraph), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Span](../span), [Switch](../switch), [Tab](../tab), [Table](../table), [Td](../td), [Tr](../tr), [U](../u), [Ul](../ul), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[A](../a), [AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [B](../b), [Case](../case), [Columns](../columns), [Contents](../contents), [Flow](../flow), [ForAll](../forall), [Function](../function), [I](../i), [Li](../li), [Loop](../loop), [Otherwise](../otherwise), [Paragraph](../paragraph), [Record](../record), [SaveXML](../savexml), [Span](../span), [Table](../table), [Td](../td), [Template](../template), [Tr](../tr), [U](../u), [Until](../until), [While](../while)


## Attributes



`select` ([XPath expressions](/programming/xpath))
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






## See also

- Commands: [While](../while), [Until](../until), [ForAll](../forall)
- Manual: [Control flow](/programming/control-flow)

