# TextBlock



Create a rectangular piece of text.



##  Child elements

[Action](../action), [Bookmark](../bookmark), [ForAll](../forall), [HTML](../html), [Ol](../ol), [Paragraph](../paragraph), [Ul](../ul), [Value](../value)

##  Parent elements

[PlaceObject](../placeobject), [SetVariable](../setvariable)


## Attributes



`parsep` (length, optional)
:   The vertical distance between two paragraphs.




`width` (number, optional)
:   Number of columns for the text. If not given, the surrounding element determines the width of the element.




## Example

```xml
<Record match="data">
  <PlaceObject>
    <TextBlock width="10">
      <Paragraph>
        <B><Value>Bold text</Value></B>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>

```





