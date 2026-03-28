# Attribute



Create an attribute for the [Element](../element) data structure that can be saved to the hard drive with [SaveXML](../savexml).



##  Child elements

(none)

##  Parent elements

[Element](../element)


## Attributes



`name` (text)
:   Name of the attribute that is created.




`select` ([XPath expressions](/manual/data-processing/xpath))
:   The contents of the attribute




## Example

```xml
<Element name="Entry">
  <Attribute name="chapter" select="@name"/>
  <Attribute name="page" select="sd:current-page()"/>
</Element>

```

creates the following structure:


```xml
<Entry chapter="(contents of @name)" page="(the current page number)" />
```





