# Until



Create a loop. All child elements are executed repeatedly until the given condition is true. The return value of Until is the concatenated return value of the child elements.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Bookmark](../bookmark), [Br](../br), [CallTemplate](../calltemplate), [ClearPage](../clearpage), [Column](../column), [CopyOf](../copyof), [Flow](../flow), [ForAll](../forall), [HTML](../html), [I](../i), [Li](../li), [LoadXML](../loadxml), [Loop](../loop), [Mark](../mark), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Ol](../ol), [Paragraph](../paragraph), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Span](../span), [Switch](../switch), [Tab](../tab), [Table](../table), [Td](../td), [Tr](../tr), [U](../u), [Ul](../ul), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [Flow](../flow), [ForAll](../forall), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [SetVariable](../setvariable), [Template](../template), [Until](../until), [While](../while)


## Attributes



`test` ([XPath expressions](/programming/xpath))
:   Every time after the loop is executed, the condition is evaluated. If it is true, the loop exits.




## Example

```xml
<Record match="data">
  <SetVariable variable="i" select="0"/>
  <Until test="$i = 4">
    <Message select="concat('$i is: ', $i)"/>
    <SetVariable variable="i" select="$i + 1"/>
  </Until>
</Record>

```

Gives the following output (in the protocol file)


```xml
Message: "$i is: 0"
Message: "$i is: 1"
Message: "$i is: 2"
Message: "$i is: 3"
```






## See also

- Commands: [While](../while), [Loop](../loop)
- Manual: [Control flow](/programming/control-flow)

