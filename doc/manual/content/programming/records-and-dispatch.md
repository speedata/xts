---
weight: 80
type: docs
linktitle: Records and dispatch
---

# Records and dispatch

A `<Record>` connects a kind of data element to the layout code that should
handle it. Dispatching data to records is how an XTS layout walks your XML tree:
the root element triggers a record, that record processes its children, and so
on. This is the data-driven heart of XTS.

## Records match data elements

Each `<Record>` has a `match` attribute naming the element it handles:

```xml
<Record match="catalog">
  <ProcessNode select="*"/>
</Record>

<Record match="article">
  <PlaceObject>
    <TextBlock>
      <Paragraph><Value select="@name"/></Paragraph>
    </TextBlock>
  </PlaceObject>
</Record>
```

When XTS processes an element, it finds the best-matching record and runs its
body. Inside the body, the **context** is that element, so XPath expressions like
`@name` or `description` are relative to it.

Records are collected in the setup pass, so their order in the file does not
affect *which* record exists -- only the [priority rules](#priority) below decide
which one wins when several could match.

## ProcessNode -- dispatch children to their records

`<ProcessNode>` sends each selected element to its matching record. Use it when
different child elements need different treatment:

```xml
<Record match="catalog">
  <ProcessNode select="*"/>
</Record>

<Record match="article"><!-- … --></Record>
<Record match="category-header"><!-- … --></Record>
```

This is recursive: a record reached via `<ProcessNode>` can itself call
`<ProcessNode>` to descend further. That recursion is what walks the whole tree.

## ProcessNode versus ForAll

There are two ways to iterate over child elements:

- **`<ProcessNode>`** dispatches each child to its matching `<Record>` -- different
  elements can be handled by different records.
- **`<ForAll>`** runs the *same* commands for every matched element, right where
  it is written -- ideal for tables and lists. See [Control flow](/programming/control-flow).

```xml
<!-- ProcessNode: heterogeneous data, one record per element type -->
<Record match="catalog">
  <ProcessNode select="*"/>
</Record>

<!-- ForAll: homogeneous rows, inline -->
<Record match="catalog">
  <PlaceObject>
    <Table stretch="max">
      <ForAll select="article">
        <Tr>
          <Td><Paragraph><Value select="@name"/></Paragraph></Td>
          <Td><Paragraph><Value select="@price"/></Paragraph></Td>
        </Tr>
      </ForAll>
    </Table>
  </PlaceObject>
</Record>
```

## Matching with predicates

The `match` attribute supports XPath predicates, so you can route the same
element to different records based on its attributes or content:

```xml
<Record match="item[@type='invoice']"><!-- … --></Record>
<Record match="item[@type='credit']"><!-- … --></Record>
<Record match="item"><!-- … --></Record>   <!-- fallback for any other item -->
```

## Priority {#priority}

When more than one record could match an element, XTS picks the winner by these
rules:

1. A record **with a predicate** beats a record **without** one (more specific
   wins).
2. Among records of equal specificity, the **last one defined** wins.

These rules apply to `<ProcessNode>` dispatch. The record for the **data root
element** (and for the root of a `<LoadXML>` document) is looked up by plain
element name only: predicates are ignored there, so the root record must be
defined without a predicate and without a mode.

So in the example above, an `<item type="invoice">` is handled by the first
record; an `<item type="other">` falls through to the plain `item` record.

## Modes

A record can declare a `mode`, and `<ProcessNode>` can request one, so the same
element can be processed differently in different phases (for instance, once to
build a table of contents and once to render the body):

```xml
<Record match="chapter" mode="toc"><!-- … --></Record>
<Record match="chapter"><!-- … --></Record>   <!-- default mode -->

<ProcessNode select="chapter" mode="toc"/>
<ProcessNode select="chapter"/>
```

Only records whose `mode` matches the `<ProcessNode>` request are considered.

## See also

- [Execution model](/programming/execution-model) -- the two passes and where records fit.
- [Control flow](/programming/control-flow) -- `ForAll`, `Loop`, `Switch`, and friends.
- [Data files](/programming/data-files) -- structuring the XML that records consume.
