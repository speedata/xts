---
weight: 50
type: docs
linktitle: Functions
---

# Functions

A `<Function>` defines a reusable, parameterised **value**. You call it from an
XPath expression, and it returns data -- a number, a string, or an XML structure.
Functions are pure: they build values and have no side effects on the document.
For reusable *behaviour* (output, page breaks), use a [template](/programming/templates)
instead.

## Defining a function

```xml
<Layout xmlns="urn:speedata.de/2021/xts/en"
    xmlns:sd="urn:speedata.de/2021/xtsfunctions/en"
    xmlns:fn="mynamespace">

    <Record match="data">
        <Message select="fn:add(4, 3)"/>
    </Record>

    <Function name="fn:add">
        <Param name="first"/>
        <Param name="second"/>
        <Value select="$first + $second"/>
    </Function>
</Layout>
```

Output: `7`

Key points:

- The function name needs a **namespace prefix** (`fn:` here), and that namespace
  must be declared on the `<Function>` element or one of its ancestors (usually
  the root element).
- Parameters are accessed as local variables (`$first`, `$second`).
- Parameter variables have **local scope** -- unlike ordinary XTS variables, they
  do not leak out of the function body.
- Every `<Value>` in the body contributes to the return value.

## Functions return data, not effects

A function body is evaluated **lazily by the XPath engine**: when you write
`fn:add(4, 3)` in a `select`, the engine decides when -- and how often -- to
evaluate the body. Depending on the surrounding expression it might run more than
once, in a different order than written, or not at all (short-circuiting).

For that reason a function body must be **action-free**. A layout action such as
`<PlaceObject>` inside a function would run an unpredictable number of times,
giving non-deterministic output. XTS enforces this: an action in a function body
is rejected with

```
action "PlaceObject" is not allowed in a value context
```

If you need to *do* something repeatedly or conditionally with effects, that is
what [templates](/programming/templates) and the [control-flow](/programming/control-flow) commands are
for -- they run in the imperative flow, exactly when and as often as written.

See [Data vs. action](/programming/data-and-actions) for the principle behind this rule.

## Functions returning XML

A function can build and return an XML structure, which makes it a clean way to
produce [typed data](/programming/values-and-types) from parameters:

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

The curly braces `{.}` switch into XPath mode, where `.` is the current item in
the `<ForAll>` iteration.

## Function or template?

| You want | Use |
|---|---|
| a value computed from arguments | a **function**, called from XPath |
| output / effects, reused with parameters | a [template](/programming/templates) |
| fixed data, reused as-is | a [data variable](/programming/values-and-types) + `<CopyOf>` |
