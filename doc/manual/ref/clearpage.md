# ClearPage



Finishes the current page. 



##  Child elements

(none)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Function](../function), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Template](../template), [Until](../until), [While](../while)


## Attributes



`openon` (optional)
:   Makes sure that the next page opens on the given side by inserting a blank page if necessary. Right pages have odd page numbers.



    `right`
    :    The next page is a right page (odd page number).



    `left`
    :    The next page is a left page (even page number).




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





