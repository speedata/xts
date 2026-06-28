---
weight: 30
type: docs
linktitle: Programming
aliases:
  - /manual/data-processing/
---

# Programming

XTS is a programming language for typesetting. A layout file is not a static
template that gets filled in: it is a program that runs, command by command, as
the PDF is built. This section explains that execution model and the language
features that go with it -- the data/action distinction, typed values, variables,
XPath, functions, templates, and control flow.

If you are looking for *how to make a page look a certain way* (fonts, tables,
images, the grid), see the [Layout Guide](/manual). If you want the terse,
per-command facts, see the [Reference](/reference).

## The mental model first

Two chapters set up the way of thinking that the rest build on. Read them first.

- [Execution model](/programming/execution-model) -- why XTS runs commands in order, keeps live
  layout state, and is therefore imperative (not a one-shot tree-to-PDF renderer).
- [Data vs. action](/programming/data-and-actions) -- every command either *builds a value* or
  *performs an action*. Knowing which is which explains almost everything else.

## The language

- [Values and types](/programming/values-and-types) -- `as=` sequence types, the queryable data
  band, and when things are evaluated.
- [Variables](/programming/variables) -- setting, reading, scope, evaluation time, building up
  content, dynamic names.
- [XPath basics](/programming/xpath-basics) -- a practical primer on the XPath expression
  language: paths, predicates, operators, functions.
- [XPath in XTS](/programming/xpath) -- where XPath is used and the `sd:` layout functions.
- [Maps and arrays](/programming/maps-and-arrays) -- XPath 3.1 structured data: maps, arrays,
  and the `?` lookup operator.
- [Functions](/programming/functions) -- reusable, side-effect-free `<Function>` definitions that
  return data.
- [Templates](/programming/templates) -- reusable *behaviour* with `<Template>` / `<CallTemplate>`,
  including layout actions.
- [Control flow](/programming/control-flow) -- `ForAll`, `Loop`, `While`/`Until`, `Switch`/`Case`.
- [Records and dispatch](/programming/records-and-dispatch) -- matching data elements to records,
  modes, and priority.
- [Data files](/programming/data-files) -- structuring the XML you feed to XTS.
