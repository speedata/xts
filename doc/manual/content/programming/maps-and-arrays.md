---
weight: 95
type: docs
linktitle: Maps and arrays
---

# Maps and arrays

XTS uses an XPath 3.1 engine, so besides plain sequences you have two structured
data types: **arrays** (ordered, indexable lists) and **maps** (key/value
dictionaries). They can be nested freely, stored in variables, passed to
[functions](/programming/functions), and built from your data. They are the idiomatic
replacement for the old "dynamic variable name" trick.

## Sequences, arrays, and maps

It helps to keep three things apart:

- A **sequence** is XPath's flat, always-unnested list of items. `(1, 2, 3)` is a
  sequence; sequences never contain other sequences.
- An **array** is a single item that *holds* an ordered list, and arrays can
  nest. `[1, 2, 3]` is one array containing three members.
- A **map** is a single item that holds key/value pairs. `map { 'a': 1 }` has one
  entry.

## Arrays

Build an array with a square-bracket constructor or the `array { … }`
constructor:

```xml
<SetVariable variable="nums" select="[10, 20, 30]"/>
<SetVariable variable="more" select="array { 1, 2, 3, 4 }"/>
```

Read a member by position with the `?` lookup operator (1-based), or get them all
with the `?*` wildcard:

```xml
<Value select="$nums?2"/>                              <!-- 20 -->
<Value select="string-join(for $n in $nums?* return string($n), ', ')"/>
```

Arrays nest, and lookups chain:

```xml
<SetVariable variable="people" select="[ map{'name':'Bob'}, map{'name':'Eve'} ]"/>
<Value select="$people?1?name"/>                       <!-- Bob -->
```

## Maps

Build a map with the `map { … }` constructor; keys and values are separated by a
colon, pairs by commas:

```xml
<SetVariable variable="prices" select="map { 'apple': 30, 'pear': 45 }"/>
```

Look up a value by key with `?`. If the key is not a simple name or a number (a
string with spaces, a computed value), put it in parentheses:

```xml
<Value select="$prices?apple"/>                        <!-- 30 -->
<Value select="$prices?('apple')"/>                    <!-- 30 -->
```

## The lookup operator `?`

The `?` operator works on both maps and arrays and is the everyday way to read
them:

| Expression | Meaning |
|---|---|
| `$map?key` | value for key `key` |
| `$map?('a key')` | value for a computed or non-name key |
| `$array?3` | the third member (1-based) |
| `$array?*` | all members as a sequence |
| `$a?1?name` | chained lookup into nested structures |

Constructors (`[ … ]`, `map { … }`, `array { … }`) and the `?` operator are part
of the XPath syntax, so they work **without any namespace declaration**.

## Function libraries

For anything beyond construction and lookup there are two standard function
libraries. Unlike the constructors, these functions are **prefixed**, so you must
declare their namespaces on the `<Layout>` root:

```xml
<Layout xmlns="urn:speedata.de/2021/xts/en"
    xmlns:sd="urn:speedata.de/2021/xtsfunctions/en"
    xmlns:map="http://www.w3.org/2005/xpath-functions/map"
    xmlns:array="http://www.w3.org/2005/xpath-functions/array">
```

Without the declaration you get `Could not find namespace for prefix "map"` (or
`"array"`).

Common map functions: `map:get`, `map:put`, `map:contains`, `map:keys`,
`map:size`, `map:remove`, `map:merge`, `map:entry`, `map:find`, `map:for-each`.

```xml
<Value select="map:size($prices)"/>                    <!-- 2 -->
<Value select="string-join(map:keys($prices), ', ')"/> <!-- apple, pear -->
<SetVariable variable="prices" select="map:put($prices, 'plum', 25)"/>
```

Common array functions: `array:get`, `array:size`, `array:head`, `array:tail`,
`array:append`, `array:subarray`, `array:insert-before`, `array:remove`,
`array:put`, `array:reverse`, `array:join`, `array:flatten`, `array:for-each`,
`array:filter`, `array:sort`, `array:fold-left`, `array:fold-right`.

```xml
<Value select="array:size($nums)"/>                    <!-- 3 -->
<Value select="array:get($nums, 1)"/>                  <!-- 10 -->
```

These are the standard XPath 3.1 `map:` and `array:` functions; see the
[goxpath documentation](https://github.com/speedata/goxpath) for the full
signatures.

{{< callout type="info" >}}
Maps and arrays are **immutable**. Functions like `map:put` or `array:append`
return a *new* value rather than modifying the original -- assign the result back
to a variable if you want to keep it.
{{< /callout >}}

## Building structures from data

Because arrays and maps are ordinary XPath values, you can build them from your
input and reuse them:

```xml
<!-- a lookup table from the data -->
<SetVariable variable="byid"
  select="map:merge(for $a in //article return map:entry($a/@id, $a/@name))"/>

<!-- later, resolve an id to a name -->
<Value select="$byid?('A-42')"/>
```

## See also

- [XPath in XTS](/programming/xpath) -- where XPath is used and the `sd:` layout functions.
- [Variables](/programming/variables) -- storing values, including maps and arrays.
- [Functions](/programming/functions) -- passing structured values in and out.
