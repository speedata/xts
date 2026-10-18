# Tab



Insert a tab. The text after it continues at the next tab stop set with the CSS property `-bag-tab-stops`. Without tab stops, the tab advances by `tab-size`. Unlike a tab character in the text, this command is never collapsed to a space.



##  Child elements

(none)

##  Parent elements

[A](../a), [B](../b), [Case](../case), [ForAll](../forall), [I](../i), [Li](../li), [Loop](../loop), [Otherwise](../otherwise), [Paragraph](../paragraph), [Span](../span), [U](../u), [Until](../until), [While](../while)


## Attributes
(none)

## Example

```xml
<StyleSheet>
  p.toc { -bag-tab-stops: 8mm, 100% end leader(dotted) }
</StyleSheet>

<Record match="data">
  <PlaceObject>
    <TextBlock>
      <Paragraph class="toc">
        <Value>1</Value><Tab/><Value>Introduction</Value><Tab/><Value>3</Value>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>

```






## See also

- Commands: [Br](../br), [Paragraph](../paragraph)
- Manual: [Tab Stops](/manual/text-and-styling/tab-stops)

