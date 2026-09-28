# Contributing to XTS

Thank you for taking the time to send a patch. XTS is developed in a
private mono repository and mirrored to this one, which is why pull
requests are squash merged and the commit that lands here does not carry
your commit history. Your name stays in the squashed commit.

## Before you start

Whether a change works is only half the question. The other half is
whether XTS should carry it for years: every command and attribute stays
in the manual, the schema and the code, and someone has to maintain it.
So some changes are welcome as a pull request right away, and some need
an agreement first.

**Send a pull request directly for:**

- **A bug fix.** XTS crashes, does something other than the manual or
  the CSS specification says, or gives a different PDF from run to run
  or from platform to platform.
- **CSS as the specification describes it.** A property or value that
  XTS does not support yet, implemented as the CSS specification defines
  it. Where XTS deliberately deviates, for example with a different
  default so that existing documents do not change, say so in the pull
  request.
- **Documentation, examples and tests.**

**Open an issue first, and wait for an answer before you write the code,
for:**

- a new command, a new attribute, or a new value of an attribute in the
  layout language,
- a change of existing behavior that is not a bug fix, including a
  changed default,
- a change that needs pull requests in more than one repository
  (boxes and glue, htmlbag and XTS).

Describe the problem in the issue, not only the solution: what your
layout does today, why the existing commands cannot do it, and what you
propose. It is much cheaper to agree on the shape of a feature before
anyone implements it. "No" and "not now" are possible answers, and they
say nothing about the quality of the idea.

**Keep pull requests small**, one change each. A pull request that
depends on an unreleased change in boxes and glue or htmlbag waits until
that change is released; name the dependency at the top of the
description.

## What a change needs

A pull request that changes what a layout author can see needs three
things, ideally in the same commit:

1. **A QA case** under `qa/`: a directory with `layout.xml`, `data.xml`
   and a `reference.pdf` that shows the new or fixed behavior. Keep it
   small, one page is enough. Make the reference with

   ```
   rake "regenerateqa[mycase]"
   ```

   which runs the binary under test the way `xts compare` does, so an
   unchanged build reproduces it byte for byte. Check the PDF by eye
   before you commit it. `rake qa` must stay green. If an existing
   reference changes on purpose, regenerate it the same way and say why
   in the pull request.

2. **A changelog entry** at the top of `doc/changelog.xml`: an `<entry>`
   of its own for each change, with one sentence in `summary`, the
   details in the element text and the next release as `version`. Never
   add a second `<en>` to an existing entry, it would replace the first
   one; `go test ./helper/changelog/` refuses that. The comment at the
   top of the file describes the format. Leave `sha1` out, it is filled
   in when the commit is known. Run `rake doc` afterwards. When the
   changelog conflicts after a rebase, resolve `doc/changelog.xml` and
   run `rake doc` again rather than merging the generated page by hand.

3. **The reference documentation** when a command or attribute is added
   or changes: edit `doc/commands-xml/commands.xml`, then run
   `rake schema` and `rake doc` and commit the regenerated files under
   `schema/` and `doc/manual/`. The prose chapters under
   `doc/manual/content` are welcome too, but not required.

Pure refactorings and dependency updates need none of this.

## Building and testing

```
rake build          # bin/xts from the working tree
rake qa             # typeset every case under qa/ and compare it with its reference
go test ./...
```

`xts compare` first compares checksums. Only a PDF that differs from its
reference is rendered and compared as images, which needs ImageMagick
and Ghostscript on the PATH.

## Releases

Merged changes are collected and released together, not one release per
pull request. A change that depends on a new version of boxes and glue
or htmlbag ships with the next XTS release after that version.

## Reporting bugs

Open an issue with a minimal `layout.xml` and `data.xml` that show the
problem, the output of `xts version`, and what you expected to see. A
PDF or screenshot of the wrong output helps.
