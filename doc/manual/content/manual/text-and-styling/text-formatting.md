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
  <TextBlock>
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
  </TextBlock>
</PlaceObject>
```

![text markup in layout](/manual/img/14-fonts.png)
<figcaption>Bold, italic, and nested bold-italic. Underline works with <code>&lt;U&gt;</code>.</figcaption>

These commands nest freely -- `<I><B>...</B></I>` gives you bold italic.

## HTML markup

If you prefer HTML-style formatting, use `<HTML>`:

```xml
<PlaceObject>
  <TextBlock>
    <Paragraph>
      <HTML>A wonderful <b>serenity</b>
          has taken possession
          <i>of my <b>entire soul,</b></i>
          like these sweet mornings.
      </HTML>
    </Paragraph>
  </TextBlock>
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

`<TextBlock>` is a rectangular area that holds one or more `<Paragraph>` elements. Text blocks don't break across pages -- they're placed as a single unit. For text that breaks between lines and continues in the next frame or on the next page, put the paragraphs into a [`<Flow>`](../../core-concepts/flow) instead. A text block is ideal for:

- Page numbers and headers
- Short descriptions and captions
- Column titles

Each paragraph can have its own class or inline style:

```xml
<TextBlock>
  <Paragraph style="color: green">
    <Value>green text</Value>
  </Paragraph>
  <Paragraph>
    <Value>default text</Value>
  </Paragraph>
</TextBlock>
```

## Spans

For inline styling within a paragraph, use `<Span>`:

```xml
<StyleSheet>
    .highlight { color: firebrick; font-weight: bold; }
</StyleSheet>

<Paragraph>
    <Value>Regular text </Value>
    <Span class="highlight">
        <Value>highlighted</Value>
    </Span>
    <Value> and back to regular.</Value>
</Paragraph>
```

Spans support `class`, `style`, and `id` attributes, just like in HTML. Properties that change the font or the text color work on a span, and so does `background-color`: the box behind the text is as high as the font and follows the text across line breaks. With `position: relative`, `top` moves a span down and `bottom` moves it up by a fixed amount, without moving the lines or the text around it.

## Line breaks

Force a line break with `<Br/>`:

```xml
<Paragraph>
    <Value>First line</Value>
    <Br/>
    <Value>Second line</Value>
</Paragraph>
```

## Alignment and spacing

`text-align` sets how the lines of a paragraph are aligned. The default is `start`, which is left for left-to-right text.

- With `left`, `right`, `center`, `start` or `end` the text is ragged: every line is set at its natural width, the word spaces keep their size, and a word that does not fit moves to the next line.
- With `justify` the word spaces stretch and shrink so that every line but the last fills the measure. Only justified text uses the font expansion (`-bag-font-expansion`), which widens or narrows the glyphs a little to even out the spaces.

```xml
<StyleSheet>
    .body { text-align: justify; hyphens: auto; }
    .caption { text-align: center; }
</StyleSheet>
```

`letter-spacing` adds space after every character, the spaces between words included, so a letter-spaced line keeps its words apart. For headings in capitals, a small value such as `0.05em` is common.

`-bag-horizontal-scale` draws the glyphs narrower or wider by a fixed factor, given as a percentage or a number (`90%` or `0.9`). The line breaker sees the scaled width, so more text fits on a line. The value is inherited but does not compound: `100%` inside a scaled element returns to normal width.

```xml
<Paragraph style="-bag-horizontal-scale: 90%">
    <Value>A long product name that should fit on one line</Value>
</Paragraph>
```

All text properties are listed in the [CSS reference](/reference/css-properties).

## CSS styling

All text elements support CSS styling via classes and inline styles. See [CSS and HTML](../css-html) for the full story.
