# LoadXML



Load an XML file previously written by [SaveXML](../savexml) (attribute name) or a well formed XML file (attribute href). The regular data processing is interrupted and the contents of the data file is taken as a data source. If the file does not exist, the call to [LoadXML](../loadxml) is ignored.



##  Child elements

(none)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Function](../function), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Until](../until), [While](../while)


## Attributes



`href` (text, optional)
:   Filename of the XML file to load. Example: `myfile.xml`.




`name` (text, optional)
:   Name of the data file. Example: toc




## Example

```xml
<Record element="articles">
  <LoadXML name="toc"/>
  <ClearPage/>
  <ProcessNode select="article"/>
</Record>

```





