# DefineMasterPage



Define a master page. A master page is chosen depending on the criterion given with the attribute “test”.



##  Child elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [ForAll](../forall), [PositioningArea](../positioningarea)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`margin` (text, optional)
:   Set the margin of the page. One to four values can be provided, similar to CSS. If the attribute is omitted, the margins are taken from the CSS `@page` rule of the same name (see remarks), or 1cm without such a rule.




`name` (text)
:   Name of the master page. It is used as a selection for [ClearPage](../clearpage) and couples the master page to the CSS rule `@page name` in a [StyleSheet](../stylesheet).




`test` ([XPath expressions](/manual/data-processing/xpath))
:   If this xpath expression evaluates to true, this page is taken as a master page.




## Remarks
The contents of the element at [AtPageCreation](../atpagecreation) is executed, as soon as something will be placed on the page. The commands inside [AtPageShipout](../atpageshipout) are executed when switching to a new page.

When creating a new page, all page types are tried in reversed order. That means that the later defined master pages have a higher priority. This is important if more than one test in a Masterpage definition evaluates to true.

A CSS rule `@page name { ... }` in a [StyleSheet](../stylesheet) belongs to the master page with the same name, the generic `@page { ... }` rule is the base for all master pages. The page margin boxes of that rule (`@top-left`, `@bottom-center` and so on) are placed in the page margin when the page is written to the PDF, so `counter(page)` is the final page number. They never occupy grid cells. The `margin` declaration of the rule is used only when the `margin` attribute is missing; if both are given and differ, the attribute wins and a warning is issued. The CSS page selectors `:first`, `:left` and `:right` have no effect, use the attribute `test` instead.


## Example

```xml
<DefineMasterPage name="right page" margin="1cm" test=" sd:odd( sd:current-page() ) "/>
```
```xml
<DefineMasterPage name="left page" margin="1cm" test=" sd:even( sd:current-page() ) "/>
```
```xml
<DefineMasterPage name="main part right" margin="1cm" test=" sd:odd( sd:current-page() ) and $chapter='main' "/>
```
```xml
<DefineMasterPage name="right page" margin="1cm" test="sd:odd( sd:current-page() )">
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
</DefineMasterPage>
```





