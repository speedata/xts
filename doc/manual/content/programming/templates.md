---
weight: 60
type: docs
linktitle: Templates
---

# Templates

A **template** is a named, reusable block of layout code that you invoke where
you need it. Unlike a [function](/programming/functions), a template may contain layout
**actions** (such as `<PlaceObject>` or `<ClearPage>`): it runs in the normal
imperative flow, exactly once, at the call site, in document order. Templates are
the right tool for reusing *behaviour*.

## Defining and calling a template

Define a template at the top level with [`<Template>`](/reference/commands/template)
and invoke it with [`<CallTemplate>`](/reference/commands/calltemplate):

```xml
<Template name="headline">
  <Param name="text"/>
  <PlaceObject>
    <TextBlock>
      <Paragraph><B><Value select="$text"/></B></Paragraph>
    </TextBlock>
  </PlaceObject>
</Template>

<Record match="chapter">
  <CallTemplate name="headline">
    <Param name="text" select="@title"/>
  </CallTemplate>
</Record>
```

Like other definitions, `<Template>` is collected in the setup pass, so it can be
defined anywhere in the layout file -- before or after the records that call it.

## Parameters

Parameters are passed with [`<Param>`](/reference/commands/param) children:

- On a `<CallTemplate>`, a `<Param>`'s `select` expression is evaluated in the
  **caller's** context and passed as the argument.
- On a `<Template>` definition, a `<Param>` declares a parameter and may give a
  **default** via `select`, used when the caller omits that parameter.

```xml
<Template name="box">
  <Param name="title"/>
  <Param name="color" select="'black'"/>   <!-- default -->
  <PlaceObject>
    <TextBlock>
      <Paragraph style="color: { $color }">
        <Value select="$title"/>
      </Paragraph>
    </TextBlock>
  </PlaceObject>
</Template>

<CallTemplate name="box">
  <Param name="title" select="@name"/>
  <!-- color omitted, defaults to 'black' -->
</CallTemplate>
```

### Scoped bindings

Parameter bindings are **scoped to the template body**. Before the body runs, XTS
saves the current value of each parameter variable, binds the parameters, runs
the body, and then restores the previous values. A parameter named like an
existing variable shadows it only for the duration of the call:

```xml
<SetVariable variable="text" select="'outer'"/>

<CallTemplate name="headline">
  <Param name="text" select="'inner'"/>   <!-- $text is 'inner' inside the body -->
</CallTemplate>

<!-- here $text is 'outer' again -->
```

## Templates versus functions

Templates and functions look similar but serve opposite purposes:

| | [`<Function>`](/programming/functions) | [`<Template>`](/programming/templates) |
|---|---|---|
| Produces | a data **value** | a side **effect** (output) |
| May contain actions | no (action-free) | yes |
| Called from | an XPath expression (`fn:foo(…)`) | the layout flow (`<CallTemplate>`) |
| Runs | lazily, when the XPath engine evaluates it | once, at the call site, in order |
| Context | frozen (caller's arguments) | the live call-site context |

Rule of thumb: if you want a **value** to compute and reuse, write a function. If
you want to **do something** -- place a box, start a page, draw a header -- write a
template.

## When to reach for what

- Reusing fixed content? Store it in a [data variable](/programming/values-and-types) and
  splice it with `<CopyOf>`.
- Reusing a computed value? Write a [function](/programming/functions).
- Reusing behaviour with output? Write a template.

See [Data vs. action](/programming/data-and-actions) for the principle behind this split.
