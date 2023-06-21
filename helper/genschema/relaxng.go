package genschema

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

const (
	// RELAXNG is the Relax NG Namespace
	RELAXNG string = "http://relaxng.org/ns/structure/1.0"
)

var (
	refElement      xml.StartElement
	emptyElement    xml.StartElement
	valueElement    xml.StartElement
	optionalElement xml.StartElement
	choiceElement   xml.StartElement
)

func init() {
	refElement = xml.StartElement{Name: xml.Name{Local: "ref"}}
	emptyElement = xml.StartElement{Name: xml.Name{Local: "empty"}}
	valueElement = xml.StartElement{Name: xml.Name{Local: "value"}}
	optionalElement = xml.StartElement{Name: xml.Name{Local: "optional"}}
	choiceElement = xml.StartElement{Name: xml.Name{Local: "choice"}}
}

// writeChildElements writes the child elements from this command to the encoder.
func writeChildElements(commands *commandsXML, enc *xml.Encoder, children []byte, lang string) {
	if len(children) == 0 {
		enc.EncodeToken(emptyElement.Copy())
		enc.EncodeToken(emptyElement.End())
		return
	}
	buf := bytes.NewBuffer(children)
	dec := xml.NewDecoder(buf)
	for {
		tok, err := dec.Token()
		if err != nil {
			return
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "cmd":
				ref := refElement.Copy()
				for _, attr := range v.Attr {
					if attr.Name.Local == "name" {
						ref.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "e_" + attr.Value}}
					}
				}
				enc.EncodeToken(ref)
			case "description":
			case "choice":
				enc.EncodeToken(choiceElement.Copy())
				for _, attribute := range v.Attr {
					if attribute.Name.Local == lang {
						enc.EncodeToken(valueElement.Copy())
						enc.EncodeToken(xml.CharData(attribute.Value))
						enc.EncodeToken(valueElement.End())
					}
				}
			case "reference":
				for _, attr := range v.Attr {
					if attr.Name.Local == "name" {
						if attr.Value == "html" {
							htmlSe := xml.StartElement{Name: xml.Name{Local: "ref"}}
							htmlSe.Attr = append(htmlSe.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: "html"})
							enc.EncodeToken(htmlSe)
							enc.EncodeToken(htmlSe.End())
						} else {
							writeChildElements(commands, enc, commands.getDefine(attr.Value), lang)
						}
					}
				}
			default:
				enc.EncodeToken(v.Copy())
			}

		case xml.EndElement:
			switch v.Name.Local {
			case "cmd":
				enc.EncodeToken(refElement.End())
			case "choice":
				enc.EncodeToken(choiceElement.End())
			default:
				enc.EncodeToken(v)
			}
		}
	}
}

func genRelaxNGSchema(commands *commandsXML, lang string, allowForeignNodes bool) ([]byte, error) {
	var outbuf bytes.Buffer
	var interleave, group xml.StartElement

	enc := xml.NewEncoder(&outbuf)
	enc.Indent("", "   ")

	grammar := xml.StartElement{Name: xml.Name{Local: "grammar", Space: RELAXNG}}
	grammar.Attr = []xml.Attr{
		{Name: xml.Name{Local: "xmlns:a"}, Value: "http://relaxng.org/ns/compatibility/annotations/1.0"},
		{Name: xml.Name{Local: "xmlns:sch"}, Value: "http://purl.oclc.org/dsdl/schematron"},
		{Name: xml.Name{Local: "ns"}, Value: SDNAMESPACE},
		{Name: xml.Name{Local: "datatypeLibrary"}, Value: "http://www.w3.org/2001/XMLSchema-datatypes"},
	}

	enc.EncodeToken(xml.Comment("Do not edit this file. Auto generated from commands.xml with xtsphelper."))
	enc.EncodeToken(xml.CharData("\n"))
	enc.EncodeToken(grammar)
	sch := xml.StartElement{Name: xml.Name{Local: "sch:ns"}}
	sch.Attr = []xml.Attr{
		{Name: xml.Name{Local: "prefix"}, Value: "t"},
		{Name: xml.Name{Local: "uri"}, Value: SDNAMESPACE},
	}
	enc.EncodeToken(sch)
	enc.EncodeToken(sch.End())

	start := xml.StartElement{Name: xml.Name{Local: "start"}}
	enc.EncodeToken(start)

	choice := xml.StartElement{Name: xml.Name{Local: "choice"}}
	enc.EncodeToken(choice)

	refLayout := xml.StartElement{Name: xml.Name{Local: "ref"}}
	refLayout.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "e_Layout"}}
	// refInclude := xml.StartElement{Name: xml.Name{Local: "ref"}}
	// refInclude.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "e_Include"}}

	enc.EncodeToken(refLayout)
	enc.EncodeToken(refLayout.End())
	// enc.EncodeToken(refInclude)
	// enc.EncodeToken(refInclude.End())
	enc.EncodeToken(choice.End())
	enc.EncodeToken(start.End())

	attributeElement := xml.StartElement{Name: xml.Name{Local: "attribute"}}

	for _, cmd := range commands.Commands {
		enc.Flush()
		for _, r := range cmd.Rules {
			if r.Lang == lang {
				outbuf.WriteString(r.Rules)
			}
		}
		def := xml.StartElement{Name: xml.Name{Local: "define"}}
		def.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "e_" + cmd.Name}}
		enc.EncodeToken(def)

		elt := xml.StartElement{Name: xml.Name{Local: "element"}}
		elt.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: cmd.Name}}
		enc.EncodeToken(elt)

		doc := xml.StartElement{Name: xml.Name{Local: "a:documentation"}}
		enc.EncodeToken(doc)
		enc.EncodeToken(xml.CharData(cmd.getCommandDescription(lang)))
		enc.EncodeToken(doc.End())

		// if the child elements contents is "empty", there is no need for allowing foreign nodes (1/2)
		if cmd.Name != "Include" && len(cmd.Childelements.Text) > 0 {
			interleave = xml.StartElement{Name: xml.Name{Local: "interleave"}}
			enc.EncodeToken(interleave)

			group = xml.StartElement{Name: xml.Name{Local: "group"}}
			enc.EncodeToken(group)
		}

		for _, attr := range cmd.Attributes {
			if attr.Optional == "yes" {
				enc.EncodeToken(optionalElement.Copy())
			}

			attelt := attributeElement.Copy()
			attelt.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: attr.Name}}
			enc.EncodeToken(attelt)

			doc := xml.StartElement{Name: xml.Name{Local: "a:documentation"}}
			enc.EncodeToken(doc)
			enc.EncodeToken(xml.CharData(attr.GetDescription(lang)))
			enc.EncodeToken(doc.End())

			if len(attr.Choice) > 0 {
				enc.EncodeToken(choiceElement.Copy())
				for _, choice := range attr.Choice {
					enc.EncodeToken(valueElement.Copy())
					enc.EncodeToken(xml.CharData(choice.Name))
					enc.EncodeToken(valueElement.End())

					doc := xml.StartElement{Name: xml.Name{Local: "a:documentation"}}
					enc.EncodeToken(doc)
					enc.EncodeToken(xml.CharData(choice.GetDescription(lang)))
					enc.EncodeToken(doc.End())

				}
				if attr.AllowXPath == "yes" {
					data := xml.StartElement{Name: xml.Name{Local: "data"}}
					data.Attr = []xml.Attr{{Name: xml.Name{Local: "type"}, Value: "string"}}
					enc.EncodeToken(data)
					param := xml.StartElement{Name: xml.Name{Local: "param"}}
					param.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "pattern"}}
					enc.EncodeToken(param)
					enc.EncodeToken(xml.CharData(`\{.+\}`))
					enc.EncodeToken(param.End())
					enc.EncodeToken(data.End())
				}

				enc.EncodeToken(choiceElement.End())
			} else if attr.Type == "yesnonumber" {
				data := xml.StartElement{Name: xml.Name{Local: "data"}}
				data.Attr = []xml.Attr{{Name: xml.Name{Local: "type"}, Value: "string"}}
				enc.EncodeToken(data)
				param := xml.StartElement{Name: xml.Name{Local: "param"}}
				param.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "pattern"}}
				enc.EncodeToken(param)
				enc.EncodeToken(xml.CharData(`[0-9]+|yes|no`))
				enc.EncodeToken(param.End())
				enc.EncodeToken(data.End())
			} else if attr.Type == "boolean" {
				enc.EncodeToken(choiceElement.Copy())
				enc.EncodeToken(valueElement.Copy())
				enc.EncodeToken(xml.CharData("yes"))
				enc.EncodeToken(valueElement.End())
				enc.EncodeToken(valueElement.Copy())
				enc.EncodeToken(xml.CharData("no"))
				enc.EncodeToken(valueElement.End())
				enc.EncodeToken(choiceElement.End())
			}

			if attr.Reference.Name != "" {
				d := commands.DefineAttrs
				for _, attrdefinition := range d {
					if attr.Reference.Name == attrdefinition.Name {
						enc.EncodeToken(choiceElement.Copy())
						if attr.AllowXPath == "yes" {
							data := xml.StartElement{Name: xml.Name{Local: "data"}}
							data.Attr = []xml.Attr{{Name: xml.Name{Local: "type"}, Value: "string"}}
							enc.EncodeToken(data)
							param := xml.StartElement{Name: xml.Name{Local: "param"}}
							param.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "pattern"}}
							enc.EncodeToken(param)
							enc.EncodeToken(xml.CharData(`\{.+\}`))
							enc.EncodeToken(param.End())
							enc.EncodeToken(data.End())
						}
						if attrdefinition.Name == "languages" {
							dl := commands.DefineList
							for _, deflist := range dl {
								if deflist.Name == "languagesshortcodes" {
									data := xml.StartElement{Name: xml.Name{Local: "data"}}
									data.Attr = []xml.Attr{{Name: xml.Name{Local: "type"}, Value: "string"}}
									enc.EncodeToken(data)
									param := xml.StartElement{Name: xml.Name{Local: "param"}}
									param.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "pattern"}}
									enc.EncodeToken(param)
									enc.EncodeToken(xml.CharData(string(deflist.Text)))
									enc.EncodeToken(param.End())
									enc.EncodeToken(data.End())
								}
							}
						}
						for _, choice := range attrdefinition.Choices {
							enc.EncodeToken(valueElement.Copy())
							enc.EncodeToken(xml.CharData(choice.Name))
							enc.EncodeToken(valueElement.End())

							doc := xml.StartElement{Name: xml.Name{Local: "a:documentation"}}
							enc.EncodeToken(doc)
							enc.EncodeToken(xml.CharData(choice.GetDescription(lang)))
							enc.EncodeToken(doc.End())

						}
						enc.EncodeToken(choiceElement.End())
					}
				}
			}
			enc.EncodeToken(attelt.End())
			if attr.Optional == "yes" {
				enc.EncodeToken(optionalElement.Copy().End())
			}
		}
		writeChildElements(commands, enc, cmd.Childelements.Text, lang)

		// if the child elements contents is "empty", there is no need for allowing foreign nodes (2/2)
		if cmd.Name != "Include" && len(cmd.Childelements.Text) > 0 {
			enc.EncodeToken(group.End())
			if allowForeignNodes {
				ref := xml.StartElement{Name: xml.Name{Local: "ref"}}
				ref.Attr = []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "foreign-nodes"}}
				enc.EncodeToken(ref)
				enc.EncodeToken(ref.End())
			}
			enc.EncodeToken(interleave.End())
		}
		enc.EncodeToken(elt.End())
		enc.EncodeToken(def.End())
	}
	enc.Flush()
	fmt.Fprint(&outbuf, `
	<!-- allow HTML in <Value> ... </Value> -->
	<define name="html">
		<zeroOrMore>
		   <choice>
			  <element name="a">
				 <attribute name="href"/>
				 <ref name="html"/>
			  </element>
			  <element name="b"><ref name="html" /></element>
			  <element name="br"><empty /></element>
			  <element name="code"><ref name="html" /></element>
			  <element name="i"><ref name="html" /></element>
			  <element name="kbd"><ref name="html" /></element>
			  <element name="li"><ref name="html" /></element>
			  <element name="p"><ref name="html" /></element>
			  <element name="span"><ref name="html" /><oneOrMore><attribute><anyName/></attribute></oneOrMore></element>
			  <element name="u"><ref name="html" /></element>
			  <element name="ul"><ref name="html" /></element>
			  <text></text>
		   </choice>
		</zeroOrMore>
	</define>

`)
	if allowForeignNodes {
		enc.Flush()
		// See feature request #144
		fmt.Fprintln(&outbuf, fmt.Sprintf(`
	<!-- This pattern allows any element from any namespace -->
	<define name="anything">
      <zeroOrMore>
         <choice>
            <element>
               <anyName/>
               <ref name="anything"/>
            </element>
            <attribute>
               <anyName/>
            </attribute>
            <text/>
         </choice>
      </zeroOrMore>
   </define>
   <define name="foreign-elements">
      <zeroOrMore>
         <element>
            <anyName>
               <except>
                  <nsName ns=""/>
                  <nsName ns="%s"/>
                  <nsName ns="%s"/>
               </except>
            </anyName>
            <ref name="anything"/>
         </element>
      </zeroOrMore>
   </define>
   <define name="foreign-attributes">
      <zeroOrMore>
         <attribute>
            <anyName>
               <except>
                  <nsName ns=""/>
                  <nsName ns="%s"/>
                  <nsName ns="%s"/>
               </except>
            </anyName>
         </attribute>
      </zeroOrMore>
   </define>
   <define name="foreign-nodes">
      <zeroOrMore>
         <choice>
            <ref name="foreign-attributes"/>
            <ref name="foreign-elements"/>
         </choice>
      </zeroOrMore>
   </define>`, SDNAMESPACE, FNNAMESPACE, SDNAMESPACE, FNNAMESPACE))
	}

	enc.EncodeToken(grammar.End())
	enc.EncodeToken(xml.CharData("\n"))
	enc.Flush()
	return outbuf.Bytes(), nil
}
