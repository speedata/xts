---
weight: 10
type: docs
linktitle: runtime
---

# runtime

The runtime module provides access to project settings, XSLT transformations, file lookup, and external commands.

```lua
runtime = require("runtime")
```

## Properties

### projectdir

The current project directory as a string.

```lua
print(runtime.projectdir)
-- /home/user/myproject
```

### variables

A table of variables passed via the `-v` command line flag. Supports read and write.

```lua
-- xts -v myvar=hello
print(runtime.variables.myvar)  -- "hello"

runtime.variables.output = "draft"
```

### options

Configuration settings (read/write). These correspond to entries in the configuration file.

```lua
print(runtime.options.mode)

runtime.options.runs = "2"
```

### finalizer

A callback function that runs after PDF creation.

```lua
runtime.finalizer = function()
    print("PDF created successfully")
end
```

### log

Logging functions at different levels.

```lua
runtime.log.debug("detailed info")
runtime.log.info("processing started")
runtime.log.warn("missing optional data")
runtime.log.error("something went wrong")
```

## Functions

### execute(args)

Runs an external program. The first entry in the table is the command, followed by arguments.

```lua
ok, msg = runtime.execute({"ls", "-l", "data"})
if not ok then
    print(msg)
end
```

### find_file(name)

Finds a file in the XTS search path. Returns the absolute path or `nil` if not found.

```lua
path = runtime.find_file("layout.xml")
if path then
    print("Found: " .. path)
end
```
