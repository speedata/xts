# A



Insert hyperlink to a URL.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Bookmark](../bookmark), [Br](../br), [CopyOf](../copyof), [ForAll](../forall), [HTML](../html), [I](../i), [Loop](../loop), [Ol](../ol), [Span](../span), [Switch](../switch), [Tab](../tab), [U](../u), [Ul](../ul), [Value](../value)

##  Parent elements

[A](../a), [B](../b), [Case](../case), [ForAll](../forall), [I](../i), [Li](../li), [Loop](../loop), [Otherwise](../otherwise), [Paragraph](../paragraph), [Span](../span), [U](../u), [Until](../until), [While](../while)


## Attributes



`href` (text, optional)
:   The target of the hyperlink (URI). Example: `https://www.speedata.de`




`link` (text, optional)
:   The target of the document link (a Mark). Example: `article123`




`page` (number, optional)
:   The target page number. Each page automatically gets a destination, so you can link to a specific page. Example: `2`




## Example

```xml
<PlaceObject>
  <TextBlock>
    <Paragraph><Value>See the </Value>
      <A href="https://www.speedata.de">
        <Value>homepage</Value>
      </A>
      <Value> for more information.</Value>
    </Paragraph>
  </TextBlock>
</PlaceObject>
```






## See also

- Commands: [Action](../action), [Mark](../mark), [Bookmark](../bookmark)
- Manual: [PDF Options: Links](/manual/advanced/pdf-options#links)

