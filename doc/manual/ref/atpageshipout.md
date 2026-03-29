# AtPageShipout



The enclosed instructions will be executed when the page is placed into the PDF file. Used in [DefineMasterPage](../definemasterpage).



##  Child elements

[ClearPage](../clearpage), [ForAll](../forall), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Slate](../slate), [Switch](../switch), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[DefineMasterPage](../definemasterpage)


## Attributes
(none)

## Example

```xml
<AtPageShipout>
  <PlaceObject column="1" row="20">
    <TextBlock>
      <Paragraph>
        <Value select="sd:current-page()"/>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
</AtPageShipout>

```





