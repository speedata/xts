# CallTemplate



Invoke a named [Template](../template). Parameters are passed with [Param](../param) children, whose `select` expressions are evaluated in the calling context. The parameter bindings are scoped to the template body.



##  Child elements

[Param](../param)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Template](../template), [Until](../until), [While](../while)


## Attributes



`name` (text)
:   The name of the template to call.




## Example

```xml
<Template name="headline">
  <Param name="text"/>
  <PlaceObject>
    <TextBlock>
      <Paragraph><B><Value select="$text"/></B></Paragraph>
    </TextBlock>
  </PlaceObject>
</Template>

<Record match="chapter">
  <CallTemplate name="headline">
    <Param name="text" select="@title"/>
  </CallTemplate>
</Record>

```





