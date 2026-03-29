# Mark



Sets an invisible mark into the output. This is helpful when you want to know on which page the mark is placed on.



##  Child elements

(none)

##  Parent elements

[Action](../action)


## Attributes



`append` (optional)
:   When yes, append the current page to the previous values of the mark. Useful to get page ranges in an index. (Default is no.)



    `yes`
    :    Append the page number to the previous value of the mark.



    `no`
    :    Replace the previous value.




`pdftarget` (yes or no, optional)
:   Set a pdf target that can be referenced by [A](../a)




`select` ([XPath expressions](/manual/data-processing/xpath))
:   The name of the mark to be set.




`shiftup` (length, optional)
:   Raise the position of the hyperlink anchor by this amount.




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





## Info

Marks get saved for subsequent runs.




