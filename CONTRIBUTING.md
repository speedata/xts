# Contributing to XTS

Thank you for taking the time to send a patch. XTS is developed in a
private mono repository and mirrored to this one, which is why pull
requests are squash merged and the commit that lands here does not carry
your commit history. Your name stays in the squashed commit.

## What a change needs

A pull request that changes what a layout author can see needs three
things, ideally in the same commit:

1. **A QA case** under `qa/`: a directory with `layout.xml`, `data.xml`
   and a `reference.pdf` that shows the new or fixed behaviour. Keep it
   small, one page is enough. Generate the reference with the binary
   under test and check the PDF by eye before you commit it:

   ```
   rake build
   cd qa/mycase && ../../bin/xts --suppressinfo --jobname reference && ../../bin/xts --jobname reference clean
   ```

   `rake qa` must stay green. If an existing reference changes on
   purpose, say why in the pull request.

2. **A changelog entry** at the top of `doc/changelog.xml`: one sentence
   in `summary`, the details in the element text, the next release as
   `version`. The comment at the top of the file describes the format.
   Leave `sha1` out, it is filled in when the commit is known.

3. **The reference documentation** when a command or attribute is added
   or changes: edit `doc/commands-xml/commands.xml`, then run
   `rake schema` and `rake doc` and commit the regenerated files under
   `schema/` and `doc/manual/`. The prose chapters under
   `doc/manual/content` are welcome too, but not required.

Pure refactorings and dependency updates need none of this.

## Building and testing

```
rake build          # bin/xts from the working tree
rake qa             # render every case under qa/ and compare with its reference
go test ./...
```

`xts compare` needs ImageMagick and Ghostscript on the PATH.

## Reporting bugs

Open an issue with a minimal `layout.xml` and `data.xml` that show the
problem, the output of `xts version`, and what you expected to see. A
PDF or screenshot of the wrong output helps.
