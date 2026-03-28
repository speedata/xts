---
weight: 40
type: docs
linktitle: Lists
---

# Lists

XTS supports ordered and unordered lists, both through XTS elements and HTML markup.

## XTS list elements

```xml
<PlaceObject>
    <Textblock>
        <Ul>
            <Li><Value>First item</Value></Li>
            <Li><Value>Second item</Value></Li>
            <Li><Value>Third item</Value></Li>
        </Ul>
    </Textblock>
</PlaceObject>
```

For numbered lists, use `<Ol>`:

```xml
<Ol>
    <Li><Value>Step one</Value></Li>
    <Li><Value>Step two</Value></Li>
    <Li><Value>Step three</Value></Li>
</Ol>
```

## HTML lists

You can also create lists inside `<HTML>` blocks:

```xml
<HTML>
    <ul>
        <li>Apples</li>
        <li>Oranges</li>
        <li>Bananas</li>
    </ul>
    <ol>
        <li>First</li>
        <li>Second</li>
    </ol>
</HTML>
```

## Styling lists with CSS

```xml
<Stylesheet>
    ul { list-style-type: disc; }
    ol { list-style-type: decimal; }
    li { padding-left: 0; margin-bottom: 4pt; }
</Stylesheet>
```

Lists can appear inside table cells, text blocks, and anywhere else that accepts content.
