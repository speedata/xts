---
weight: 30
type: docs
linktitle: Variables and Functions
---

# Variables and Functions

## Variables

All variables in XTS are **globally visible** by default. This is intentional -- since XTS executes layout code as it builds the PDF (including page hooks like `<AtPageShipout>`), variables must be accessible everywhere.

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
    <Textblock>
        <Value select="$greeting"/>
    </Textblock>
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
    <Value select="$tablecolumns"/>
    <Tr>...</Tr>
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

This also works for building XML structures:

```xml
<SetVariable variable="toc">
    <Value select="$toc"/>
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
<SetVariable variable="tmp"><Value select="$greeting"/></SetVariable>
<SetVariable variable="greeting"><Value>cruel</Value></SetVariable>
<!-- $tmp is still "nice" -->
```

This means variables must not contain output commands like `<PlaceObject>` -- those would execute immediately during assignment.

### Simulating arrays

XTS doesn't have built-in arrays, but you can simulate them with dynamic variable names:

```xml
<SetVariable variable="{ concat('item', 1) }" select="'First'"/>
<SetVariable variable="{ concat('item', 2) }" select="'Second'"/>

<!-- Read back -->
<Message select="sd:variable(('item', 1))"/>
<Message select="sd:variable(('item', 2))"/>
```

The `sd:variable()` function concatenates its arguments into a variable name and returns the value.

## Functions

You can define custom functions in the layout file:

```xml
<Layout xmlns="urn:speedata.de/2021/xts/en"
    xmlns:sd="urn:speedata.de/2021/xtsfunctions/en"
    xmlns:fn="mynamespace">

    <Record element="data">
        <Message select="fn:add(4, 3)"/>
    </Record>

    <Function name="fn:add">
        <Param name="first"/>
        <Param name="second"/>
        <Value select="$first + $second"/>
    </Function>
</Layout>
```

Output: `7.000000`

Key points:

- The function namespace (`fn:` in this example) must be declared on the root element.
- Parameters are accessed as local variables (`$first`, `$second`).
- Variables inside functions have **local scope** (unlike the rest of XTS).
- Every `<Value>` in the function body contributes to the return value.
- Functions can output objects too, not just return values.

### Functions returning XML

Functions can build and return XML structures:

```xml
<Function name="fn:cols">
    <Param name="colspec"/>
    <Columns>
        <ForAll select="$colspec">
            <Column width="{.}"/>
        </ForAll>
    </Columns>
</Function>

<!-- Call with a sequence -->
<Table>
    <Value select="fn:cols(('2cm', '3cm'))"/>
    <Tr>
        <Td><Paragraph><Value>two cm</Value></Paragraph></Td>
        <Td><Paragraph><Value>three cm</Value></Paragraph></Td>
    </Tr>
</Table>
```

The curly braces `{.}` switch into XPath mode, where `.` is the current item in the `<ForAll>` iteration.
