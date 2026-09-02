# PositioningArea



Describes an area which contains one or more frames. Elements can be placed within these frames.

The frames can be wrapped in a [Switch](../switch) element. The Switch is evaluated each time a new page is created, so a Case test can use `sd:current-page()` to select different frames for even and odd pages. If no branch matches, the area does not exist on that page.



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
```xml
<PositioningArea name="text">
  <Switch>
    <Case test="sd:odd(sd:current-page())">
      <PositioningFrame width="12" height="30" column="2" row="2"/>
    </Case>
    <Otherwise>
      <PositioningFrame width="12" height="30" column="16" row="2"/>
    </Otherwise>
  </Switch>
</PositioningArea>
```





