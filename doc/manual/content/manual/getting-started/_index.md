---
weight: 10
type: docs
linktitle: Getting Started
---

# Getting Started

So you want to turn XML data into beautiful PDFs? You're in the right place. This chapter gets you up and running in a few minutes: install XTS, create a tiny project, and see your first PDF appear.

## What is XTS?

XTS is a non-interactive tool that takes two inputs -- a **data file** (XML) and a **layout file** (also XML) -- and produces a PDF. There's no graphical interface. You describe *what goes where* in the layout, point XTS at your data, and it does the rest.

![XML to PDF schema](/images/xmltopdf.png)

This strict separation of data and layout is the key idea. Your data can come from a database, a PIM system, an API -- anything that can produce XML. The layout file is where you define page sizes, grids, fonts, tables, and all the visual design. XTS combines the two and generates the PDF fully automatically.

**Typical use cases:**

- Product catalogs
- Price lists and data sheets
- Travel guides
- Invoices and reports
- Any document where content changes but the design stays the same

## Installing XTS

### Pre-built packages (recommended)

1. Go to the [download page](https://github.com/speedata/xts/releases/latest) and grab the package for your operating system.
2. Unzip it anywhere you like -- no admin rights needed.
3. Add the `bin/` directory to your `PATH`.

That's it. You should now be able to run `xts` in your terminal.

### Building from source

If you prefer to build from source:

```
git clone https://github.com/speedata/xts.git
cd xts
rake build
```

The `xts` binary ends up in `bin/`. You'll need Go 1.21+ and Ruby/rake (or check the `Rakefile` for the raw `go build` command).

## Hello, World!

Let's create your first PDF. The fastest way:

```
xts new helloworld
cd helloworld
xts
```

This creates a directory with two files and runs XTS on them. Open `xts.pdf` and you should see "Hello, user!". But let's look at what's inside.

**data.xml** -- the data:
```xml
<data name="user"></data>
```

Just a root element with a `name` attribute. In a real project, this would be your product data, article list, or whatever you're typesetting.

**layout.xml** -- the layout:
```xml
<Layout xmlns="urn:speedata.de/2021/xts/en"
    xmlns:sd="urn:speedata.de/2021/xtsfunctions/en">

    <Record element="data">
        <PlaceObject>
            <Textblock>
                <Paragraph>
                    <Value select="concat('Hello, ', @name, '!')" />
                </Paragraph>
            </Textblock>
        </PlaceObject>
    </Record>
</Layout>
```

Here's what's happening:

1. `<Layout>` is the root element. The two `xmlns` declarations set up the XTS namespace and the XPath function namespace (`sd:`).
2. `<Record element="data">` says: "When you encounter a `<data>` element in the data file, execute these commands."
3. `<PlaceObject>` places something on the page.
4. `<Textblock>` is a rectangular text area (no page breaks).
5. `<Paragraph>` holds one paragraph of text.
6. `<Value select="...">` evaluates an XPath expression -- here it concatenates "Hello, " with the `name` attribute and "!".

Run `xts` in the directory, and you get your PDF:

![hello world](img/helloworld.png)

Congratulations -- that's your first XTS document!

## Examples

There's a separate repository with more complete examples on [GitHub](https://github.com/speedata/xts-examples). Browse through them to see what's possible: data-driven tables, multi-page layouts, custom fonts, and more.

## What's next?

Now that you have XTS running, head to [Core Concepts](../core-concepts) to understand how the grid, positioning, and data processing work together.
