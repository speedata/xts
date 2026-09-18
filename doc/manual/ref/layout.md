# Layout



This command is the root element in the Layout instructions.



##  Child elements

[AttachFile](../attachfile), [DefineColor](../definecolor), [DefineMasterPage](../definemasterpage), [Function](../function), [Message](../message), [Options](../options), [PDFOptions](../pdfoptions), [PageFormat](../pageformat), [Record](../record), [Section](../section), [SetGrid](../setgrid), [SetVariable](../setvariable), [StyleSheet](../stylesheet), [Template](../template), [Trace](../trace)

##  Parent elements

(none)


## Attributes



`version` (number, optional)
:   The version of XTS the layout was written for, such as `0.1` or `0.1.2`. The parts are compared from the left with the version of the running XTS, so `0.1` is satisfied by 0.1.0, 0.1.5 and 0.2.0. When the layout asks for a newer version, XTS stops with an error instead of producing a document that silently lacks the features the layout relies on. Development builds of XTS accept every version. See the chapter Versions and Compatibility in the manual.




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





