# Value



Contains a text value that is passed to the surrounding element (always as plain text).



##  Child elements

(none)

##  Parent elements

[A](../a), [AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [B](../b), [Case](../case), [Columns](../columns), [Contents](../contents), [ForAll](../forall), [Function](../function), [I](../i), [Li](../li), [Loop](../loop), [Message](../message), [Otherwise](../otherwise), [Paragraph](../paragraph), [PlaceObject](../placeobject), [Record](../record), [SaveXML](../savexml), [SetVariable](../setvariable), [Span](../span), [Table](../table), [Td](../td), [Textblock](../textblock), [Tr](../tr), [U](../u), [Until](../until), [While](../while)


## Attributes



`select` ([XPath expressions](/manual/data-processing/xpath), optional)
:   Value to be passed to the outer element.




## Remarks
The value can be passed to the outer element either as an XPath expression or as the contents of this element.

The result is always treated as text (markup is not preserved).


## Example

```xml
<Record element="data">
  <PlaceObject>
   <Textblock>
     <Paragraph>
       <Value select="@name"/>
       <Value>, symbol=</Value>
       <Value select="@symbol"/>
     </Paragraph>
   </Textblock>
  </PlaceObject>
</Record>

```

hich is the same as


```xml
<Record element="data">
  <PlaceObject>
   <Textblock>
     <Paragraph>
       <Value select="concat(@name, ', symbol=', @symbol)"/>
     </Paragraph>
   </Textblock>
  </PlaceObject>
</Record>

```





