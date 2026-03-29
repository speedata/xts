# SetVariable



Associates a value with a variable name. The value can be a simple value or a more complex one consisting of several elements.



##  Child elements

[CopyOf](../copyof), [Element](../element), [Paragraph](../paragraph), [Table](../table), [TextBlock](../textblock), [Until](../until), [Value](../value), [While](../while)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Function](../function), [Layout](../layout), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Section](../section), [Until](../until), [While](../while)


## Attributes



`select` ([XPath expressions](/manual/data-processing/xpath), optional)
:   The value of the contents.




`trace` (optional)
:   Show information about the assignment in the log file.



    `yes`
    :    Show information.



    `no`
    :    Don't show information (default).




`variable` (text)
:   The name of the variable that holds the contents.




## Remarks
Variables have global scope.


## Example

```xml
<Record match="product">
  <SetVariable variable="wd" select="5"/>
  <PlaceObject>
    <TextBlock width="{ $wd }">
      <Paragraph>
        <Value select="$articlenumber"/>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>

```

The following example shows a more complex scenario: you can collect complex elements in a variable.


```xml
<Record match="products">
  <SetVariable variable="articletext"/>
  <ProcessNode select="article"/>
  <PlaceObject>
    <TextBlock>
      <Value select=" $articletext "/>
    </TextBlock>
  </PlaceObject>
</Record>

<Record match="article">
  <SetVariable variable="articletext">
    <!-- the previous contents is added -->
    <Value select="$articletext"/>
    <Paragraph>
      <Value select=" @description "/>
    </Paragraph>
  </SetVariable>
</Record>

```





