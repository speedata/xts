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
<Record element="data">
  <PlaceObject>
    <Textblock>
      <Paragraph><Value>This is page 1</Value></Paragraph>
    </Textblock>
  </PlaceObject>
  <ClearPage openon="right"/>
  <PlaceObject>
    <Textblock>
      <Paragraph><Value>And this is page 3</Value></Paragraph>
    </Textblock>
  </PlaceObject>
</Record>

```





