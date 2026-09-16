package genschema

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// LSPNAMESPACE is the namespace for xml-lsp schema annotations. The
// annotations are consumed by the speedata language server (vscode-speedata)
// and ignored by ordinary Relax NG validators.
const LSPNAMESPACE string = "urn:xml-lsp:annotations"

// lspSymbolXML is a <defines> or <references> rule from the <lspannotations>
// block. Command and Attribute are space separated lists of names.
type lspSymbolXML struct {
	Command   string `xml:"command,attr"`
	Attribute string `xml:"attribute,attr"`
	Symbol    string `xml:"symbol,attr"`
	Form      string `xml:"form,attr"`
}

// lspBuiltinXML lists names of a symbol that exist without being defined in
// the layout, such as the CSS color names.
type lspBuiltinXML struct {
	Symbol string `xml:"symbol,attr"`
	Names  string `xml:"names,attr"`
}

// lspFormatXML carries formatter hints for a command.
type lspFormatXML struct {
	Command    string `xml:"command,attr"`
	Preserve   string `xml:"preserve,attr"`
	BlankLines string `xml:"blank-lines,attr"`
	Inline     string `xml:"inline,attr"`
}

// lspDocSymbolXML turns a command into a document symbol (outline entry).
type lspDocSymbolXML struct {
	Command string `xml:"command,attr"`
	Kind    string `xml:"kind,attr"`
	Label   string `xml:"label,attr"`
	Detail  string `xml:"detail,attr"`
}

// lspExclusiveXML states that the listed attributes and the element content
// are mutually exclusive.
type lspExclusiveXML struct {
	Command    string `xml:"command,attr"`
	Attributes string `xml:"attributes,attr"`
	Content    string `xml:"content,attr"`
}

// lspWhenXML is a conditional attribute rule: when Attribute has one of the
// values in Value (or is absent, if Value is empty), the attributes in
// Requires must be present and the attributes in Forbids must not.
type lspWhenXML struct {
	Command   string `xml:"command,attr"`
	Attribute string `xml:"attribute,attr"`
	Value     string `xml:"value,attr"`
	Requires  string `xml:"requires,attr"`
	Forbids   string `xml:"forbids,attr"`
}

// lspNamespaceXML is a namespace prefix the editor should offer for xmlns
// completion.
type lspNamespaceXML struct {
	Prefix string `xml:"prefix,attr"`
	URI    string `xml:"uri,attr"`
}

type lspAnnotationsXML struct {
	Defines    []lspSymbolXML    `xml:"defines"`
	References []lspSymbolXML    `xml:"references"`
	Builtins   []lspBuiltinXML   `xml:"builtin"`
	Formats    []lspFormatXML    `xml:"format"`
	Symbols    []lspDocSymbolXML `xml:"symbol"`
	Exclusives []lspExclusiveXML `xml:"exclusive"`
	Whens      []lspWhenXML      `xml:"when"`
	Namespaces []lspNamespaceXML `xml:"namespace"`
}

// containsField reports whether name is one of the space separated fields in
// list.
func containsField(list, name string) bool {
	for _, f := range strings.Fields(list) {
		if f == name {
			return true
		}
	}
	return false
}

// matches reports whether the rule applies to the attribute of the command.
// With specific set, only rules that name a command are considered, otherwise
// only rules without a command.
func (r *lspSymbolXML) matches(command, attribute string, specific bool) bool {
	if specific != (r.Command != "") {
		return false
	}
	if specific && !containsField(r.Command, command) {
		return false
	}
	return containsField(r.Attribute, attribute)
}

// lspSymbolAnnotation returns the annotation kind ("defines" or
// "references"), the symbol and the form for the attribute of the command, or
// an empty kind if no rule matches. Rules naming a command win over generic
// rules, so for example Mark/select keeps its defines annotation while select
// on all other commands is a variable reference. Each attribute gets at most
// one annotation.
func (c *commandsXML) lspSymbolAnnotation(command, attribute string) (kind, symbol, form string) {
	for _, specific := range []bool{true, false} {
		for i := range c.LspAnnotations.Defines {
			if r := &c.LspAnnotations.Defines[i]; r.matches(command, attribute, specific) {
				return "defines", r.Symbol, r.Form
			}
		}
		for i := range c.LspAnnotations.References {
			if r := &c.LspAnnotations.References[i]; r.matches(command, attribute, specific) {
				return "references", r.Symbol, r.Form
			}
		}
	}
	return "", "", ""
}

// lspFormatAnnotation returns the format rule for the command or nil.
func (c *commandsXML) lspFormatAnnotation(command string) *lspFormatXML {
	for i := range c.LspAnnotations.Formats {
		if containsField(c.LspAnnotations.Formats[i].Command, command) {
			return &c.LspAnnotations.Formats[i]
		}
	}
	return nil
}

// lspDocSymbolAnnotation returns the document symbol rule for the command or
// nil.
func (c *commandsXML) lspDocSymbolAnnotation(command string) *lspDocSymbolXML {
	for i := range c.LspAnnotations.Symbols {
		if containsField(c.LspAnnotations.Symbols[i].Command, command) {
			return &c.LspAnnotations.Symbols[i]
		}
	}
	return nil
}

// lspExclusiveAnnotations returns all exclusive rules for the command.
func (c *commandsXML) lspExclusiveAnnotations(command string) []*lspExclusiveXML {
	var ret []*lspExclusiveXML
	for i := range c.LspAnnotations.Exclusives {
		if containsField(c.LspAnnotations.Exclusives[i].Command, command) {
			ret = append(ret, &c.LspAnnotations.Exclusives[i])
		}
	}
	return ret
}

// lspWhenAnnotations returns all conditional attribute rules for the command.
func (c *commandsXML) lspWhenAnnotations(command string) []*lspWhenXML {
	var ret []*lspWhenXML
	for i := range c.LspAnnotations.Whens {
		if containsField(c.LspAnnotations.Whens[i].Command, command) {
			ret = append(ret, &c.LspAnnotations.Whens[i])
		}
	}
	return ret
}

// lspDocURL returns the URL of the reference page of the command in the
// online manual. The manual is English only, so both schema languages get
// the same link.
func lspDocURL(command string) string {
	return DOCBASE + strings.ToLower(command) + "/"
}

// hasAttribute reports whether the command declares the attribute.
func (c *commandsxmlCommand) hasAttribute(name string) bool {
	for _, a := range c.Attributes {
		if a.Name == name {
			return true
		}
	}
	return false
}

// warnUnmatchedLspRules prints a warning for each rule in <lspannotations>
// that names a command or attribute that does not exist. Such a rule is
// almost always a typo and would otherwise silently produce no annotation.
func warnUnmatchedLspRules(c *commandsXML) {
	byName := map[string]*commandsxmlCommand{}
	for i := range c.Commands {
		byName[c.Commands[i].Name] = &c.Commands[i]
	}
	warn := func(format string, args ...any) {
		fmt.Fprintf(os.Stderr, "genschema: "+format+"\n", args...)
	}
	checkCommands := func(rule, list string) []*commandsxmlCommand {
		var ret []*commandsxmlCommand
		for _, name := range strings.Fields(list) {
			cmd, ok := byName[name]
			if !ok {
				warn("lsp %s rule names unknown command %q", rule, name)
				continue
			}
			ret = append(ret, cmd)
		}
		return ret
	}
	// checkAttributes warns when an attribute in list is missing on every
	// command in cmds. Generic rules (empty command list) are checked against
	// all commands.
	checkAttributes := func(rule string, cmds []*commandsxmlCommand, list string) {
		if cmds == nil {
			for i := range c.Commands {
				cmds = append(cmds, &c.Commands[i])
			}
		}
		for _, attr := range strings.Fields(list) {
			found := false
			for _, cmd := range cmds {
				if cmd.hasAttribute(attr) {
					found = true
					break
				}
			}
			if !found {
				warn("lsp %s rule: attribute %q matches nothing", rule, attr)
			}
		}
	}
	symbolRule := func(kind string, rules []lspSymbolXML) {
		for _, r := range rules {
			var cmds []*commandsxmlCommand
			if r.Command != "" {
				cmds = checkCommands(kind, r.Command)
			}
			checkAttributes(kind, cmds, r.Attribute)
		}
	}
	symbolRule("defines", c.LspAnnotations.Defines)
	symbolRule("references", c.LspAnnotations.References)
	for _, r := range c.LspAnnotations.Formats {
		checkCommands("format", r.Command)
	}
	for _, r := range c.LspAnnotations.Symbols {
		checkCommands("symbol", r.Command)
	}
	for _, r := range c.LspAnnotations.Exclusives {
		cmds := checkCommands("exclusive", r.Command)
		checkAttributes("exclusive", cmds, r.Attributes)
	}
	for _, r := range c.LspAnnotations.Whens {
		cmds := checkCommands("when", r.Command)
		checkAttributes("when", cmds, r.Attribute)
		checkAttributes("when", cmds, r.Requires)
		checkAttributes("when", cmds, r.Forbids)
	}
}

// lspElement returns an empty start element in the lsp namespace with the
// given attributes. Attributes with an empty value are omitted.
func lspElement(name string, attrs ...xml.Attr) xml.StartElement {
	elt := xml.StartElement{Name: xml.Name{Local: "lsp:" + name}}
	for _, a := range attrs {
		if a.Value != "" {
			elt.Attr = append(elt.Attr, a)
		}
	}
	return elt
}

func attr(name, value string) xml.Attr {
	return xml.Attr{Name: xml.Name{Local: name}, Value: value}
}

func encodeEmpty(enc *xml.Encoder, elt xml.StartElement) {
	enc.EncodeToken(elt)
	enc.EncodeToken(elt.End())
}

// writeLspGrammarAnnotations writes the annotations that belong to the
// grammar as a whole: namespace prefixes for completion and the built-in
// symbol names.
func writeLspGrammarAnnotations(c *commandsXML, enc *xml.Encoder) {
	for _, ns := range c.LspAnnotations.Namespaces {
		encodeEmpty(enc, lspElement("namespace", attr("prefix", ns.Prefix), attr("uri", ns.URI)))
	}
	for _, b := range c.LspAnnotations.Builtins {
		encodeEmpty(enc, lspElement("builtin", attr("symbol", b.Symbol), attr("names", strings.Join(strings.Fields(b.Names), " "))))
	}
}

// writeLspElementAnnotations writes the annotations of a command directly
// after its a:documentation element.
func writeLspElementAnnotations(c *commandsXML, enc *xml.Encoder, command string) {
	encodeEmpty(enc, lspElement("doc", attr("href", lspDocURL(command))))
	if f := c.lspFormatAnnotation(command); f != nil {
		yes := func(v string) string {
			if v == "yes" {
				return "true"
			}
			return ""
		}
		encodeEmpty(enc, lspElement("format", attr("preserve", yes(f.Preserve)), attr("blank-lines", yes(f.BlankLines)), attr("inline", yes(f.Inline))))
	}
	if s := c.lspDocSymbolAnnotation(command); s != nil {
		encodeEmpty(enc, lspElement("symbol", attr("kind", s.Kind), attr("label", s.Label), attr("detail", s.Detail)))
	}
	for _, e := range c.lspExclusiveAnnotations(command) {
		content := ""
		if e.Content == "yes" {
			content = "true"
		}
		encodeEmpty(enc, lspElement("exclusive", attr("attributes", e.Attributes), attr("content", content)))
	}
	for _, w := range c.lspWhenAnnotations(command) {
		// A rule with several values expands to one annotation per value. A
		// rule without a value applies when the attribute is absent.
		values := strings.Fields(w.Value)
		if len(values) == 0 {
			values = []string{""}
		}
		for _, v := range values {
			encodeEmpty(enc, lspElement("when", attr("attribute", w.Attribute), attr("value", v), attr("requires", w.Requires), attr("forbids", w.Forbids)))
		}
	}
}

// writeLspAttributeAnnotation writes the defines/references annotation of an
// attribute directly after its a:documentation element.
func writeLspAttributeAnnotation(c *commandsXML, enc *xml.Encoder, command, attribute string) {
	if kind, symbol, form := c.lspSymbolAnnotation(command, attribute); kind != "" {
		encodeEmpty(enc, lspElement(kind, attr("symbol", symbol), attr("form", form)))
	}
}
