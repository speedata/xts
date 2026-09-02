# PositioningArea



Describes an area which contains one or more frames. Elements can be placed within these frames.



##  Child elements

[PositioningFrame](../positioningframe), [Switch](../switch)

##  Parent elements

[DefineMasterPage](../definemasterpage)


## Attributes



`name` (text)
:   Name of the area.




## Example

```xml
<DefineMasterPage name="right page" margin="1cm" test="sd:odd( sd:current-page() )">
  <PositioningArea name="frame1">
    <PositioningFrame width="12" height="30" column="2" row="2"/>
    <PositioningFrame width="12" height="30" column="16" row="2"/>
  </PositioningArea>
</DefineMasterPage>
```





