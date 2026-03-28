# While



Create a loop. All child elements are executed as long as the condition in the test attribute evaluates to true.



##  Child elements

[ClearPage](../clearpage), [Column](../column), [ForAll](../forall), [HTML](../html), [Li](../li), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Paragraph](../paragraph), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Switch](../switch), [Td](../td), [Tr](../tr), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [SetVariable](../setvariable), [Until](../until), [While](../while)


## Attributes



`test` ([XPath expressions](/manual/data-processing/xpath))
:   Every time before the the loop is executed, this condition must evaluate to true. See the command [Until](../until) for a loop with an exit test.




## Example


The following example creates a textblock with three times the contents 'Text Text Text '.


```xml

<Record element="data">
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
        <Textblock>
            <Paragraph><Value select="$text"/></Paragraph>
        </Textblock>
    </PlaceObject>
</Record>

```





