---
weight: 20
type: docs
linktitle: Lua Filter
---

# Lua Filter / Preprocessing

Sometimes you need to transform data before creating the PDF -- convert CSV to XML, validate input, or run an XSLT transformation. XTS supports Lua scripts that run *before* the publishing process.

## Running a Lua script

Via command line:

```
xts --filter myfile.lua
```

Or via configuration file:

```toml
filter = "myfile.lua"
```

The script executes before PDF generation. It must be within the XTS search path.

## Available modules

- [runtime](runtime) -- Project settings, file lookup, external commands
- [xml](xml) -- Create, read, and transform XML files
- [csv](csv) -- Read CSV files
- [xlsx](xlsx) -- Read Excel spreadsheets
- [http](http) -- HTTP requests
