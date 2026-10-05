# Mark



Sets an invisible mark into the output. This is helpful when you want to know on which page the mark is placed on.



##  Child elements

(none)

##  Parent elements

[Action](../action), [Flow](../flow)


## Attributes



`pdftarget` (yes or no, optional)
:   Set a pdf target that can be referenced by [A](../a)




`select` ([XPath expressions](/programming/xpath))
:   The name of the mark to be set.




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
  <Message select="sd:page-number('textstart')"></Message>
</Record>

```





## Info

Marks get saved for subsequent runs.





## See also

- Commands: [Action](../action), [Bookmark](../bookmark)
- Manual: [PDF Options: Page destinations](/manual/advanced/pdf-options#page-destinations), [XPath Functions: Page and position](/reference/xpath-functions#page-and-position)

