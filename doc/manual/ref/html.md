# HTML



Insert HTML content, either inline or retrieved via an XPath expression.



##  Child elements

(none)

##  Parent elements

[A](../a), [B](../b), [Case](../case), [ForAll](../forall), [I](../i), [Li](../li), [Loop](../loop), [Otherwise](../otherwise), [Paragraph](../paragraph), [PlaceObject](../placeobject), [Span](../span), [Td](../td), [Textblock](../textblock), [U](../u), [Until](../until), [While](../while)


## Attributes



`expand-text` (yes or no, optional)
:   If set to "yes", expressions in curly braces {expr} within the HTML content are evaluated as XPath expressions (similar to XSLT 3.0 text value templates). Use {{ and }} for literal curly braces. Default is "no".




`select` ([XPath expressions](/manual/data-processing/xpath), optional)
:   XPath expression that yields an XML/HTML fragment.




## Example

```xml
<Record element="data">
  <PlaceObject>
    <Textblock>
      <HTML select="/data/htmlcontent"/>
    </Textblock>
  </PlaceObject>
</Record>

```

```xml
<SetVariable variable="product" select="'Widget'"/>
<PlaceObject>
  <Textblock>
    <HTML expand-text="yes">
      <p>Product: {$product}</p>
      <p>Literal braces: {{ and }}</p>
    </HTML>
  </Textblock>
</PlaceObject>

```





