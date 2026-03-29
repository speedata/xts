# PositioningArea



Describes an area which contains one or more frames. Elements can be placed within these frames.



##  Child elements

[PositioningFrame](../positioningframe), [Switch](../switch)

##  Parent elements

[DefineMasterPage](../definemasterpage)


## Attributes



`framecolor` (text, optional)
:   Set the color of the frame in grid=yes mode. Defaults to 'red'




`name` (text)
:   Name of the area.




## Example

```xml
<Masterpage name="right page" test="sd:odd( sd:current-page() )">
  <Margin left="1cm" right="1cm" top="1cm" bottom="1cm"/>
  <PositioningArea name="frame1">
    <PositioningFrame width="12" height="30" column="2" row="2"/>
    <PositioningFrame width="12" height="30" column="16" row="2"/>
  </PositioningArea>
</Masterpage>
```





