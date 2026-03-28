---
weight: 10
type: docs
linktitle: Working with Tables
---

# Working with Tables

## Basic structure

A table is built from `<Table>`, `<Tr>` (rows), and `<Td>` (cells):

```xml
<PlaceObject>
    <Table>
        <Tr>
            <Td><Paragraph><Value>Cell 1</Value></Paragraph></Td>
            <Td><Paragraph><Value>Cell 2</Value></Paragraph></Td>
        </Tr>
        <Tr>
            <Td><Paragraph><Value>Cell 3</Value></Paragraph></Td>
            <Td><Paragraph><Value>Cell 4</Value></Paragraph></Td>
        </Tr>
    </Table>
</PlaceObject>
```

## Defining column widths

Without explicit widths, columns share the available space equally. Use `<Columns>` for control:

```xml
<Table>
    <Columns>
        <Column width="3cm"/>
        <Column width="5cm"/>
        <Column width="*"/>
    </Columns>
    <!-- rows here -->
</Table>
```

### Width options

| Value | Meaning |
|-------|---------|
| `3cm`, `50mm`, `2in` | Fixed width |
| `*` | Takes all remaining space |
| `2*` | Takes twice the share of remaining space |
| `20%` | Percentage of table width |

Proportional columns are handy:

```xml
<Columns>
    <Column width="1*"/>  <!-- 1/4 -->
    <Column width="2*"/>  <!-- 2/4 -->
    <Column width="1*"/>  <!-- 1/4 -->
</Columns>
```

## Table headers

`<Tablehead>` defines rows that repeat on every page when a table spans multiple pages:

```xml
<Table>
    <Columns>
        <Column width="2cm"/>
        <Column width="*"/>
        <Column width="3cm"/>
    </Columns>
    <Tablehead>
        <Tr>
            <Td><Paragraph><Value>ID</Value></Paragraph></Td>
            <Td><Paragraph><Value>Product</Value></Paragraph></Td>
            <Td><Paragraph><Value>Price</Value></Paragraph></Td>
        </Tr>
    </Tablehead>
    <ForAll select="product">
        <Tr>
            <Td><Paragraph><Value select="@id"/></Paragraph></Td>
            <Td><Paragraph><Value select="@name"/></Paragraph></Td>
            <Td><Paragraph><Value select="@price"/></Paragraph></Td>
        </Tr>
    </ForAll>
</Table>
```

## Spanning rows and columns

```xml
<!-- Column span -->
<Td colspan="3">
    <Paragraph><Value>This spans 3 columns</Value></Paragraph>
</Td>

<!-- Row span -->
<Td rowspan="2">
    <Paragraph><Value>Spans 2 rows</Value></Paragraph>
</Td>
```

## Styling tables with CSS

Tables support full CSS styling:

```xml
<Stylesheet>
    table { font-family: sans; font-size: 10pt; }
    thead { font-weight: bold; background-color: #e0e0e0; }
    td { padding: 4pt 8pt; border-bottom: 0.5pt solid #ccc; }
    tr:nth-child(even) { background-color: #f8f8f8; }
</Stylesheet>
```

### Cell alignment

```xml
<Stylesheet>
    .right { text-align: right; }
    .center { text-align: center; }
    .top { vertical-align: top; }
    .middle { vertical-align: middle; }
</Stylesheet>

<Td class="right"><Paragraph><Value>99.95</Value></Paragraph></Td>
```

### Tabular numbers for financial data

```xml
<Stylesheet>
    .numbers {
        text-align: right;
        font-feature-settings: "tnum", "lnum";
    }
</Stylesheet>
```

## Cell content

Table cells can hold more than just text:

```xml
<!-- Multiple paragraphs -->
<Td>
    <Paragraph><Value>First line</Value></Paragraph>
    <Paragraph><Value>Second line</Value></Paragraph>
</Td>

<!-- Images -->
<Td>
    <Image file="logo.pdf" width="2cm"/>
</Td>

<!-- Nested tables -->
<Td>
    <Table>
        <Tr><Td><Paragraph><Value>Nested</Value></Paragraph></Td></Tr>
    </Table>
</Td>

<!-- Lists -->
<Td>
    <Ul>
        <Li><Value>Item 1</Value></Li>
        <Li><Value>Item 2</Value></Li>
    </Ul>
</Td>
```

Individual cells are never split across pages -- they're always rendered as a single rectangular box.

## Data-driven tables

Generate tables from XML data using `<ForAll>`:

```xml title="data.xml"
<inventory>
    <item sku="A001" name="Widget" qty="150" price="9.99"/>
    <item sku="A002" name="Gadget" qty="75" price="24.99"/>
    <item sku="A003" name="Gizmo" qty="200" price="14.99"/>
</inventory>
```

```xml title="layout.xml"
<Stylesheet>
    table { font-family: sans; font-size: 10pt; }
    thead { font-weight: bold; background-color: #333; color: white; }
    td { padding: 4pt 8pt; border-bottom: 0.5pt solid #ddd; }
    .right { text-align: right; }
</Stylesheet>

<Record element="inventory">
    <PlaceObject>
        <Table>
            <Columns>
                <Column width="2cm"/>
                <Column width="*"/>
                <Column width="2cm"/>
                <Column width="2.5cm"/>
            </Columns>
            <Tablehead>
                <Tr>
                    <Td><Paragraph><Value>SKU</Value></Paragraph></Td>
                    <Td><Paragraph><Value>Product</Value></Paragraph></Td>
                    <Td class="right"><Paragraph><Value>Qty</Value></Paragraph></Td>
                    <Td class="right"><Paragraph><Value>Price</Value></Paragraph></Td>
                </Tr>
            </Tablehead>
            <ForAll select="item">
                <Tr>
                    <Td><Paragraph><Value select="@sku"/></Paragraph></Td>
                    <Td><Paragraph><Value select="@name"/></Paragraph></Td>
                    <Td class="right"><Paragraph><Value select="@qty"/></Paragraph></Td>
                    <Td class="right"><Paragraph><Value select="concat(@price, ' EUR')"/></Paragraph></Td>
                </Tr>
            </ForAll>
        </Table>
    </PlaceObject>
</Record>
```

## Conditional row styling

Apply different styles based on data values:

```xml
<Stylesheet>
    .low-stock { background-color: #ffcccc; }
    .in-stock { background-color: #ccffcc; }
</Stylesheet>

<ForAll select="item">
    <Tr class="{ if (@qty &lt; 100) then 'low-stock' else 'in-stock' }">
        <!-- cells -->
    </Tr>
</ForAll>
```

Or use `<Switch>/<Case>` inside cells for more complex logic:

```xml
<Td>
    <Switch>
        <Case test="@qty &lt; 100">
            <Paragraph style="color: red;"><Value select="@qty"/></Paragraph>
        </Case>
        <Otherwise>
            <Paragraph><Value select="@qty"/></Paragraph>
        </Otherwise>
    </Switch>
</Td>
```

## Table width

```xml
<!-- Fixed width -->
<Table width="15cm">...</Table>

<!-- Grid cells -->
<Table width="10">...</Table>

<!-- Full available width via CSS -->
<Stylesheet>
    table { width: 100%; }
</Stylesheet>
```

## See also

- [Table reference](/reference/commands/table), [Tr](/reference/commands/tr), [Td](/reference/commands/td), [Columns](/reference/commands/columns), [Tablehead](/reference/commands/tablehead)
