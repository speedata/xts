# While



Create a loop. All child elements are executed as long as the condition in the test attribute evaluates to true.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Bookmark](../bookmark), [Br](../br), [CallTemplate](../calltemplate), [ClearPage](../clearpage), [Column](../column), [CopyOf](../copyof), [Flow](../flow), [ForAll](../forall), [HTML](../html), [I](../i), [Li](../li), [LoadXML](../loadxml), [Loop](../loop), [Mark](../mark), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Ol](../ol), [Paragraph](../paragraph), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Span](../span), [Switch](../switch), [Tab](../tab), [Table](../table), [Td](../td), [Tr](../tr), [U](../u), [Ul](../ul), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [Flow](../flow), [ForAll](../forall), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [SetVariable](../setvariable), [Template](../template), [Until](../until), [While](../while)


## Attributes



`test` ([XPath expressions](/programming/xpath))
:   Every time before the the loop is executed, this condition must evaluate to true. See the command [Until](../until) for a loop with an exit test.




## Example


The following example creates a textblock with three times the contents 'Text Text Text '.


```xml

<Record match="data">
    <SetVariable variable="counter" select="1"/>
    <SetVariable variable="text" select="''"/>
    <While test=" $counter &lt;= 3 "> <!-- less or equal -->
        <SetVariable variable="counter" select=" $counter + 1"/>
        <SetVariable variable="text">
            <Value select="$text"/>
            <Value select="'Text '"/>
      </SetVariable>
    </While>
    <PlaceObject>
        <TextBlock>
            <Paragraph><Value select="$text"/></Paragraph>
        </TextBlock>
    </PlaceObject>
</Record>

```






## See also

- Commands: [Until](../until), [Loop](../loop), [ForAll](../forall)
- Manual: [Control flow](/programming/control-flow)

