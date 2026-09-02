# Grid



Set the grid of the surrounding [Slate](../slate). Settings not given here keep the value copied from the page grid. The page grid itself is set with [SetGrid](../setgrid).



##  Child elements

(none)

##  Parent elements

[Slate](../slate)


## Attributes



`dx` (optional)
:   Gap between two grid cells (horizontal).




`dy` (optional)
:   Gap between two grid cells (vertical).




`height` (length, optional)
:   The height of a grid cell in the slate.




`nx` (number, optional)
:   The number of grid cells of the slate in horizontal direction. It limits the area for the cursor when objects are placed without a column.




`ny` (number, optional)
:   The number of grid cells of the slate in vertical direction. It limits the area for the cursor when objects are placed without a row.




`width` (length, optional)
:   The width of a grid cell in the slate.




## Example

```xml
<Slate name="calendar">
  <Grid width="5mm" height="5mm"/>
  <Contents>
    <PlaceObject column="3" row="2">
      <TextBlock width="10">
        <Paragraph><Value>Positioned on the 5mm slate grid</Value></Paragraph>
      </TextBlock>
    </PlaceObject>
  </Contents>
</Slate>
```





