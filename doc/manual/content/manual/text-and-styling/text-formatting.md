---
weight: 20
type: docs
linktitle: Text Formatting
---

# Text Formatting

XTS gives you several ways to format text: XTS commands (`<B>`, `<I>`, `<U>`), HTML markup inside `<HTML>`, and CSS classes/styles. Use whichever fits your situation.

## Bold, italic, underline

The most direct way to switch fonts is with the inline commands:

```xml
<PlaceObject>
  <Textblock>
    <Paragraph>
      <Value>A wonderful </Value>
      <B><Value>serenity</Value></B>
      <Value> has taken possession </Value>
      <I><Value>of my</Value>
        <Value> </Value>
        <B><Value>entire soul,</Value></B>
      </I>
      <Value> like these sweet mornings.</Value>
    </Paragraph>
  </Textblock>
</PlaceObject>
```

![text markup in layout](/manual/img/14-fonts.png)
<figcaption>Bold, italic, and nested bold-italic. Underline works with <code>&lt;U&gt;</code>.</figcaption>

These commands nest freely -- `<I><B>...</B></I>` gives you bold italic.

## HTML markup

If you prefer HTML-style formatting, use `<HTML>`:

```xml
<PlaceObject>
  <Textblock>
    <Paragraph>
      <HTML>A wonderful <b>serenity</b>
          has taken possession
          <i>of my <b>entire soul,</b></i>
          like these sweet mornings.
      </HTML>
    </Paragraph>
  </Textblock>
</PlaceObject>
```

The result is identical. You can also load HTML from your data file:

```xml
<HTML select="."/>
```

with data like:

```xml
<data>A wonderful <b>serenity</b> has taken possession
  <i>of my <b>entire soul,</b></i> like these sweet
  mornings.</data>
```

Tags can be uppercase (`<B>`) or lowercase (`<b>`).

{{< callout type="info" >}}
If your data contains raw HTML (not well-formed XML), use `sd:decode-html()` to interpret it:
`<HTML select="sd:decode-html(description)"/>`
{{< /callout >}}

## Paragraphs and text blocks

`<Textblock>` is a rectangular area that holds one or more `<Paragraph>` elements. Text blocks don't break across pages -- they're placed as a single unit. This makes them ideal for:

- Page numbers and headers
- Short descriptions and captions
- Column titles

Each paragraph can have its own class or inline style:

```xml
<Textblock>
  <Paragraph style="color: green">
    <Value>green text</Value>
  </Paragraph>
  <Paragraph>
    <Value>default text</Value>
  </Paragraph>
</Textblock>
```

## Spans

For inline styling within a paragraph, use `<Span>`:

```xml
<Paragraph>
    <Value>Regular text </Value>
    <Span class="highlight">
        <Value>highlighted</Value>
    </Span>
    <Value> and back to regular.</Value>
</Paragraph>
```

Spans support `class`, `style`, and `id` attributes, just like in HTML.

## Line breaks

Force a line break with `<Br/>`:

```xml
<Paragraph>
    <Value>First line</Value>
    <Br/>
    <Value>Second line</Value>
</Paragraph>
```

## CSS styling

All text elements support CSS styling via classes and inline styles. See [CSS and HTML](../css-html) for the full story.
