# AtPageCreation



The contents of the element [AtPageCreation](../atpagecreation) is executed the first time the page is accessed. This is used in [DefineMasterPage](../definemasterpage).



##  Child elements

[AttachFile](../attachfile), [ClearPage](../clearpage), [ForAll](../forall), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Slate](../slate), [Switch](../switch), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[DefineMasterPage](../definemasterpage)


## Attributes
(none)

## Example

```xml
<AtPageCreation>
  <PlaceObject column="1" row="1">
    <TextBlock>
      <Paragraph>
        <Value select="$pageheader"/>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
</AtPageCreation>

```





