# Slate



Create a virtual page (a slate) that behaves like a real page but is not placed into the PDF. Think of it as a magic slate — you draw on it, inspect the result, and then place it on the page.



##  Child elements

[Contents](../contents)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Contents](../contents), [Function](../function), [Record](../record)


## Attributes



`name` (text)
:   Name of the slate that is created.




## Example

```xml
<Record element="data">
  <Slate name="sidebar">
    <!-- Optional, taken from the current page -->
    <Grid width="10mm" height="10mm"/>
    <Contents>
      <PlaceObject column="3" row="2">
        <Textblock width="14">
          <Paragraph>
            <Value>Text</Value>
          </Paragraph>
        </Textblock>
      </PlaceObject>
      <PlaceObject column="2" row="4">
        <Textblock width="14">
          <Paragraph>
            <Value>Next text</Value>
          </Paragraph>
        </Textblock>
      </PlaceObject>
    </Contents>
  </Slate>
  <PlaceObject slate="sidebar" row="1" />
</Record>

```





