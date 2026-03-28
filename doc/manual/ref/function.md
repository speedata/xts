# Function



Define a function



##  Child elements

[ClearPage](../clearpage), [Column](../column), [Columns](../columns), [ForAll](../forall), [LoadXML](../loadxml), [Loop](../loop), [Message](../message), [NextFrame](../nextframe), [NextRow](../nextrow), [Param](../param), [PlaceObject](../placeobject), [ProcessNode](../processnode), [SaveXML](../savexml), [SetVariable](../setvariable), [Slate](../slate), [Switch](../switch), [Value](../value)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`name` (text, optional)
:   The name of the function (with namespace prefix).




## Example

```xml
<Layout xmlns="urn:speedata.de/2021/xts/en"
    xmlns:sd="urn:speedata.de/2021/xtsfunctions/en"
    xmlns:fn="mynamespace"
    >

  <Record element="data">
    <PlaceObject>
        <Textblock>
            <Paragraph>
                <Value select="fn:add(3,4)"></Value>
            </Paragraph>
        </Textblock>
    </PlaceObject>
</Record>

...

<Function name="fn:add">
    <Param name="a" />
    <Param name="b" />
    <Value select="$a + $b" />
</Function>

```

Print out the number 7.







