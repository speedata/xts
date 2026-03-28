# Action



Associates an action with a text. Once the text is placed on the page, the associated action will be executed. The action can be compared to an invisible character. When the publisher outputs the character, the corresponding instructions will be run.



##  Child elements

[Mark](../mark)

##  Parent elements

[A](../a), [B](../b), [I](../i), [Li](../li), [Paragraph](../paragraph), [Span](../span), [Textblock](../textblock), [U](../u)


## Attributes
(none)

## Example

```xml
<Pageformat width="210mm" height="4cm"/>

<Record element="data">
  <PlaceObject>
    <Textblock>
      <Paragraph>
        <Value>
          Row
          Row
          Row
          Row
        </Value>
      </Paragraph>
    </Textblock>
    <Textblock>
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
    </Textblock>
  </PlaceObject>
  <ClearPage/>
  <Message select="sd:pagenumber('textstart')"></Message>
</Record>

```





