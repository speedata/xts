# SaveXML



Saves an element/attribute structure to be used in the next publisher run. The contents must have a tree structure.



##  Child elements

[Element](../element), [ForAll](../forall), [Loop](../loop), [Value](../value)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Function](../function), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Until](../until), [While](../while)


## Attributes



`elementname` (text)
:   Name of the root element that surrounds the elements given by the child elements.




`href` (text, optional)
:   Name of the file. Example: toc.xml




`name` (text, optional)
:   Name of the file. Example: toc




`select` ([XPath expressions](/manual/data-processing/xpath), optional)
:   Alternative to giving the data structure in the child elements.




## Example

```xml
  <Record match="data">
    <SetVariable variable="attributesvar">
      <Attribute name="att1" select="'Hello'" />
      <Attribute name="att2" select="123" />
    </SetVariable>

    <SaveXML name="toc" elementname="root" attributes="$attributesvar">
      <Element name="child">
        <Attribute name="attchild" select="999"/>
      </Element>
    </SaveXML>
  </Record>

```

This code saves an XML file to the disc which has this structure:


```xml
<root att1="Hello" att2="123">
 <child attchild="999"/>
</root>

```

```xml
<SaveXML name="toc" elementname="Contents">
  <Value select="$contents"/>
</SaveXML>

```

is equivalent to


```xml
<SaveXML name="toc" elementname="Contents" select="$contents"/>
```





