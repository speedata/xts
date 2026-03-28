# Until



Create a loop. All child elements are executed repeatedly until the given condition is true. The return value of Until is the concatenated return value of the child elements.



##  Child elements

[ClearPage](../clearpage), [Column](../column), [ForAll](../forall), [HTML](../html), [Li](../li), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Paragraph](../paragraph), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Switch](../switch), [Td](../td), [Tr](../tr), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [SetVariable](../setvariable), [Until](../until), [While](../while)


## Attributes



`test` ([XPath expressions](/manual/data-processing/xpath))
:   Every time after the the loop is executed, the condition is evaluated. If it is true, the loop exits.




## Example

```xml
<Record element="data">
  <SetVariable variable="i" select="0"/>
  <Until test="$i = 4">
    <Message select="concat('$i is: ', $i)"/>
    <SetVariable variable="i" select="$i + 1"/>
  </Until>
</Record>

```

Gives the following output (in the protocol file)


```xml
Message: "$i is: 0.000000"
Message: "$i is: 1.000000"
Message: "$i is: 2.000000"
Message: "$i is: 3.000000"
```





