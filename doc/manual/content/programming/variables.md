---
weight: 40
type: docs
linktitle: Variables
aliases:
  - /manual/data-processing/variables-and-functions/
---

# Variables

All variables in XTS are **globally visible** by default. This is intentional -- since XTS executes layout code as it builds the PDF (including page hooks like `<AtPageShipout>`), variables must be accessible everywhere. The exception is parameters inside a [`<Function>`](/programming/functions) or [`<Template>`](/programming/templates), which are scoped to the body.

### Setting variables

```xml
<SetVariable variable="count" select="42"/>
<SetVariable variable="name" select="'hello'"/>
```

### Reading variables

Use `$variablename` in XPath expressions:

```xml
<Value select="$count"/>
<Message select="concat('Name is: ', $name)"/>
```

### Storing complex content

Variables can hold not just simple values but entire XML structures:

```xml
<SetVariable variable="greeting">
    <Paragraph>
        <Value>Hello, world!</Value>
    </Paragraph>
</SetVariable>

<PlaceObject>
    <TextBlock>
        <CopyOf select="$greeting"/>
    </TextBlock>
</PlaceObject>
```

A practical use case is storing table column definitions for reuse:

```xml
<SetVariable variable="tablecolumns">
    <Columns>
        <Column width="1cm"/>
        <Column width="4mm"/>
        <Column width="1cm"/>
    </Columns>
</SetVariable>

<Table>
    <CopyOf select="$tablecolumns"/>
    <Tr><!-- ... --></Tr>
</Table>
```

### Appending to variables

Build up content incrementally:

```xml
<SetVariable variable="foo">
    <Value>Hello</Value>
</SetVariable>

<SetVariable variable="foo">
    <Value select="$foo"/>
    <Value>, world!</Value>
</SetVariable>
<!-- $foo is now "Hello, world!" -->
```

This also works for building XML structures -- use `<CopyOf>` to keep the
collected nodes intact (`<Value select>` would flatten them to text):

```xml
<SetVariable variable="toc">
    <CopyOf select="$toc"/>
    <Element name="entry">
        <Attribute name="title" select="@name"/>
        <Attribute name="page" select="sd:current-page()"/>
    </Element>
</SetVariable>
```

### Evaluation time

Variable contents with child elements are evaluated **immediately** when `<SetVariable>` is executed. So this:

```xml
<SetVariable variable="greeting"><Value>nice</Value></SetVariable>
<SetVariable variable="tmp"><CopyOf select="$greeting"/></SetVariable>
<SetVariable variable="greeting"><Value>cruel</Value></SetVariable>
<!-- $tmp is still "nice" -->
```

This means variables must not contain output commands like `<PlaceObject>` -- those would execute immediately during assignment. Declaring the variable's type with `as` makes this a checked rule rather than a silent trap: see [Values and types](/programming/values-and-types) and [Data vs. action](/programming/data-and-actions).

### Collections

For lists and dictionaries, use XPath's built-in **arrays** and **maps** rather
than juggling many variables -- they are real values you can store, nest, query,
and pass around:

```xml
<SetVariable variable="nums" select="[10, 20, 30]"/>
<SetVariable variable="prices" select="map { 'apple': 30, 'pear': 45 }"/>

<Value select="$nums?2"/>        <!-- 20 -->
<Value select="$prices?apple"/>  <!-- 30 -->
```

See [Maps and arrays](/programming/maps-and-arrays) for the full story.

### Dynamic variable names

The variable name itself can be computed, which is occasionally useful when the
name depends on the data:

```xml
<SetVariable variable="{ concat('item', 1) }" select="'First'"/>
<SetVariable variable="{ concat('item', 2) }" select="'Second'"/>

<!-- Read back: sd:variable() joins its arguments into a name and returns the value -->
<Message select="sd:variable(('item', 1))"/>
<Message select="sd:variable(('item', 2))"/>
```

Reach for this only when you genuinely need dynamically named *global* variables.
For ordinary collections, a [map or array](/programming/maps-and-arrays) is clearer and keeps
the data in one value.

## See also

- [Maps and arrays](/programming/maps-and-arrays) -- arrays, maps, and the `?` lookup operator.
- [Values and types](/programming/values-and-types) -- typing a variable with `as` and the
  queryable data band.
- [Functions](/programming/functions) -- parameterised, reusable values.
- [Templates](/programming/templates) -- parameterised, reusable behaviour.
