# Element



Create a data structure that can be used to save on the hard-drive between consecutive runs (with [SaveXML](../savexml)).



##  Child elements

[Attribute](../attribute), [Element](../element)

##  Parent elements

[Element](../element), [Message](../message), [SaveXML](../savexml), [SetVariable](../setvariable)


## Attributes



`name` (text)
:   Name of the element that gets created.




## Example

```xml
<SetVariable variable="articles">
  <Element name="articlelist">
    <Attribute name="name" select=" @name "/>
    <Attribute name="page" select="sd:current-page()"/>
  </Element>
</SetVariable>

```





