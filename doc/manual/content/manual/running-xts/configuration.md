---
weight: 20
type: docs
linktitle: Configuration
---

# Configuration File

Instead of passing flags on the command line every time, you can create a `xts.cfg` file in [TOML format](https://toml.io/en/):

```toml title="xts.cfg"
# Use a custom data file
data = "products.xml"
jobname = "catalog"
runs = 2

# Pass variables to the layout
[variables]
edition = "spring"
year = 2026
```

## All configuration keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `data` | string | `"data.xml"` | Data file name |
| `layout` | string | `"layout.xml"` | Layout file name |
| `dummy` | boolean | `false` | Use `<data/>` instead of data file |
| `extradir` | string[] | `[]` | Additional search directories |
| `filter` | string | `""` | Lua script to run before publishing |
| `jobname` | string | `"xts"` | Output file name (without `.pdf`) |
| `loglevel` | string | `"info"` | Log level: debug, info, warn, error |
| `runs` | integer | `1` | Number of publishing runs |
| `systemfonts` | boolean | `false` | Use system fonts |
| `quiet` | boolean | `false` | Suppress console output |
| `verbose` | boolean | `false` | Extra debug output |
| `suppressinfo` | boolean | `false` | Reproducible PDF output |
| `trace` | string[] | `[]` | Traces: `"grid"`, `"gridallocation"` |
| `mode` | string[] | `[]` | Active modes |
| `variables` | table | | Key-value pairs accessible in the layout |

Variables defined here are accessible in the layout via the command line variable mechanism.
