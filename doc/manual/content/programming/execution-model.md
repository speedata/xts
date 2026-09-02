---
weight: 10
type: docs
linktitle: Execution model
---

# The execution model

The single most important thing to understand about XTS is this: **a layout file
is a program that runs**. Commands execute one after another, in document order,
and they have effects. The cursor moves down the page, variables change, pages
are created. XTS is not a template engine that fills in blanks, and it is not a
declarative "build a tree, then render it once" system like XSL-FO.

## Why XTS is imperative

There is one central use case that shapes the whole design:

> Produce some output, look at where you are on the page, and decide what to do
> next.

"If there is enough room, do this, otherwise do that." Answering that question
requires **live layout state while the program runs** -- the current page, the
remaining space, the position of the cursor. A system that builds a complete
tree and renders it in one shot cannot do this, because the decision depends on
results that only exist *during* rendering.

So XTS has an imperative core, in the tradition of TeX and the speedata
Publisher. Commands run in order, they change state, and that state is
observable. This is the reason XTS exists, and it is why the language looks the
way it does.

## What happens when you run XTS

When you run `xts`, processing happens in two passes over the layout:

1. **Setup pass.** XTS reads the layout file and runs the top-level commands.
   Definitions register themselves: `<Record>`, `<Function>`, `<Template>`,
   `<DefineColor>`, `<DefineMasterPage>`, `<PageFormat>`, `<SetGrid>`,
   `<Options>`, and so on. Nothing is placed on a page yet.
2. **Processing pass.** XTS reads the data file, looks at its root element, finds
   the matching `<Record>`, and runs it. From there your records call each other
   via [`<ProcessNode>`](/programming/records-and-dispatch) and `<ForAll>`, and output
   commands such as `<PlaceObject>` actually emit content.

```xml title="layout.xml"
<Layout xmlns="urn:speedata.de/2021/xts/en"
  xmlns:sd="urn:speedata.de/2021/xtsfunctions/en">

  <!-- setup pass: this record is registered, not run yet -->
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

</Layout>
```

```xml title="data.xml"
<catalog>
  <article name="Widget A"/>
  <article name="Widget B"/>
</catalog>
```

Because definitions are collected in the setup pass, the *order* in which you
write `<Record>`, `<Function>`, and `<Template>` blocks does not matter -- they
are all known before the processing pass begins. The order of commands *inside* a
record, on the other hand, matters very much: that is the imperative part.

## Live state and introspection

Because the layout runs as the PDF is built, you can ask questions about the
current state and react to the answers. The `sd:` XPath functions expose this
live state -- for example `sd:current-page()`, `sd:current-row()`, or
`sd:number-of-columns()`. See [XPath in XTS](/programming/xpath) for the full list.

This is what makes "observe, then decide" possible:

```xml
<PlaceObject><TextBlock><!-- … --></TextBlock></PlaceObject>
<Message select="concat('now on page ', sd:current-page())"/>
```

Forward references -- "the total number of pages", "the page where chapter 3
starts" -- are a different problem. They are not available as live state, because
they depend on output that has not happened yet. XTS handles them with multiple
passes and an auxiliary file (the same mechanism used for a table of contents);
only the *current* state is a live query.

## Two languages in one

A layout file actually mixes two languages:

- The **layout language**: the XML commands (`<PlaceObject>`, `<ForAll>`,
  `<SetVariable>`, `<Switch>` …). This is the imperative layer.
- **XPath**: the expressions inside `select` and `test` attributes, and inside
  `{ … }` braces. This is a pure, side-effect-free expression language.

Keeping these two straight is the key to the next chapter: some commands build
*values* that XPath can query, while others perform *actions* that change the
document. That distinction is the subject of [Data vs. action](/programming/data-and-actions).
