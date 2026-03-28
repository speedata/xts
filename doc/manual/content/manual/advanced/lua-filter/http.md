---
weight: 50
type: docs
linktitle: http
---

# http

The http module makes HTTP requests and returns response objects.

```lua
http = require("http")
```

## Functions

### get(url, options)

```lua
http = require("http")

response, err = http.get("https://api.example.com/products")
if not response then
    print(err)
    os.exit(-1)
end

print(response.status_code)
print(response.body)
```

### post(url, options)

```lua
response, err = http.post("https://api.example.com/data", {
    headers = {
        ["Content-Type"] = "application/json",
    },
    body = '{"name": "Widget"}',
})
```

### put(url, options), patch(url, options), delete(url, options), head(url, options)

Same signature as `get` and `post`.

### request(method, url, options)

Generic request with explicit HTTP method.

```lua
response, err = http.request("OPTIONS", "https://api.example.com/data")
```

## Options

All request functions accept an optional options table:

| Option | Description | Example |
|--------|-------------|---------|
| `headers` | Table of HTTP headers | `{["Accept"] = "text/xml"}` |
| `body` | Request body string | `'<data/>'` |
| `query` | Raw query string | `"page=2&limit=10"` |
| `cookies` | Table of cookies | `{session = "abc123"}` |
| `timeout` | Timeout in seconds or duration string | `30` or `"5s"` |
| `auth` | Basic auth table | `{user = "admin", pass = "secret"}` |

## Response object

The response object has the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `status_code` | number | HTTP status code |
| `body` | string | Response body |
| `body_size` | number | Body size in bytes |
| `url` | string | Final URL (after redirects) |
| `headers` | table | Response headers |
| `cookies` | table | Response cookies |

### Example: fetch XML data for typesetting

```lua
http = require("http")
xml = require("xml")

response, err = http.get("https://api.example.com/catalog.xml", {
    headers = {
        ["Accept"] = "application/xml",
        ["Authorization"] = "Bearer " .. os.getenv("API_TOKEN"),
    },
    timeout = 10,
})

if not response or response.status_code ~= 200 then
    runtime = require("runtime")
    runtime.log.error("Failed to fetch catalog: " .. (err or response.status_code))
    os.exit(-1)
end

-- Write response body as data.xml for the publishing run
f = io.open("data.xml", "w")
f:write(response.body)
f:close()
```
