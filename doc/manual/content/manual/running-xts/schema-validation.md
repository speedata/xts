---
weight: 40
type: docs
linktitle: Schema Validation
---

# Schema Validation

XTS ships with XML schema files (RELAX NG and XSD) that describe the layout language. A good XML editor can use these schemas to provide auto-complete, validation, and inline documentation -- making layout editing much faster and less error-prone.

## What schemas give you

- **Auto-complete**: Type `<` and see all valid commands
- **Attribute suggestions**: See which attributes are available and their allowed values
- **Inline documentation**: Read command descriptions without leaving your editor
- **Instant validation**: Catch syntax errors as you type

## Schema files

The schema files are in the `share/schema/` directory of your XTS installation:

| File | Description |
|------|-------------|
| `layoutschema-en.rng` | RELAX NG, English documentation |
| `layoutschema-de.rng` | RELAX NG, German documentation |
| `layoutschema-en.xsd` | XML Schema, English documentation |
| `layoutschema-de.xsd` | XML Schema, German documentation |

## Visual Studio Code (recommended)

The easiest way to get schema support in VS Code is the **speedata** extension. It supports both the speedata Publisher and XTS layout files out of the box -- no manual schema configuration needed.

1. Open the Extensions marketplace and search for `speedata`.
2. Install the **Speedata Publisher** extension (`vscode-speedata`). It supports XTS as well.
3. That's it. Layout files with the XTS namespace are automatically recognized.

![speedata extension in VS Code marketplace](img/vscode-speedata-extension.png)

The extension provides auto-complete, validation, and inline documentation based on RELAX NG.

### Alternative: Red Hat XML extension

If you prefer the generic XML extension, you can configure it manually:

1. Install the **XML** extension by Red Hat from the marketplace.
2. Open VS Code settings and find `xml.catalogs`.
3. Add the path to `catalog-schema-en.xml` from your XTS installation.

![VS Code XML extension](img/vscode-xml-redhat.png)

![VS Code catalog setting](img/vscode-xml-catalog.png)

Once configured, any layout file with the namespace `urn:speedata.de/2021/xts/en` gets full auto-complete:

![VS Code auto-complete](img/vscode-sample-layout.png)

## oXygen XML

1. Open **Preferences > Document Type Association**.
2. Click **New** to create a new association.
3. Add a namespace rule for `urn:speedata.de/2021/xts/en`.
4. Set the schema type to RELAX NG and select `layoutschema-en.rng`.

![oXygen namespace](img/oxygen-namespace.png)

![oXygen schema](img/oxygen-schema.png)

After setup, you get auto-complete for both elements and attributes:

![oXygen elements](img/oxygen-elements.png)

![oXygen attributes](img/oxygen-attributes.png)

## Other editors

Any XML editor that supports RELAX NG or XSD can use the XTS schemas:

- [XMLSpy](https://www.altova.com/xml-editor/) (Windows)
- [XML Blueprint](https://www.xmlblueprint.com/) (Windows)
- [GNU Emacs](https://www.gnu.org/software/emacs/) with nxml-mode (cross-platform, free)
- [jEdit](http://www.jedit.org) (cross-platform, free)
