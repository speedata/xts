---
weight: 85
type: docs
linktitle: XPath basics
---

# XPath basics

XPath is the expression language XTS uses to read your data: it appears in every
`select` and `test` attribute and inside `{ … }` braces. It is a W3C standard,
but it has its own model and a handful of rules that surprise newcomers. This
chapter is a practical primer -- enough to be productive. For *where* XPath is
used in XTS and the `sd:` layout functions, see [XPath in XTS](/programming/xpath);
for structured data, see [Maps and arrays](/programming/maps-and-arrays).

All examples below run against this data file:

```xml title="data.xml"
<catalog type="demo">
  <article id="a1" price="30" stock="5">Apple</article>
  <article id="a2" price="45" stock="0">Pear</article>
  <article id="a3" price="20" stock="12">Plum</article>
</catalog>
```

## The data model and the context

Your data file is a **tree** of elements, attributes, and text. Every XPath
expression is evaluated relative to a **context node** -- "where you are" in the
tree. Inside a `<Record match="catalog">`, the context is the `<catalog>`
element; inside a `<ForAll select="article">`, the context becomes each
`<article>` in turn.

A second idea runs through everything: **every value is a sequence**. A path that
selects three elements yields a sequence of three nodes; a single number is a
sequence of one. Sequences never nest (for nested data use
[arrays](/programming/maps-and-arrays)).

## Paths

A path navigates the tree from the context node. The two everyday steps are a
child element name and an attribute with `@`:

```xml
<Value select="count(article)"/>   <!-- 3   : child elements named article -->
<Value select="@type"/>            <!-- demo: the type attribute of the context -->
```

Other building blocks:

| Step | Meaning |
|---|---|
| `article` | child elements named `article` |
| `@price` | the `price` attribute |
| `.` | the context node itself |
| `..` | the parent |
| `*` | all child elements (any name) |
| `//article` | all `article` descendants, at any depth |
| `article/@id` | the `id` attribute of each child `article` |

A leading `/` makes a path absolute (from the document root) instead of relative
to the context.

## Predicates

A predicate in `[ … ]` filters a sequence. It can test a condition or pick by
position:

```xml
<Value select="article[@stock = '0']/@id"/>   <!-- a2  : the out-of-stock one -->
<Value select="article[1]/@id"/>              <!-- a1  : the first  (1-based!) -->
<Value select="article[last()]/@id"/>         <!-- a3  : the last -->
<Value select="string-join(article[position() &lt; 3]/@id, ',')"/>  <!-- a1,a2 -->
```

{{< callout type="warning" >}}
XPath positions are **1-based**: the first item is `[1]`, not `[0]`.
{{< /callout >}}

## Values and operators

XPath values are nodes, strings, numbers, and booleans. The usual operators work:

```xml
<Value select="2 + 3 * 4"/>    <!-- 14 -->
<Value select="7 div 2"/>      <!-- 3.5  (division is 'div', not '/') -->
<Value select="7 mod 2"/>      <!-- 1 -->
```

Comparisons and boolean logic:

```xml
<Value select="article[@price = 45]/@id"/>                          <!-- a2 -->
<Value select="count(article[@stock &gt; 0 and @price &lt; 40])"/>  <!-- 2 -->
```

Note the spelled-out words: `div`, `mod`, `and`, `or`. The `|` operator unions two
node sequences:

```xml
<Value select="count(article[1] | article[3])"/>   <!-- 2 -->
```

{{< callout type="info" >}}
`=` is a *general* comparison: when either side is a sequence it is true if **any**
pair matches. XPath 2.0+ also has *value* comparisons -- `eq`, `ne`, `lt`, `le`,
`gt`, `ge` -- which require exactly one item on each side and are stricter. Use
`eq` when you mean "these two single values are equal".
{{< /callout >}}

## Common functions

A small set of functions covers most needs:

```xml
<Value select="concat('id-', @id)"/>                  <!-- id-a1 -->
<Value select="count(article)"/>                      <!-- 3 -->
<Value select="string-join(article/@id, ', ')"/>      <!-- a1, a2, a3 -->
<Value select="substring('Hello', 1, 3)"/>            <!-- Hel  (1-based start) -->
<Value select="string-length('abc')"/>                <!-- 3 -->
<Value select="normalize-space('  a   b  ')"/>        <!-- a b -->
```

Also frequently used: `string()`, `number()`, `not()`, `position()`, `last()`,
`sum()`, `contains()`, `starts-with()`, `substring-before()` / `-after()`. The
context node taken as a string yields its text content:

```xml
<Value select="article[1]"/>          <!-- Apple : the element's text -->
<Value select="string(article[1])"/>  <!-- Apple : explicit, same result -->
```

For the complete list of supported functions see the
[XPath Functions Reference](/reference/xpath-functions) and the
[goxpath documentation](https://github.com/speedata/goxpath).

## Sequences, iteration, and conditionals

Build a sequence with commas, iterate with `for`, branch with `if`:

```xml
<Value select="count((1, 2, 3))"/>                                  <!-- 3 -->
<Value select="string-join(for $a in article return $a/@id, '-')"/> <!-- a1-a2-a3 -->
<Value select="if (@type = 'demo') then 'sample' else 'live'"/>     <!-- sample -->
```

These keep logic inside a single expression, which is often cleaner than the
layout-level [control-flow](/programming/control-flow) commands for small
computations.

## XTS-specific things to watch

A few rules trip up XPath newcomers in XTS specifically:

- **The context changes.** Inside `<ForAll>` or a matched `<Record>`, `.` is the current element, `@id` is its attribute, and `..` is its parent. Paths are always relative to *here*.
- **Escape `<` and `&`.** XPath lives inside XML attributes, so `<` must be written `&lt;` and `&` as `&amp;` (`>` is allowed as-is, but `&gt;` is fine too).
- **Braces switch into XPath.** In almost every attribute, `{ … }` evaluates an XPath expression: `width="{ $wd }"`. Only a few attributes that are XPath expressions themselves are exempt (`select` on `<SetVariable>`, `<ForAll>`, `<Loop>` and `<Attribute>`, `test` on `<Case>`, `<While>` and `<Until>`). See [Variables](/programming/variables).
- **Bodies are evaluated eagerly.** A `<SetVariable>` body is computed when the command runs, with the context frozen at that point. See [Values and types](/programming/values-and-types).

For example, inside a `<ForAll>` the context node is each element in turn, so `.`,
`@id`, and `..` all refer to the current `article`:

```xml
<ForAll select="article">
  <Message select="concat(@id, ' = ', ., ' (in ', ../@type, ')')"/>
  <!-- a1 = Apple (in demo), a2 = Pear (in demo), ... -->
</ForAll>
```

## See also

- [XPath in XTS](/programming/xpath) -- where XPath is used and the `sd:` functions.
- [Maps and arrays](/programming/maps-and-arrays) -- XPath 3.1 structured data.
- [Control flow](/programming/control-flow) -- the layout-level loops and conditionals.
