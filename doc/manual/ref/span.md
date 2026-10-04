# Span



Surround text by styling options.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Br](../br), [CopyOf](../copyof), [HTML](../html), [I](../i), [Ol](../ol), [Span](../span), [Tab](../tab), [U](../u), [Ul](../ul), [Value](../value)

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
<StyleSheet>
  .green { background-color: lightgreen; }
</StyleSheet>

<Record match="data">
  <PlaceObject>
    <TextBlock>
      <Paragraph>
        <Span class="green"><Value>green</Value></Span>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>

```






## See also

- Commands: [B](../b), [I](../i), [U](../u), [Paragraph](../paragraph)
- Manual: [Text Formatting: Spans](/manual/text-and-styling/text-formatting#spans)

