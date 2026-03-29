# Action



Associates an action with a text. Once the text is placed on the page, the associated action will be executed. The action can be compared to an invisible character. When the publisher outputs the character, the corresponding instructions will be run.



##  Child elements

[Mark](../mark)

##  Parent elements

[A](../a), [B](../b), [I](../i), [Li](../li), [Paragraph](../paragraph), [Span](../span), [TextBlock](../textblock), [U](../u)


## Attributes
(none)

## Example

```xml
<PageFormat width="210mm" height="4cm"/>

<Record match="data">
  <PlaceObject>
    <TextBlock>
      <Paragraph>
        <Value>
          Row
          Row
          Row
          Row
        </Value>
      </Paragraph>
    </TextBlock>
    <TextBlock>
      <Action>
        <Mark select="'textstart'"/>
      </Action>
      <Paragraph>
        <Value>
          Row
          Row
          Row
        </Value>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
  <ClearPage/>
  <Message select="sd:pagenumber('textstart')"></Message>
</Record>

```





