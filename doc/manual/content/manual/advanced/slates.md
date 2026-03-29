---
weight: 30
type: docs
linktitle: Slates
---

# Slates

<style>
.slate-images { display: flex; gap: 1rem; align-items: center; flex-wrap: wrap; }
.slate-images img { width: 45%; }
</style>

<figure>
<div class="slate-images">
<img src="/manual/advanced/img/slate.webp" alt="A classic school slate with wooden frame" />
<img src="/manual/advanced/img/magic-slate.webp" alt="A magic drawing slate (JIKKY) with drawings" />
</div>
<figcaption>From school slate to magic drawing tablet -- the idea behind Slates in XTS: draw content on an independent surface, measure it, place it on the page, or simply discard it.<br/>
<small>Left: Hannes Grobe, <a href="https://commons.wikimedia.org/wiki/File:Slate_hg.jpg">Wikimedia Commons</a>, CC BY 3.0. Right: Tatsuo Yamashita, CC BY 2.0.</small>
</figcaption>
</figure>

A Slate is a virtual layout surface -- like a magic drawing tablet, you sketch content on it without it appearing on the page, inspect the result, and then place it wherever you want. A slate has its own grid and cursor, independent of the page.

## Creating and placing a slate

```xml
<Slate name="sidebar">
    <Grid width="5mm" height="12pt"/>
    <Contents>
        <PlaceObject>
            <TextBlock>
                <Paragraph><Value>Sidebar content</Value></Paragraph>
            </TextBlock>
        </PlaceObject>
    </Contents>
</Slate>

<!-- Place the slate on the page -->
<PlaceObject slate="sidebar"/>
```

## Why use slates?

- **Independent grid**: A slate can use a finer or coarser grid than the page.
- **Measure before placing**: Use `sd:slate-width('name')` and `sd:slate-height('name')` to query a slate's dimensions before deciding where to put it.
- **Reuse**: Place the same slate multiple times.
- **Discard**: If the content doesn't fit or isn't needed, simply don't place it -- nothing ends up in the PDF.

## Querying slate dimensions

```xml
<Slate name="card">
    <Contents>
        <!-- build the card content -->
    </Contents>
</Slate>

<!-- Check if it fits -->
<Switch>
    <Case test="sd:slate-height('card', 'cm') &lt; 5">
        <PlaceObject slate="card"/>
    </Case>
    <Otherwise>
        <ClearPage/>
        <PlaceObject slate="card"/>
    </Otherwise>
</Switch>
```

## See also

- [Slate reference](/reference/commands/slate)
- [Contents reference](/reference/commands/contents)
