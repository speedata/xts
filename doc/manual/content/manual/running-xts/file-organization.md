---
weight: 30
type: docs
linktitle: File Organization
---

# File Organization

How XTS finds layout files, data files, images, fonts, and stylesheets.

## Lookup rules

1. **Relative or absolute paths** -- always work directly
2. **URLs** -- fetched on the fly (no caching)
3. **Paths are relative to the current working directory**
4. **CSS paths are relative to the CSS file**, not the working directory
5. **The `--extradir` lookup** -- file names without paths are searched in these directories

## Relative and absolute paths

```xml
<!-- Relative to working directory -->
<Image href="img/ocean.pdf" width="2"/>

<!-- Absolute path -->
<Image href="/Users/myuser/assets/ocean.pdf" width="2"/>
```

## URLs

```xml
<Image href="https://placekitten.com/200/300" width="2"/>
```

## CSS path resolution

CSS `url()` references are relative to the CSS file:

![file organization](/manual/img/metaprocss.png)

```xml
<StyleSheet href="css/metapro.css"/>
```

Inside `metapro.css`:

```css
@font-face {
    font-family: "MetaPro";
    src: url("../fonts/ff-metapro-normal.otf");
}
```

The `../fonts/` path is relative to the `css/` directory.

## Adding search directories

Use `--extradir` to register directories for path-free lookups:

```
xts --extradir=/path/to/assets
```

![folder with assets](/manual/img/fileorgassets.png)

Now you can reference files by name alone:

```xml
<Image href="logo.png"/>
<Image href="jupiter.jpg"/>
```

XTS searches the extra directory and all its subdirectories recursively.

Multiple directories can be added in the config file:

```toml
extradir = ["/path/to/images", "/path/to/fonts"]
```
