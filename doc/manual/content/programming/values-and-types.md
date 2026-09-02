---
weight: 30
type: docs
linktitle: Values and types
---

# Values and types

A constructor builds a value, and you can give that value a **type** when you
bind it to a variable. The type is declared with the `as` attribute on
[`<SetVariable>`](/reference/commands/setvariable), using the sequence-type
syntax familiar from XSLT 2.0+ and XQuery.

## Declaring a type with `as`

```xml
<SetVariable variable="head" as="element(Column)*">
  <ForAll select="cols">
    <Column width="{@w}"/>
  </ForAll>
</SetVariable>
<!-- queryable: count($head), $head/@width, $head[3] -->
```

The `as` attribute does two things:

1. It declares that the variable holds a **data value** (not an action). The body
   is dispatched in a value context, so any layout action inside it is rejected.
   See [Data vs. action](/programming/data-and-actions).
2. It documents the shape of the value for the reader (and, in a later XTS
   version, will be checked against the produced value).

{{< callout type="info" >}}
Today, `as` enforces the data/action boundary: a body containing an action is
rejected. Full sequence-type checking of the produced value (verifying it really
is, say, `element(Column)*` with the right cardinality) is planned for a later
release. Writing `as` now is forward-compatible.
{{< /callout >}}

## Sequence-type syntax

The `as` value follows the XPath/XQuery sequence-type grammar:

- **Atomic types** are prefixed with `xs:`: `xs:integer`, `xs:string`,
  `xs:boolean`, `xs:double`.
- **Node kind tests** have no prefix: `element()`, `element(Column)`, `node()`,
  `attribute()`, `text()`.
- **Occurrence indicators** follow the type:
  - nothing -- exactly one
  - `?` -- zero or one
  - `*` -- zero or more
  - `+` -- one or more

So `element(Column)*` is "zero or more `Column` elements", `xs:string` is "exactly
one string", and `node()?` is "an optional node".

## The queryable data band

The point of typed values is that constructor output is **real, queryable XML
data** -- the same XPath data model (XDM) as your input data file. A variable
built from constructors can be navigated just like data:

```xml
<SetVariable variable="head" as="element(Column)*">
  <ForAll select="cols">
    <Column width="{@w}"/>
  </ForAll>
</SetVariable>

<Value select="count($head)"/>        <!-- how many columns -->
<Value select="$head[2]/@width"/>      <!-- width of the second one -->
```

This is why `<Column>` produces a `Column` element you can address, rather than an
opaque internal object. Data you build behaves exactly like data you load.

## Eager versus deferred evaluation

The type also determines *when* a value is computed:

- A **value type** (`element(...)*`, `xs:string`, …) is evaluated **immediately**,
  when the `<SetVariable>` runs. The context is frozen at that moment. In the
  example above, `{@w}` is read against the context that is current at the
  `<SetVariable>`, so `$head` is safe to reuse later without surprises -- there is
  no "context at the use site" footgun.
- A **function** (see [Functions](/programming/functions)) is **deferred**: its body runs only
  when the function is called, against the caller's arguments.

This is the practical reason a data variable must not contain output actions: it
is evaluated eagerly, so an action inside it would fire at assignment time, not
where you expect. The data/action rule turns that mistake into an explicit error.

## Reuse with `<CopyOf>`

Once you have a typed value, splice it into a consuming command with
[`<CopyOf>`](/reference/commands/copyof), which copies the nodes as-is
(analogous to `xsl:copy-of`):

```xml
<Table>
  <Columns><CopyOf select="$head"/></Columns>
  <!-- … -->
</Table>
```

See [Data vs. action](/programming/data-and-actions) for how this "data into a consuming
command" bridge replaces the need for any explicit evaluation step.
