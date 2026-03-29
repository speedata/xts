---
weight: 20
type: docs
linktitle: Loops and Conditions
---

# Loops and Conditions

XTS has two programming levels: the **layout language** (XML commands like `<Loop>` and `<Switch>`) and **XPath** (expressions inside `select` and `test` attributes). This page covers the layout language constructs.

## ForAll -- iterating over data

`<ForAll>` runs the same commands for every matching element:

```xml
<ForAll select="article">
    <Tr>
        <Td><Paragraph><Value select="@name"/></Paragraph></Td>
        <Td><Paragraph><Value select="@price"/></Paragraph></Td>
    </Tr>
</ForAll>
```

Inside `<ForAll>`, the context switches to each matched element, so `@name` refers to the current article's name attribute.

## ProcessNode -- dispatching to records

`<ProcessNode>` sends each child element to its matching `<Record>`:

```xml
<Record match="catalog">
    <ProcessNode select="*"/>
</Record>

<Record match="article">
    <!-- handle each article -->
</Record>
```

Use `<ProcessNode>` when different elements need different treatment. Use `<ForAll>` when they all get the same treatment.

## Loop -- counting

`<Loop>` runs a fixed number of times:

```xml
<Loop select="10">
    <!-- runs 10 times -->
    <!-- $loopcounter is 1, 2, ..., 10 -->
</Loop>
```

The loop counter is stored in `$_loopcounter` by default, or you can name it: `<Loop select="10" variable="i">`.

## While and Until -- conditional loops

```xml
<SetVariable variable="i" select="1"/>
<While test="$i &lt;= 4">
    <PlaceObject>
        <TextBlock>
            <Paragraph><Value select="$i"/></Paragraph>
        </TextBlock>
    </PlaceObject>
    <SetVariable variable="i" select="$i + 1"/>
</While>
```

This outputs 1, 2, 3, 4. Don't forget to increment the variable, or you'll get an infinite loop.

{{< callout type="info" >}}
In XML, `<` must be written as `&lt;` inside attribute values. So `$i <= 4` becomes `$i &lt;= 4`.
{{< /callout >}}

`<Until>` is the inverse -- it runs *until* the condition becomes true:

```xml
<SetVariable variable="i" select="1"/>
<Until test="$i &gt; 4">
    <!-- runs while i <= 4 -->
    <SetVariable variable="i" select="$i + 1"/>
</Until>
```

## Switch/Case -- conditional logic

This works like `switch/case` in most programming languages:

```xml
<Switch>
    <Case test="$i = 1">
        <!-- when i is 1 -->
    </Case>
    <Case test="$i = 2">
        <!-- when i is 2 -->
    </Case>
    <Otherwise>
        <!-- when nothing else matches -->
    </Otherwise>
</Switch>
```

The `test` attribute expects an XPath expression that evaluates to `true()` or `false()`. Only the first matching `<Case>` is executed. `<Otherwise>` is optional and runs when no case matches.

## Practical example

A data-driven table with conditional styling:

```xml
<Record match="inventory">
    <PlaceObject>
        <Table stretch="max">
            <TableHead>
                <Tr>
                    <Td><Paragraph><Value>Product</Value></Paragraph></Td>
                    <Td><Paragraph><Value>Stock</Value></Paragraph></Td>
                </Tr>
            </TableHead>
            <ForAll select="item">
                <Tr>
                    <Td><Paragraph><Value select="@name"/></Paragraph></Td>
                    <Td>
                        <Switch>
                            <Case test="@qty &lt; 10">
                                <Paragraph style="color: red; font-weight: bold;">
                                    <Value select="@qty"/>
                                </Paragraph>
                            </Case>
                            <Otherwise>
                                <Paragraph><Value select="@qty"/></Paragraph>
                            </Otherwise>
                        </Switch>
                    </Td>
                </Tr>
            </ForAll>
        </Table>
    </PlaceObject>
</Record>
```
