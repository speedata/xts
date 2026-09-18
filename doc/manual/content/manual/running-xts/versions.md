---
weight: 60
type: docs
linktitle: Versions
---

# Versions and Compatibility

XTS is released with version numbers of the form `0.MINOR.PATCH`, for
example `0.1.2`. The leading zero says that the layout language is still
being shaped and may change between minor versions.

## What a version number tells you

- **PATCH** (`0.1.1` to `0.1.2`): bug fixes only. A layout that worked
  with `0.1.1` produces the same document with `0.1.2`, apart from the
  bug that was fixed.
- **MINOR** (`0.1` to `0.2`): new commands, attributes and functions, and
  changes in behaviour. A minor release may change line breaks or page
  breaks of an existing layout, and it may remove features that were
  marked as deprecated in the previous minor release.

Every change that you can notice in a layout is listed in the
[changelog](../../../changelog), including the ones that alter the
output of existing layouts. Read it before you update XTS in a
production setup, then run your own [comparison tests](../quality-assurance).

## Deprecation

A command, attribute or function that is going to be removed is first
marked as deprecated: it keeps working for one minor release and XTS
writes a warning that names the replacement. The next minor release
removes it and the layout stops with an error. The changelog entry for
the deprecation says which release removes the feature.

## Declaring the version in the layout

The `<Layout>` element accepts a `version` attribute with the XTS
version the layout was written for:

```xml
<Layout xmlns="urn:speedata.de/2021/xts/en"
        xmlns:sd="urn:speedata.de/2021/xtsfunctions/en"
        version="0.1">
```

XTS compares the parts from the left with its own version, so `0.1` is
satisfied by 0.1.0, 0.1.5 and 0.2.0. A layout that asks for a newer
version than the running XTS stops with an error instead of producing a
document that silently lacks the features the layout relies on. An older
version passes, so the attribute records the minimum, not an exact match.
Development builds of XTS accept every version.

## Finding the version

```
xts version
```

prints the version of the binary. The [releases page](https://github.com/speedata/xts/releases)
lists every release with its changes.
