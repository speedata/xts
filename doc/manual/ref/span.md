# Span



Surround text by styling options.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Bookmark](../bookmark), [Br](../br), [CopyOf](../copyof), [ForAll](../forall), [HTML](../html), [I](../i), [Loop](../loop), [Ol](../ol), [Span](../span), [Switch](../switch), [Tab](../tab), [U](../u), [Ul](../ul), [Value](../value)

##  Parent elements

[A](../a), [B](../b), [Case](../case), [ForAll](../forall), [I](../i), [Li](../li), [Loop](../loop), [Otherwise](../otherwise), [Paragraph](../paragraph), [Span](../span), [U](../u), [Until](../until), [While](../while)


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

