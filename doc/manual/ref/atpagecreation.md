# AtPageCreation



The contents of the element [AtPageCreation](../atpagecreation) is executed the first time the page is accessed. This is used in [DefineMasterpage](../definemasterpage).



##  Child elements

[ClearPage](../clearpage), [ForAll](../forall), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Slate](../slate), [Switch](../switch), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[DefineMasterpage](../definemasterpage)


## Attributes
(none)

## Example

```xml
<AtPageCreation>
  <PlaceObject column="1" row="1">
    <Textblock>
      <Paragraph>
        <Value select="$pageheader"/>
      </Paragraph>
    </Textblock>
  </PlaceObject>
</AtPageCreation>

```





