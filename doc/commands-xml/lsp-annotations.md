# xml-lsp annotations in the layout schemas

`xtshelper genschema` writes the Relax NG schemas `schema/layoutschema-en.rng`
and `schema/layoutschema-de.rng` from `commands.xml`. Besides the usual
`a:documentation` elements, the schemas carry annotations in the namespace
`urn:xml-lsp:annotations` (prefix `lsp`). They are hints for the speedata
language server (vscode-speedata) and are ignored by Relax NG validators. The
XSD schemas do not contain them.

The rules that produce the annotations live in the `<lspannotations>` block of
`commands.xml`; only `lsp:doc` is derived from the command name. All
annotation elements are empty. The annotations are identical in the English
and the German schema, only the documentation text differs.

## Grammar level

Direct children of `<grammar>`, before `<start>`.

| Annotation | Meaning |
|---|---|
| `<lsp:namespace prefix="sd" uri="…"/>` | A namespace prefix worth offering when completing `xmlns:` on the root element. |
| `<lsp:builtin symbol="color" names="black white …"/>` | Names of a symbol that exist without a definition in the layout. `names` is space separated. |

## Element level

Direct children of `<element>`, immediately after `<a:documentation>`.

| Annotation | Meaning |
|---|---|
| `<lsp:doc href="…"/>` | Link to the reference page of the command in the online manual. Present on every command. |
| `<lsp:format preserve="true"/>`, `blank-lines="true"`, `inline="true"` | Formatter hints: keep the content verbatim, separate children with blank lines, or treat the content as inline flow. Only attributes that are set are written. |
| `<lsp:symbol kind="namespace" label="{@name}" detail="{@mode}"/>` | The element is a document symbol (outline entry). `kind` is an LSP SymbolKind name in lower camel case (`namespace`, `struct`, `function`, `method`, `object`, …). `label` and `detail` are templates in which `{@attr}` is replaced by the value of that attribute; a `detail` that is empty after substitution is dropped. |
| `<lsp:exclusive attributes="select" content="true"/>` | At most one of the listed attributes (space separated) and, with `content="true"`, the element content may be present. |
| `<lsp:when attribute="model" value="cmyk" requires="c m y k" forbids="r g b"/>` | Conditional attributes. When `attribute` has the given `value`, the attributes in `requires` must be present and the attributes in `forbids` must be absent. Without `value` the rule applies when the attribute itself is absent. A command may carry several `when` rules for the same attribute, one per value. |

## Attribute level

Direct children of `<attribute>`, immediately after `<a:documentation>` and
before the content model (`<choice>`, `<data>`). An attribute carries at most
one of them.

| Annotation | Meaning |
|---|---|
| `<lsp:defines symbol="color"/>` | The attribute value defines a name of the symbol. |
| `<lsp:references symbol="color"/>` | The attribute value refers to a name of the symbol. |

Both accept `form`:

| `form` | Value of the attribute |
|---|---|
| absent | the plain name |
| `xpath-string` | an XPath string literal, `'name'` |
| `xpath-variable` | an XPath expression that may refer to `$name` |

Symbols used by xts: `area`, `color`, `function`, `mark`, `masterpage`,
`mode`, `slate`, `template`, `variable`.
