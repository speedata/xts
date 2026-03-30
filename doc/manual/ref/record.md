# Record



Contains the instructions to be executed when a data element matches the given match expression. The Record matching the root element will be called automatically, all further data handling must be done by the user via [ProcessNode](../processnode). Match expressions support XPath predicates for conditional matching (similar to XSLT template matching).



##  Child elements

[AttachFile](../attachfile), [ClearPage](../clearpage), [ForAll](../forall), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Slate](../slate), [Switch](../switch), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`match` ([XPath expressions](/manual/data-processing/xpath))
:   An XPath match expression. This can be a simple element name (e.g. `data`) or an element name followed by an XPath predicate in brackets (e.g. `item[@type='invoice']` or `item[not(@hidden='true')]`). When multiple Records match the same element, Records with predicates take priority over those without. Among Records with predicates, the last defined one wins.




`mode` (text, optional)
:   Name of the mode that matches the mode in [ProcessNode](../processnode).




## Example

```xml
<Record match="url" mode="output">
  <PlaceObject>
    <TextBlock>
      <Paragraph>
        <A href="https://www.speedata.de"><Value>website of speedata</Value></A>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>

```
```xml
<!-- Record with predicate: only matches item elements where type='invoice' -->
<Record match="item[@type='invoice']">
  <PlaceObject>
    <TextBlock>
      <Paragraph><Value>Invoice item</Value></Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>

<!-- Fallback Record: matches all other item elements -->
<Record match="item">
  <PlaceObject>
    <TextBlock>
      <Paragraph><Value>Other item</Value></Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>

```





