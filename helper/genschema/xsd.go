package genschema

import (
	"bytes"
	"encoding/xml"
)

// XSDNS is the W3C XML Schema namespace.
const XSDNS string = "http://www.w3.org/2001/XMLSchema"

// childModel is the flattened content model of a command. The XSD output drops
// the RNG ordering/occurrence operators (interleave, group, oneOrMore, ...) and
// keeps only the set of permitted child elements plus the mixed/wildcard flags.
// That is a deliberate, pragmatic simplification: XSD cannot express RNG's
// interleave, and a flat <xs:choice minOccurs="0" maxOccurs="unbounded"> is the
// most permissive content model that is still useful for editor autocompletion
// and never produces Unique Particle Attribution conflicts.
type childModel struct {
	elements []string // permitted child element names, first-seen order, deduped
	mixed    bool     // character data allowed (RNG <text> or free content)
	wildcard bool     // free foreign/HTML content allowed (RNG <anycontent> or html)
}

// collectChildModel scans the raw inner XML of a <childelements> block (or a
// <define> body) and flattens it into cm. References to other defines are
// expanded recursively; the html reference and <anycontent> become a wildcard.
// Comments (e.g. disabled <cmd> entries) are skipped because only StartElement
// tokens are inspected.
func (commands *commandsXML) collectChildModel(raw []byte, cm *childModel, seenElt, seenDef map[string]bool) {
	dec := xml.NewDecoder(bytes.NewReader(raw))
	for {
		tok, err := dec.Token()
		if err != nil {
			return
		}
		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		switch se.Name.Local {
		case "cmd":
			for _, a := range se.Attr {
				if a.Name.Local == "name" && !seenElt[a.Value] {
					seenElt[a.Value] = true
					cm.elements = append(cm.elements, a.Value)
				}
			}
		case "reference":
			for _, a := range se.Attr {
				if a.Name.Local != "name" {
					continue
				}
				if a.Value == "html" {
					cm.wildcard = true
					cm.mixed = true
				} else if !seenDef[a.Value] {
					seenDef[a.Value] = true
					commands.collectChildModel(commands.getDefine(a.Value), cm, seenElt, seenDef)
				}
			}
		case "anycontent":
			cm.wildcard = true
			cm.mixed = true
		case "text":
			cm.mixed = true
		}
	}
}

func xsdStart(local string, attrs ...xml.Attr) xml.StartElement {
	return xml.StartElement{Name: xml.Name{Local: local}, Attr: attrs}
}

func xsdAttr(name, value string) xml.Attr {
	return xml.Attr{Name: xml.Name{Local: name}, Value: value}
}

// encEmpty writes a self-closing element (start + end with no children).
func encEmpty(enc *xml.Encoder, local string, attrs ...xml.Attr) {
	el := xsdStart(local, attrs...)
	enc.EncodeToken(el)
	enc.EncodeToken(el.End())
}

// encDoc writes an <xs:annotation><xs:documentation> block unless text is empty.
func encDoc(enc *xml.Encoder, text string) {
	if text == "" {
		return
	}
	ann := xsdStart("xs:annotation")
	doc := xsdStart("xs:documentation")
	enc.EncodeToken(ann)
	enc.EncodeToken(doc)
	enc.EncodeToken(xml.CharData(text))
	enc.EncodeToken(doc.End())
	enc.EncodeToken(ann.End())
}

// encSimpleEnum writes a string-based simpleType restricted to an enumeration.
func encSimpleEnum(enc *xml.Encoder, values []string) {
	st := xsdStart("xs:simpleType")
	rest := xsdStart("xs:restriction", xsdAttr("base", "xs:string"))
	enc.EncodeToken(st)
	enc.EncodeToken(rest)
	for _, v := range values {
		encEmpty(enc, "xs:enumeration", xsdAttr("value", v))
	}
	enc.EncodeToken(rest.End())
	enc.EncodeToken(st.End())
}

// encSimplePattern writes a string-based simpleType restricted to a pattern.
func encSimplePattern(enc *xml.Encoder, pattern string) {
	st := xsdStart("xs:simpleType")
	rest := xsdStart("xs:restriction", xsdAttr("base", "xs:string"))
	enc.EncodeToken(st)
	enc.EncodeToken(rest)
	encEmpty(enc, "xs:pattern", xsdAttr("value", pattern))
	enc.EncodeToken(rest.End())
	enc.EncodeToken(st.End())
}

// encAttribute writes a single <xs:attribute>. When the attribute also accepts
// an XPath expression in {...}, enumerations/patterns are dropped in favour of a
// plain string so that both literal values and {xpath} validate (pragmatic).
func encAttribute(enc *xml.Encoder, commands *commandsXML, attr *commandsxmlAttribute, lang string) {
	att := xsdStart("xs:attribute", xsdAttr("name", attr.Name))
	if attr.Optional == "no" {
		att.Attr = append(att.Attr, xsdAttr("use", "required"))
	}
	enc.EncodeToken(att)
	encDoc(enc, attr.GetDescription(lang))

	allowXPath := attr.AllowXPath == "yes"
	switch {
	case len(attr.Choice) > 0:
		if !allowXPath {
			vals := make([]string, 0, len(attr.Choice))
			for _, c := range attr.Choice {
				vals = append(vals, c.Name)
			}
			encSimpleEnum(enc, vals)
		}
	case attr.Type == "boolean":
		encSimpleEnum(enc, []string{"yes", "no"})
	case attr.Type == "yesnonumber":
		encSimplePattern(enc, `[0-9]+|yes|no`)
	case attr.Type == "pdfformat":
		// Kept in sync with document.ParseFormat and the RNG generator.
		encSimplePattern(enc, `\s*(PDF|PDF/A-3b|PDF/X-3|PDF/X-4|PDF/UA|PDF/UA-1|PDF/UA-2)(\s*,\s*(PDF|PDF/A-3b|PDF/X-3|PDF/X-4|PDF/UA|PDF/UA-1|PDF/UA-2))*\s*`)
	case attr.Reference.Name != "" && !allowXPath:
		for _, d := range commands.DefineAttrs {
			if d.Name == attr.Reference.Name && len(d.Choices) > 0 {
				vals := make([]string, 0, len(d.Choices))
				for _, c := range d.Choices {
					vals = append(vals, c.Name)
				}
				encSimpleEnum(enc, vals)
			}
		}
	}
	enc.EncodeToken(att.End())
}

// genXSDSchema generates an XSD layout schema for the given language directly
// from commands.xml, without an external converter. The element/attribute names
// and the target namespace are identical for every language; only the embedded
// documentation differs. See childModel for the content-model simplification.
func genXSDSchema(commands *commandsXML, lang string) ([]byte, error) {
	var outbuf bytes.Buffer
	enc := xml.NewEncoder(&outbuf)
	enc.Indent("", "  ")

	enc.EncodeToken(xml.Comment("Do not edit this file. Auto generated from commands.xml with xtshelper."))
	enc.EncodeToken(xml.CharData("\n"))

	schema := xsdStart("xs:schema",
		xsdAttr("xmlns:xs", XSDNS),
		xsdAttr("xmlns:en", SDNAMESPACE),
		xsdAttr("targetNamespace", SDNAMESPACE),
		xsdAttr("elementFormDefault", "qualified"),
	)
	enc.EncodeToken(schema)

	for _, cmd := range commands.Commands {
		cm := &childModel{}
		commands.collectChildModel(cmd.Childelements.Text, cm, map[string]bool{}, map[string]bool{})

		elt := xsdStart("xs:element", xsdAttr("name", cmd.Name))
		enc.EncodeToken(elt)
		encDoc(enc, cmd.getCommandDescription(lang))

		ct := xsdStart("xs:complexType")
		if cm.mixed {
			ct.Attr = append(ct.Attr, xsdAttr("mixed", "true"))
		}
		enc.EncodeToken(ct)

		if cm.wildcard {
			// Free content (e.g. <Value>, <HTML>): a single lax wildcard. Mixing
			// explicit element refs with a same-namespace wildcard would violate
			// XSD's Unique Particle Attribution rule, so the refs are dropped here.
			seq := xsdStart("xs:sequence")
			enc.EncodeToken(seq)
			encEmpty(enc, "xs:any",
				xsdAttr("processContents", "lax"),
				xsdAttr("minOccurs", "0"),
				xsdAttr("maxOccurs", "unbounded"))
			enc.EncodeToken(seq.End())
		} else if len(cm.elements) > 0 {
			choice := xsdStart("xs:choice", xsdAttr("minOccurs", "0"), xsdAttr("maxOccurs", "unbounded"))
			enc.EncodeToken(choice)
			for _, name := range cm.elements {
				encEmpty(enc, "xs:element", xsdAttr("ref", "en:"+name))
			}
			enc.EncodeToken(choice.End())
		}

		for i := range cmd.Attributes {
			encAttribute(enc, commands, &cmd.Attributes[i], lang)
		}

		enc.EncodeToken(ct.End())
		enc.EncodeToken(elt.End())
	}

	enc.EncodeToken(schema.End())
	enc.EncodeToken(xml.CharData("\n"))
	if err := enc.Flush(); err != nil {
		return nil, err
	}
	return outbuf.Bytes(), nil
}
