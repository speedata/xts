# A



Insert hyperlink to a URL.



##  Child elements

[A](../a), [Action](../action), [B](../b), [Br](../br), [CopyOf](../copyof), [HTML](../html), [I](../i), [Ol](../ol), [Span](../span), [U](../u), [Ul](../ul), [Value](../value)

##  Parent elements

[A](../a), [B](../b), [I](../i), [Li](../li), [Paragraph](../paragraph), [Span](../span), [U](../u)


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





