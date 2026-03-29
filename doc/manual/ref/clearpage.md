# ClearPage



Finishes the current page. 



##  Child elements

(none)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Function](../function), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Until](../until), [While](../while)


## Attributes
(none)

## Example

```xml
<Record match="data">
  <PlaceObject>
    <TextBlock>
      <Paragraph><Value>This is page 1</Value></Paragraph>
    </TextBlock>
  </PlaceObject>
  <ClearPage openon="right"/>
  <PlaceObject>
    <TextBlock>
      <Paragraph><Value>And this is page 3</Value></Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>

```





