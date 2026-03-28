---
weight: 30
type: docs
linktitle: CSS and HTML
---

# CSS and HTML

HTML markup and CSS styling are first-class citizens in XTS. You can use CSS to style XTS layout elements (like `<Paragraph>`) and HTML content within `<HTML>` blocks. If you know web CSS, you already know most of what you need.

## CSS stylesheets

### Inline stylesheets

Define CSS rules directly in your layout file:

```xml
<Stylesheet>
    p {
        font-family: serif;
        font-size: 12pt;
        line-height: 1.4;
    }
    .highlight {
        background-color: yellow;
    }
</Stylesheet>
```

### External stylesheets

For larger projects, keep CSS in separate files:

```xml
<Stylesheet href="styles/main.css"/>
<Stylesheet href="styles/typography.css"/>
```

You can include as many as you need. They're processed in order, so later rules override earlier ones (normal CSS cascade).

## CSS classes

Apply classes with the `class` attribute -- works on both XTS elements and HTML elements:

```xml
<Stylesheet>
    .product-title {
        font-size: 18pt;
        font-weight: bold;
        color: darkblue;
    }
    .price {
        font-family: monospace;
        color: green;
    }
</Stylesheet>

<PlaceObject>
    <Textblock>
        <Paragraph class="product-title">
            <Value select="@name"/>
        </Paragraph>
        <Paragraph class="price">
            <Value select="concat(@price, ' EUR')"/>
        </Paragraph>
    </Textblock>
</PlaceObject>
```

Multiple classes are space-separated: `class="centered bold large"`.

## Inline styles

For one-off styling:

```xml
<Paragraph style="color: red; font-weight: bold;">
    <Value>Warning!</Value>
</Paragraph>
```

## Supported CSS properties

Here's a quick overview of what works. For the full list, see [CSS Properties Reference](/reference/css-properties).

**Text:** `font-family`, `font-size`, `font-weight`, `font-style`, `color`, `text-align`, `text-indent`, `line-height`, `font-feature-settings`

**Box model:** `margin`, `padding`, `border`, `border-radius`, `background-color`

**Borders per side:** `border-top`, `border-bottom`, `border-left`, `border-right`

## HTML content

The `<HTML>` element lets you include HTML markup:

```xml
<PlaceObject>
    <Textblock>
        <HTML>
            <p>A wonderful <b>serenity</b> has taken possession
               <i>of my <b>entire soul,</b></i> like these sweet mornings.</p>
        </HTML>
    </Textblock>
</PlaceObject>
```

### Supported HTML elements

| Element | Description |
|---------|-------------|
| `<p>` | Paragraph |
| `<b>`, `<strong>` | Bold |
| `<i>`, `<em>` | Italic |
| `<u>` | Underline |
| `<span>` | Inline container |
| `<br>` | Line break |
| `<a>` | Link |
| `<ul>`, `<ol>`, `<li>` | Lists |
| `<table>`, `<tr>`, `<td>`, `<th>` | Tables |
| `<h1>` -- `<h6>` | Headings |
| `<div>` | Block container |
| `<pre>`, `<code>` | Preformatted / code |
| `<img>` | Image |

### Dynamic content with expand-text

By default, `{...}` expressions are *not* evaluated inside HTML. Enable them with `expand-text="yes"`:

```xml
<SetVariable variable="product" select="'Premium Widget'"/>

<PlaceObject>
    <Textblock>
        <HTML expand-text="yes">
            <p><b>{$product}</b></p>
            <p>Page: {sd:current-page()}</p>
        </HTML>
    </Textblock>
</PlaceObject>
```

For literal curly braces, double them: `{{` and `}}`.

## Loading HTML from data

Pull HTML from your data file using `select`:

```xml
<!-- Reference HTML nodes directly -->
<HTML select="/data/htmlcontent"/>

<!-- Or decode raw HTML strings -->
<HTML select="sd:decode-html(description)"/>
```

## Combining XTS elements and HTML

Mix freely within a `<Textblock>`:

```xml
<Textblock>
    <Paragraph>
        <Value>Introduction: </Value>
        <HTML><b>Important</b> information follows.</HTML>
    </Paragraph>
    <HTML>
        <ul>
            <li>First item</li>
            <li>Second item</li>
        </ul>
    </HTML>
</Textblock>
```

## Practical example

A complete product catalog card:

```xml
<Stylesheet>
    body { font-family: serif; font-size: 11pt; }
    h1 { font-size: 18pt; color: darkblue; margin-bottom: 12pt; }
    .product { border: 1pt solid #ccc; padding: 10pt; margin-bottom: 10pt; }
    .product-name { font-weight: bold; font-size: 14pt; }
    .price { color: green; font-feature-settings: "tnum"; }
    .description { font-style: italic; color: #666; }
</Stylesheet>

<Record element="products">
    <PlaceObject>
        <Textblock>
            <HTML><h1>Product Catalog</h1></HTML>
            <ForAll select="product">
                <HTML expand-text="yes">
                    <div class="product">
                        <p class="product-name">{@name}</p>
                        <p class="price">{@price} EUR</p>
                        <p class="description">{description}</p>
                    </div>
                </HTML>
            </ForAll>
        </Textblock>
    </PlaceObject>
</Record>
```
