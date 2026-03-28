---
weight: 30
type: docs
linktitle: csv
---

# csv

The csv module reads CSV files into Lua tables.

```lua
csv = require("csv")
```

## Functions

### decode(filename, options)

Reads a CSV file and returns a table of rows. Each row is a table of column values.

**Options:**

| Option | Description | Default |
|--------|-------------|---------|
| `separator` | Column separator character | `,` |
| `charset` | Character encoding (`"ISO-8859-1"` supported) | UTF-8 |
| `columns` | Table of column indices to extract | all columns |

```lua
csv = require("csv")

-- Read all columns
data = csv.decode("products.csv")

for i, row in ipairs(data) do
    print(row[1], row[2], row[3])
end
```

**With options:**

```lua
csv = require("csv")

data, msg = csv.decode("products.csv", {
    separator = ";",
    charset = "ISO-8859-1",
    columns = {1, 3, 5}  -- only columns 1, 3, and 5
})

if not data then
    print(msg)
    os.exit(-1)
end

-- data[1][1] = first row, first selected column (original column 1)
-- data[1][2] = first row, second selected column (original column 3)
-- data[1][3] = first row, third selected column (original column 5)
```
