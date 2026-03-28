# DefineMasterpage



Define a master page. A master page is chosen depending on the criterion given with the attribute “test”.



##  Child elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [ForAll](../forall), [PositioningArea](../positioningarea)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`margin` (text, optional)
:   Set the margin of the page (defaults to 1cm for each margin). One to four values can be provided, similar to CSS.




`name` (text)
:   Name of the master page. It is for informational purpose and as a selection for [ClearPage](../clearpage).




`test` ([XPath expressions](/manual/data-processing/xpath))
:   If this xpath expression evaluates to true, this page is taken as a master page.




## Remarks
The contents of the element at [AtPageCreation](../atpagecreation) is executed, as soon as something will be placed on the page. The commands inside [AtPageShipout](../atpageshipout) are executed when switching to a new page.

When creating a new page, all page types are tried in reversed order. That means that the later defined master pages have a higher priority. This is important if more than one test in a Masterpage definition evaluates to true.


## Example

```xml
<DefineMasterpage name="right page" test=" sd:odd( sd:current-page() ) "/>
```
```xml
<DefineMasterpage name="left page" test=" sd:even( sd:current-page() ) "/>
```
```xml
<DefineMasterpage name="main part right" test=" sd:odd( sd:current-page() ) and $chapter='main' "/>
```
```xml
<DefineMasterpage name="right page" test="sd:odd( sd:current-page() )">
  <Margin left="1cm" right="1cm" top="1cm" bottom="1cm"/>
  <PositioningArea name="frame1">
    <PositioningFrame width="12" height="30" column="2" row="2"/>
    <PositioningFrame width="12" height="30" column="16" row="2"/>
  </PositioningArea>
  <AtPageCreation>
    <PlaceObject column="1">
      <!-- header -->
    </PlaceObject>
  </AtPageCreation>
  <AtPageShipout>
    <PlaceObject column="1">
      <!-- footer -->
    </PlaceObject>
  </AtPageShipout>
</DefineMasterpage>
```





