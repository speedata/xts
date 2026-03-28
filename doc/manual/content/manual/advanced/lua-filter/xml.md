---
weight: 20
type: docs
linktitle: xml
---

# xml

The xml module creates and reads XML files.

```lua
xml = require("xml")
```

## Functions

### run_xslt(table) / run_xslt(stylesheet, source, out)

Runs an XSLT transformation using the built-in XSLT processor.

**Table form:**

```lua
xml = require("xml")

ok, msg = xml.run_xslt({
    stylesheet = "transform.xsl",
    source = "raw-data.xml",
    out = "data.xml",
    initialtemplate = "main",          -- optional
    params = { year = "2026", draft = "yes" }  -- optional
})
if not ok then
    print(msg)
    os.exit(-1)
end
```

**Positional form:**

```lua
ok, msg = xml.run_xslt("transform.xsl", "raw-data.xml", "data.xml")
```

| Parameter | Description |
|-----------|-------------|
| `stylesheet` | Path to the XSLT stylesheet (required) |
| `source` | Path to the source XML document |
| `out` | Path for the output file |
| `initialtemplate` | Named template to invoke instead of the default |
| `params` | Table of stylesheet parameters |

### encode_table(table, filename)

Creates an XML file from a Lua table. The filename defaults to `data.xml` if omitted.

Each table represents a node with the following string keys:

| Key | Description |
|-----|-------------|
| `type` | `"element"` or `"comment"` |
| `name` | Element name |
| `value` | Comment text (for `type = "comment"`) |
| `attribs` | Table of attributes (key-value pairs) |

Integer keys become child elements or text content.

```lua
xml = require("xml")

tbl = {
    type = "element",
    name = "catalog",
    attribs = { version = "2.0" },
    {
        type = "element",
        name = "product",
        attribs = { id = "42", category = "books" },
        "The Art of Typesetting",
    },
    {
        type = "element",
        name = "product",
        attribs = { id = "43", category = "tools" },
        "Layout Grid",
    },
    {
        type = "comment",
        value = " end of catalog ",
    },
}

ok, msg = xml.encode_table(tbl, "catalog.xml")
```

This produces:

```xml
<catalog version="2.0"><product id="42" category="books">The Art of Typesetting</product><product id="43" category="tools">Layout Grid</product><!-- end of catalog --></catalog>
```

Attribute names can contain any characters, including namespace prefixes:

```lua
{
    type = "element",
    name = "root",
    attribs = {
        ["xmlns:dc"] = "http://purl.org/dc/elements/1.1/",
        ["_special"] = "works too",
    },
}
```

### decode_xml(filename)

Reads an XML file and returns a Lua table with the same structure as described above.

```lua
xml = require("xml")

ok, root = xml.decode_xml("data.xml")
if not ok then
    print("Error: " .. root)
    os.exit(-1)
end

-- root.name contains the root element name
-- root.attribs contains the root element's attributes (if any)
-- root[1] is the first child element or text node
print(root.name)

-- iterate child elements
for i = 1, #root do
    local child = root[i]
    if type(child) == "table" then
        print(child.name, child.attribs and child.attribs.id)
    else
        print("Text: " .. child)
    end
end
```
