---
weight: 10
linktitle: How It Works
type: docs
---

# How XTS Works

At its core, XTS does one thing: it reads a data file and a layout file, then produces a PDF. But the way it does this is worth understanding, because it affects how you structure your projects.

## The two files

Every XTS project has at least two files:

- **data.xml** -- your content (product data, article text, numbers, whatever)
- **layout.xml** -- your design rules (where things go, how they look)

The data file can be structured however you like, as long as it's well-formed XML. There's no required schema, no mandatory structure. The layout file uses XTS commands in the namespace `urn:speedata.de/2021/xts/en`.

## How processing works

When you run `xts`, this is what happens:

1. XTS reads the layout file. Top-level commands that aren't `<Record>` are executed immediately -- things like `<DefineColor>`, `<StyleSheet>`, `<SetGrid>`, or `<PageFormat>`.
2. All `<Record>` commands are stored for later.
3. XTS reads the data file and looks at the root element.
4. If there's a `<Record>` whose `element` attribute matches the root element's tag name, XTS executes that record.
5. Inside the record, you typically either output content directly or use `<ProcessNode>` / `<ForAll>` to iterate over child elements.

Here's the minimal setup:

```xml title="layout.xml"
<Layout xmlns="urn:speedata.de/2021/xts/en"
  xmlns:sd="urn:speedata.de/2021/xtsfunctions/en">

  <Record match="catalog">
    <ProcessNode select="*"/>
  </Record>

  <Record match="article">
    <PlaceObject>
      <TextBlock>
        <Paragraph>
          <Value select="@name"/>
        </Paragraph>
      </TextBlock>
    </PlaceObject>
  </Record>

</Layout>
```

```xml title="data.xml"
<catalog>
  <article name="Widget A"/>
  <article name="Widget B"/>
</catalog>
```

When XTS processes this:

1. It sees the root element `<catalog>` and runs the first `<Record>`.
2. `<ProcessNode select="*"/>` tells XTS to visit each child element.
3. For each `<article>`, XTS finds the matching `<Record match="article">` and runs it.
4. Each article's name gets placed on the page.

## ProcessNode vs. ForAll

There are two ways to iterate over child elements:

**`<ProcessNode>`** dispatches each child to its matching `<Record>`. This is great when different child elements need different treatment:

```xml
<Record match="catalog">
  <ProcessNode select="*"/>
</Record>

<Record match="article">
  <!-- handle articles -->
</Record>

<Record match="category-header">
  <!-- handle headers differently -->
</Record>
```

**`<ForAll>`** runs the same commands for every matching element, right where it's used. This is great for tables and lists:

```xml
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

## Accessing data with XPath

Inside a `<Record>`, you can use XPath expressions to access the current element's attributes and children:

- `@nr` -- the value of the `nr` attribute
- `description` -- the child element called `description`
- `image/@mainimage` -- the `mainimage` attribute of the `image` child
- `concat(@price, ' EUR')` -- string concatenation

XPath expressions appear in `select` attributes and inside curly braces `{...}` in some contexts.

## Real-time processing

One thing that makes XTS special: the layout is evaluated *as the PDF is being built*. This means you can ask questions like "Is there room on this page?" or "What's the current page number?" and react accordingly. You're not pre-calculating a layout and then rendering it -- you're building the PDF in a single pass (or a few passes if you need cross-references).

## The namespace boilerplate

Every layout file starts with this:

```xml
<Layout xmlns="urn:speedata.de/2021/xts/en"
    xmlns:sd="urn:speedata.de/2021/xtsfunctions/en">
```

The first namespace is for XTS commands (`PlaceObject`, `Record`, etc.). The second one, bound to `sd:`, gives you access to XTS-specific XPath functions like `sd:current-page()` or `sd:number-of-columns()`.

## What's next?

Now that you understand the pipeline, let's look at [the grid](../grid) -- the invisible framework that keeps your layouts aligned.
