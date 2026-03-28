---
weight: 10
type: docs
linktitle: Data Files
---

# Structuring Your Data

> You can use any data format so long as it is XML.

The data file can be structured however you like -- there's no required schema. The only rule is that it must be well-formed XML. Other formats (CSV, Excel) can be converted to XML using the [Lua filter](../../advanced/lua-filter).

## Recommendations

If you can influence the data structure, here are some practical tips:

1. **Put data where you need it.** If article details are printed on the article page, the data should be structured by article, not by property type.

2. **Make variants recognizable.** If a new article group should start on a new page, make the group boundary visible in the data structure.

3. **Be explicit.** If an article number like `123-12345` encodes the group in the first three digits, break it out as an attribute instead of parsing it at runtime.

4. **Redundancy is fine.** Saving the full article number on each item is better than assembling it from group + suffix every time.

### Example structure

```xml
<productdata>
  <globalsettings>
    ...
  </globalsettings>
  <articlegroup name="interior lights" number="123">
    <article number="123-12345">
      <property1>...</property1>
      <property2>...</property2>
    </article>
    <article number="123-12346">
      <property1>...</property1>
      <property2>...</property2>
    </article>
  </articlegroup>
  <articlegroup name="exterior lights" number="124">
    <article number="124-23456">
      <property1>...</property1>
    </article>
  </articlegroup>
</productdata>
```

## Accessing data from the layout

XTS uses `<Record>` to match elements from the data file. When XTS encounters an element, it looks for a `<Record>` with a matching `element` attribute and executes the commands inside.

Within a record, you use XPath expressions to access attributes and child elements:

- `@nr` -- the `nr` attribute of the current element
- `description` -- the `description` child element
- `image/@mainimage` -- the `mainimage` attribute of the `image` child

See [How It Works](../../core-concepts/how-it-works) for the full processing model and [XPath](../xpath) for the expression language.
