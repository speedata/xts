# ProcessNode



Executes all given nodes. The elements, that are to be executed, are given with the attribute `selection`.



##  Child elements

(none)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Function](../function), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Until](../until), [While](../while)


## Attributes



`limit` (number, optional)
:   Limits the number of items processed with this command




`mode` (text, optional)
:   Name of the mode. This must match the mode at the corresponding [Record](../record) element. With this it is possible to have different rules for the same element.




`select` ([XPath expressions](/manual/data-processing/xpath))
:   Selection of child elements, that are to be processed.




## Example

```xml
<ProcessNode select="*" mode="sum" />
```





