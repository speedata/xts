# Span



Surround text by styling options.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Br](../br), [CopyOf](../copyof), [HTML](../html), [I](../i), [Ol](../ol), [Span](../span), [U](../u), [Ul](../ul), [Value](../value)

##  Parent elements

[A](../a), [B](../b), [I](../i), [Li](../li), [Paragraph](../paragraph), [Span](../span), [U](../u)


## Attributes



`class` (text, optional)
:   CSS class for this element.




`id` (text, optional)
:   CSS id for this element.




`style` (text, optional)
:   Set the CSS style.




## Example

```xml
<Stylesheet>
  .green { background-color: lightgreen; }
</Stylesheet>

<Record element="data">
  <PlaceObject>
    <Textblock>
      <Paragraph>
        <Span class="green"><Value>green</Value></Span>
      </Paragraph>
    </Textblock>
  </PlaceObject>
</Record>

```





