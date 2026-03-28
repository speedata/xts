# Record



Contains the instructions to be executed when a data element matches the given match expression. The Record matching the root element will be called automatically, all further data handling must be done by the user via [ProcessNode](../processnode). Match expressions support XPath predicates for conditional matching (similar to XSLT template matching).



##  Child elements

[ClearPage](../clearpage), [ForAll](../forall), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Slate](../slate), [Switch](../switch), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`match` ([XPath expressions](/manual/data-processing/xpath))
:   An XPath match expression. This can be a simple element name (e.g. `data`) or an element name followed by an XPath predicate in brackets (e.g. `item[@type='invoice']` or `item[not(@hidden='true')]`). When multiple Records match the same element, Records with predicates take priority over those without. Among Records with predicates, the last defined one wins.




`mode` (text, optional)
:   Name of the mode that matches the mode in [ProcessNode](../processnode).




## Examples

Simple match by element name:

```xml
<Record match="url" mode="output">
  <PlaceObject>
    <Textblock>
      <Paragraph>
        <A href="https://www.speedata.de"><Value>website of speedata</Value></A>
      </Paragraph>
    </Textblock>
  </PlaceObject>
</Record>
```

Match with XPath predicate — the Record with a predicate takes priority over the fallback:

```xml
<!-- Only matches item elements where type='invoice' -->
<Record match="item[@type='invoice']">
  <PlaceObject>
    <Textblock>
      <Paragraph><Value>Invoice item</Value></Paragraph>
    </Textblock>
  </PlaceObject>
</Record>

<!-- Fallback: matches all other item elements -->
<Record match="item">
  <PlaceObject>
    <Textblock>
      <Paragraph><Value>Other item</Value></Paragraph>
    </Textblock>
  </PlaceObject>
</Record>
```





