# Layout



This command is the root element in the Layout instructions.



##  Child elements

[AttachFile](../attachfile), [DefineColor](../definecolor), [DefineMasterPage](../definemasterpage), [Function](../function), [Message](../message), [Options](../options), [PDFOptions](../pdfoptions), [PageFormat](../pageformat), [Record](../record), [Section](../section), [SetGrid](../setgrid), [SetVariable](../setvariable), [StyleSheet](../stylesheet), [Trace](../trace)

##  Parent elements

(none)


## Attributes



`version` (number, optional)
:   Minimum publisher version required. If major or minor version differ, give a warning. Format: 1.6.12 (revision number can be left out).




## Example


This is a complete example for a layout rule set. The first part is the data file (save as `data.xml`) and the second the layout instructions (`layout.xml`).


```xml
<root>
  <elt greeting="Hello world!" />
</root>
```
```xml
<Layout xmlns="urn:speedata.de/2021/xts/en"
  xmlns:sd="urn:speedata.de/2021/xtsfunctions/en">

  <Options mainlanguage="English (USA)"/>

  <Record match="root">
    <ProcessNode select="elt"/>
  </Record>

  <Record match="elt">
    <PlaceObject>
      <TextBlock>
        <Paragraph>
          <Value select="@greeting"></Value>
        </Paragraph>
      </TextBlock>
    </PlaceObject>
  </Record>
</Layout>
```





