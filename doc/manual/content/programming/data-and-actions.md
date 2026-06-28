---
weight: 20
type: docs
linktitle: Data vs. action
---

# Data vs. action

Every command in XTS plays one of two roles: it either **builds an inert value**
or it **performs a layout action**. This distinction is the backbone of the
language. Once you can tell which role a command plays, most of the rules about
variables, functions, and templates follow naturally.

## The two roles

**Constructors** build a value. The value is inert data: you can store it in a
variable, query it with XPath, count it, navigate into it, and splice it into
another command later. Constructors do not, by themselves, put anything on a
page.

```xml
<!-- builds a value: a sequence of Column elements -->
<Columns>
  <Column width="2cm"/>
  <Column width="3cm"/>
</Columns>
```

**Actions** have an effect on the document and produce no reusable value. They
*do* something: place an object, break to a new page, define a colour, attach a
file. You cannot store an action in a variable, because there is nothing to
store.

```xml
<!-- performs an action: emits content onto the page -->
<PlaceObject>
  <TextBlock><Paragraph><Value>Hello</Value></Paragraph></TextBlock>
</PlaceObject>
```

A third group, **control flow**, is transparent: `ForAll`, `Loop`, `While`,
`Until`, `Switch`, `Record`/`ProcessNode`. These commands do not have a role of
their own -- they inherit it from whatever they contain. A `<ForAll>` that
contains constructors builds values; a `<ForAll>` that contains actions performs
actions.

## The boundary is the absence of a type

Why does the distinction matter? Because it tells you what you are allowed to do
with a command's result.

A value can be named with a type (see [Values and types](/programming/values-and-types)):

```xml
<SetVariable variable="head" as="element(Column)*">
  <Columns>
    <Column width="2cm"/>
    <Column width="3cm"/>
  </Columns>
</SetVariable>
```

An action has **no type**, because it is not a value. That is not a gap in the
type system -- it *is* the boundary. Asking "what type does `<PlaceObject>`
produce?" has no sensible answer, so:

```xml
<!-- error: an action cannot be bound to a variable -->
<SetVariable variable="x" as="item()*">
  <PlaceObject>…</PlaceObject>
</SetVariable>
```

is rejected with a clear message:

```
action "PlaceObject" is not allowed in a value context
```

When you write `as="…"` on a `<SetVariable>`, you declare "this is data". XTS
then dispatches the body in a **value context**, and any action inside it -- at
any depth, even nested inside a constructor -- is an error. This is what keeps
`$head` safe to reuse: it is guaranteed to be inert data, never a hidden effect.

## Where the boundary is enforced

The same constructor/action classification drives three rules:

1. **`<SetVariable as="…">`** rejects actions in its body. A bound value is data,
   not behaviour. See [Values and types](/programming/values-and-types).
2. **`<Function>` bodies are action-free.** A function body is evaluated lazily by
   the XPath engine -- possibly more than once, in any order, or not at all. An
   action there would run an unpredictable number of times. Functions therefore
   build values only. See [Functions](/programming/functions).
3. **`<Template>` bodies may contain actions.** A template runs in the normal
   imperative flow, exactly once, at the call site -- so effects are welcome.
   This is the home for reusable *behaviour*. See [Templates](/programming/templates).

## Reuse: data versus behaviour

The data/action split gives you a clean way to choose a reuse mechanism:

| What you are reusing | Mechanism |
|---|---|
| Fixed content (a header that is always the same) | a **data variable** (`as="element()*"`) plus `<CopyOf>` |
| Parameterised *data* | a named [`<Function>`](/programming/functions), called from XPath |
| Parameterised *behaviour* with effects | a named [`<Template>`](/programming/templates) via `<CallTemplate>` |

## The bridge from data to a command

There is no special "evaluate this data now" command in XTS. Stored data simply
goes into a consuming command as an argument:

```xml
<SetVariable variable="head" as="element(Column)*">
  <Columns>
    <Column width="2cm"/>
    <Column width="3cm"/>
  </Columns>
</SetVariable>

<PlaceObject>
  <Table>
    <Columns><CopyOf select="$head"/></Columns>
    <Tr><Td>…</Td><Td>…</Td></Tr>
  </Table>
</PlaceObject>
```

`<Columns>` and `<Column>` are constructors, so `$head` holds queryable data that
you can inspect (`count($head)`, `$head[2]/@width`) and splice into the table
with `<CopyOf>`. The table consumes that data; the placement is the action.
