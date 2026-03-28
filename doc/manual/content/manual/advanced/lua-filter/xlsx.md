---
weight: 40
type: docs
linktitle: xlsx
---

# xlsx

The xlsx module reads Microsoft Excel files (.xlsx).

```lua
xlsx = require("xlsx")
```

## Functions

### open(filename)

Opens an Excel file and returns a spreadsheet object.

```lua
xlsx = require("xlsx")

spreadsheet, err = xlsx.open("report.xlsx")
if not spreadsheet then
    print(err)
    os.exit(-1)
end

-- Number of worksheets
print(#spreadsheet)

-- Access first worksheet (1-based)
ws = spreadsheet[1]
```

### string_to_date(string)

Converts an Excel date number (as string) to a table with date/time fields.

```lua
xlsx = require("xlsx")

d = xlsx.string_to_date("45678")
print(d.year, d.month, d.day)
print(d.hour, d.minute, d.second)
```

## Worksheet object

A worksheet provides cell access and metadata.

### Reading cells

Call the worksheet with `(row, column)` to read a cell value (1-based).

```lua
ws = spreadsheet[1]

-- Cell A1 (row 1, column 1)
val = ws(1, 1)

-- Cell C5 (row 5, column 3)
val = ws(5, 3)
```

### Properties

| Property | Description |
|----------|-------------|
| `ws.name` | Worksheet name |
| `ws.minrow` | First row with data |
| `ws.maxrow` | Last row with data |
| `ws.mincol` | First column with data |
| `ws.maxcol` | Last column with data |

### Example: iterate all cells

```lua
xlsx = require("xlsx")
xml = require("xml")

spreadsheet, err = xlsx.open("products.xlsx")
if not spreadsheet then
    print(err)
    os.exit(-1)
end

ws = spreadsheet[1]

-- Build XML from spreadsheet data
root = {
    type = "element",
    name = "data",
}

for row = ws.minrow, ws.maxrow do
    local item = {
        type = "element",
        name = "row",
    }
    for col = ws.mincol, ws.maxcol do
        item[#item + 1] = {
            type = "element",
            name = "cell",
            ws(row, col),
        }
    end
    root[#root + 1] = item
end

xml.encode_table(root, "data.xml")
```
