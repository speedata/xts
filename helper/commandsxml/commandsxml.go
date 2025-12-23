// Package commandsxml reads the commands.xml file.
package commandsxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	multipleSpace       *regexp.Regexp
	everysecondast      *regexp.Regexp
	everysecondbacktick *regexp.Regexp
)

func init() {
	multipleSpace = regexp.MustCompile(`\s+`)
	everysecondast = regexp.MustCompile(`(?s)(.*?)\*(.*?)\*`)
	everysecondbacktick = regexp.MustCompile("(?s)(.*?\\S)`(\\*)`")
}

type para struct {
	commands *Commands
	Text     []byte `xml:",innerxml"`
}

func (p *para) Markdown() string {
	ret := []string{}
	c := p.commands
	r := bytes.NewReader(p.Text)
	dec := xml.NewDecoder(r)

	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "cmd":
				var x *Command
				var cmdname string
				for _, attribute := range v.Attr {
					if attribute.Name.Local == "name" {
						x = c.commandsEn[attribute.Value]
						if x == nil {
							fmt.Printf("There is an unknown cmd in the para section of %q\n", attribute.Value)
							os.Exit(-1)
						}
						cmdname = x.Name
					}
				}
				ret = append(ret, fmt.Sprintf(`[%s](%s)`, cmdname, x.CmdLink()))
			case "tt":
				ret = append(ret, "`")
			}
		case xml.CharData:
			ret = append(ret, string(v.Copy()))
		case xml.EndElement:
			switch v.Name.Local {
			case "tt":
				ret = append(ret, "`")
			}
		}
	}
	ret = append(ret, "\n\n")
	a := strings.Join(ret, "")
	a = everysecondast.ReplaceAllString(a, "$1\\*$2*")
	a = everysecondbacktick.ReplaceAllString(a, "$1`$2`")
	a = strings.Replace(a, "&", "\\&", -1)
	return a
}

func (p *para) String() string {
	ret := []string{}
	r := bytes.NewReader(p.Text)
	dec := xml.NewDecoder(r)
outer:
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "cmd":
				var cmdname string
				for _, attribute := range v.Attr {
					if attribute.Name.Local == "name" {
						cmdname = attribute.Value
					}
				}
				ret = append(ret, cmdname)
			}
		case xml.CharData:
			ret = append(ret, string(v.Copy()))
		case xml.EndElement:
			switch v.Name.Local {
			case "description":
				break outer
			}
		}
	}
	return multipleSpace.ReplaceAllString(strings.Join(ret, ""), " ")
}

type define struct {
	Name string `xml:"name,attr"`
	Text []byte `xml:",innerxml"`
}

// DescriptionText returns the description of the attribute without markup.
func (c *Choice) DescriptionText() string {
	return descriptiontext(c.commands, c.descriptionEn.Text)
}

// Choice represents alternative attribute values
type Choice struct {
	commands      *Commands
	Text          []byte `xml:",innerxml"`
	Name          string `xml:"en,attr"`
	Pro           bool
	descriptionEn *description
	descriptionDe *description
}

// Attribute has all information about each attribute
type Attribute struct {
	Choice        []*Choice
	Name          string
	CSS           string
	Since         string
	Type          string
	Optional      bool
	AllowXPath    bool
	Pro           bool
	commands      *Commands
	command       *Command
	descriptionEn *description
	descriptionDe *description
}

// UnmarshalXML fills the Choice value
func (c *Choice) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "description":
				d := &description{}
				d.commands = c.commands
				dec.DecodeElement(d, &v)
				switch d.Lang {
				case "en":
					c.descriptionEn = d
				case "de":
					c.descriptionDe = d
				}
			}
		}
	}
}

// Attlink returns a string in the form of cmd-commandname-attribute e.g. cmd-setvariable-select
func (a *Attribute) Attlink() string {
	cmd := a.command
	ret := []string{}
	ret = append(ret, cmd.CmdLink())
	tmp := strings.ToLower(a.Name)
	tmp = strings.Replace(tmp, ":", "_", -1)
	ret = append(ret, tmp)
	return strings.Join(ret, "_")
}

// UnmarshalXML fills the attribute from the given XML segment
func (a *Attribute) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "description":
				d := &description{}
				d.commands = a.commands
				dec.DecodeElement(d, &v)
				switch d.Lang {
				case "en":
					a.descriptionEn = d
				case "de":
					a.descriptionDe = d
				}
			case "choice":
				c := &Choice{}
				c.commands = a.commands
				dec.DecodeElement(c, &v)
				for _, attribute := range v.Attr {
					switch attribute.Name.Local {
					case "en":
						c.Name = attribute.Value
					case "pro":
						c.Pro = attribute.Value == "yes"
					}
				}
				a.Choice = append(a.Choice, c)
			}
		}
	}
}

// DescriptionText returns the description of the attribute without markup.
func (a *Attribute) DescriptionText() string {
	return descriptiontext(a.commands, a.descriptionEn.Text)
}

// DescriptionMarkdown returns the description of the attribute as an asciidoctor blob.
func (a *Attribute) DescriptionMarkdown() string {
	var ret []string
	ret = append(ret, a.descriptionEn.Markdown())
	var name string
	var desc string
	for _, c := range a.Choice {
		name = c.Name
		desc = c.descriptionEn.Markdown()
		if c.Pro {
			ret = append(ret, "[.profeature]")
		}
		ret = append(ret, "\n    `"+name+"`\n    :    "+desc)
	}
	return string(strings.Join(ret, "\n"))
}

// Childelement has all child elements of a command.
type Childelement struct {
	commands *Commands
	Text     []byte `xml:",innerxml"`
}

type example struct {
	commands *Commands
	Lang     string `xml:"http://www.w3.org/XML/1998/namespace lang,attr"`
	Text     []byte `xml:",innerxml"`
}

type seealso struct {
	commands *Commands
	Text     []byte `xml:",innerxml"`
}

type description struct {
	commands *Commands
	Lang     string `xml:"http://www.w3.org/XML/1998/namespace lang,attr"`
	Text     []byte `xml:",innerxml"`
}

// Markdown returns the description in markdown format.
func (d *description) Markdown() string {
	if d == nil {
		return ""
	}
	r := bytes.NewReader(d.Text)
	dec := xml.NewDecoder(r)
	var ret []string
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "para":
				p := &para{}
				p.commands = d.commands
				err = dec.DecodeElement(p, &v)
				if err != nil {
					panic(err)
				}
				ret = append(ret, p.Markdown())

			}
		}
	}
	return strings.Join(ret, "")
}

func (d *description) String() string {
	if d == nil {
		return ""
	}
	r := bytes.NewReader(d.Text)
	dec := xml.NewDecoder(r)
	var ret []string
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "para":
				p := &para{}
				p.commands = d.commands
				err = dec.DecodeElement(p, &v)
				if err != nil {
					panic(err)
				}
				ret = append(ret, p.String())

			}
		}
	}
	return strings.Join(ret, "")
}

// SchematronRules represents schematron rules
type SchematronRules struct {
	Lang  string `xml:"lang,attr"`
	Rules string `xml:",innerxml"`
}

// Command has information about a command
type Command struct {
	Attr           []*Attribute
	Name           string
	CSS            string
	Since          string
	Pro            bool
	Deprecated     bool
	Rules          []SchematronRules `xml:"rules"`
	childelement   *Childelement
	remarkEn       *description
	remarkDe       *description
	infoEn         *description
	infoDe         *description
	descriptionEn  *description
	descriptionDe  *description
	commands       *Commands
	parentelements map[*Command]bool
	examplesEn     []*example
	examplesDe     []*example
	children       map[string][]*Command
	seealso        *seealso
}

// Parents returns all parent commands
func (c *Command) Parents() []*Command {

	var cmds []*Command
	for k := range c.parentelements {
		cmds = append(cmds, k)
	}
	sort.Sort(commandsbyen{cmds})

	return cmds
}

func (c *Command) String() string {
	return c.Name
}

// UnmarshalXML fills the command from the XML segment
func (c *Command) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "description":
				d := &description{}
				d.commands = c.commands
				dec.DecodeElement(d, &v)
				switch d.Lang {
				case "en":
					c.descriptionEn = d
				case "de":
					c.descriptionDe = d

				}

			case "info":
				d := &description{}
				d.commands = c.commands
				dec.DecodeElement(d, &v)
				switch d.Lang {
				case "en":
					c.infoEn = d
				case "de":
					c.infoDe = d
				}
			case "remark":
				d := &description{}
				d.commands = c.commands
				dec.DecodeElement(d, &v)
				switch d.Lang {
				case "en":
					c.remarkEn = d
				case "de":
					c.remarkDe = d
				}

			case "attribute":
				a := &Attribute{}
				a.commands = c.commands
				a.command = c
				dec.DecodeElement(a, &v)
				for _, attribute := range v.Attr {
					switch attribute.Name.Local {
					case "en":
						a.Name = attribute.Value
					case "css":
						a.CSS = attribute.Value
					case "since":
						a.Since = attribute.Value
					case "optional":
						a.Optional = attribute.Value == "yes"
					case "type":
						a.Type = attribute.Value
					case "allowxpath":
						a.AllowXPath = attribute.Value == "yes"
					case "pro":
						a.Pro = attribute.Value == "yes"
					}
				}

				c.Attr = append(c.Attr, a)
			case "childelements":
				child := &Childelement{}
				child.commands = c.commands
				dec.DecodeElement(child, &v)
				c.childelement = child
			case "example":
				e := &example{}
				e.commands = c.commands
				dec.DecodeElement(e, &v)
				switch e.Lang {
				case "en":
					c.examplesEn = append(c.examplesEn, e)
				case "de":
					c.examplesDe = append(c.examplesDe, e)
				}
			case "seealso":
				e := &seealso{}
				e.commands = c.commands
				dec.DecodeElement(e, &v)
				c.seealso = e
			case "rules":
				r := &SchematronRules{}
				dec.DecodeElement(r, &v)
				c.Rules = append(c.Rules, *r)
			}
		}
	}
}

// MDlink returns the command name with ".md"
func (c *Command) MDlink() string {
	if c == nil {
		return ""
	}
	tmp := url.URL{Path: strings.ToLower(c.Name)}
	filenameSansExtension := tmp.String()
	filenameSansExtension = strings.Replace(filenameSansExtension, "-", "_", -1)
	return filenameSansExtension + ".md"
}

// Htmllink returns a text such as "mycmd.html"
func (c *Command) Htmllink() string {
	if c == nil {
		return ""
	}
	tmp := url.URL{Path: strings.ToLower(c.Name)}
	filenameSansExtension := tmp.String()
	return filenameSansExtension + ".html"
}

// CmdLink returns a text such as cmd-atpageshipout
func (c *Command) CmdLink() string {
	if c == nil {
		return ""
	}
	tmp := url.URL{Path: strings.ToLower(c.Name)}
	filenameSansExtension := tmp.String()
	filenameSansExtension = strings.Replace(filenameSansExtension, "-", "_", -1)
	return "../" + filenameSansExtension
}

// DescriptionMarkdown returns the description of the command as a asciidoctor
// blob.
func (c *Command) DescriptionMarkdown() string {
	return c.descriptionEn.Markdown()
}

// RemarkMarkdown returns the remark section as a formatted asciidoctor blob.
func (c *Command) RemarkMarkdown() string {
	lang := "en"
	var ret string
	switch lang {
	case "en":
		ret = c.remarkEn.Markdown()
	case "de":
		ret = c.remarkDe.Markdown()
	default:
		ret = ""
	}
	return ret
}

// InfoMarkdown returns the info section as a asciidoctor blob.
func (c *Command) InfoMarkdown() string {
	lang := "en"
	var r *bytes.Reader
	switch lang {
	case "en":
		if x := c.infoEn; x != nil {
			r = bytes.NewReader(x.Text)
		} else {
			return ""
		}
	case "de":
		if x := c.infoDe; x != nil {
			r = bytes.NewReader(x.Text)
		} else {
			return ""
		}
	}

	var ret []string
	dec := xml.NewDecoder(r)

	inListing := false
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "listing":
				inListing = true
			case "image":
				var fn, wd string
				for _, a := range v.Attr {
					wd = "auto"
					if a.Name.Local == "file" {
						fn = a.Value
					} else if a.Name.Local == "width" {
						wd = fmt.Sprintf(`%s`, a.Value)
					}
				}
				ret = append(ret, fmt.Sprintf("\nimage::%s[width=%s]\n", fn, wd))
			case "para":
				p := &para{}
				p.commands = c.commands
				err = dec.DecodeElement(p, &v)
				if err != nil {
					panic(err)
				}
				ret = append(ret, "\n")
				ret = append(ret, p.Markdown())
				ret = append(ret, "\n")
			}
		case xml.CharData:
			if inListing {
				ret = append(ret, "\n```xml\n")
				ret = append(ret, string(v))
				ret = append(ret, "\n```\n")

			}
		case xml.EndElement:
			switch v.Name.Local {
			case "listing":
				inListing = false
			}
		}
	}
	return strings.Join(ret, "")
}

func descriptiontext(c *Commands, text []byte) string {
	r := bytes.NewReader(text)
	dec := xml.NewDecoder(r)
	var ret []string

	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "para":
				p := &para{}
				p.commands = c
				err = dec.DecodeElement(p, &v)
				if err != nil {
					panic(err)
				}
				ret = append(ret, p.String())
			}
		}
	}
	return strings.Join(ret, "")

}

// DescriptionText returns the description as text.
func (c *Command) DescriptionText(lang string) string {
	switch lang {
	case "en":
		return descriptiontext(c.commands, c.descriptionEn.Text)
	case "de":
		return descriptiontext(c.commands, c.descriptionDe.Text)
	default:
		return ""
	}
}

type reference struct {
	longnameEn  string
	longnameDe  string
	pagename    string
	chaptername string
}

// Attributes returns all attributes for the command
func (c *Command) Attributes() []*Attribute {
	sort.Sort(attributesbyen{c.Attr})
	ret := make([]*Attribute, len(c.Attr))
	copy(ret, c.Attr)
	return ret
}

// ExampleMarkdown returns the examples section as an markdown blob.
func (c *Command) ExampleMarkdown() string {
	lang := "en"
	var r *bytes.Reader
	switch lang {
	case "en":
		if x := c.examplesEn; len(x) != 0 {
			r = bytes.NewReader(x[0].Text)
		} else {
			return ""
		}
	case "de":
		if x := c.examplesDe; len(x) != 0 {
			r = bytes.NewReader(x[0].Text)
		} else {
			return ""
		}
	default:
		return ""
	}
	var ret []string
	dec := xml.NewDecoder(r)

	inListing := false
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "listing":
				inListing = true
			case "image":
				var fn, wd string
				for _, a := range v.Attr {
					wd = "auto"
					if a.Name.Local == "file" {
						fn = a.Value
					} else if a.Name.Local == "width" {
						wd = fmt.Sprintf(`%s`, a.Value)
					}
				}
				ret = append(ret, fmt.Sprintf("\n![](../img/%s){: style=\"width=%s\"; }\n", fn, wd))
			case "para":
				p := &para{}
				p.commands = c.commands
				err = dec.DecodeElement(p, &v)
				if err != nil {
					panic(err)
				}
				ret = append(ret, "\n")
				ret = append(ret, p.Markdown())
				ret = append(ret, "\n")
			}
		case xml.CharData:
			if inListing {
				ret = append(ret, "```xml\n")
				ret = append(ret, string(v))
				ret = append(ret, "\n```\n")
			}
		case xml.EndElement:
			switch v.Name.Local {
			case "listing":
				inListing = false
			}
		}
	}
	return strings.Join(ret, "")
}

func getchildren(c *Commands, dec *xml.Decoder) []*Command {
	var cmds []*Command
	for {
		tok, err := dec.Token()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch eltname := v.Name.Local; eltname {
			case "cmd":
				var cmdname string
				for _, attribute := range v.Attr {
					if attribute.Name.Local == "name" {
						cmdname = attribute.Value
					}
				}
				cmds = append(cmds, c.commandsEn[cmdname])
			case "reference":
				var refname string
				for _, attribute := range v.Attr {
					if attribute.Name.Local == "name" {
						refname = attribute.Value
					}
				}
				if refname != "html" {
					dec = xml.NewDecoder(bytes.NewReader(c.defines[refname].Text))
					x := getchildren(c, dec)
					for _, command := range x {
						cmds = append(cmds, command)
					}
				}
			default:
			}
		}
	}
	return cmds
}

// Command returns a Command structure for the command named in commandname.
func (c *Commands) Command(commandname string) *Command {
	return c.commandsEn[commandname]
}

// Childelements returns a list of commands that are allowed within this command.
func (c *Command) Childelements() []*Command {
	if c == nil {
		return nil
	}
	x := c.children["en"]
	if x != nil {
		return x
	}

	r := bytes.NewReader(c.childelement.Text)
	dec := xml.NewDecoder(r)

	cmds := getchildren(c.commands, dec)

	for _, v := range cmds {
		v.parentelements[c] = true
	}
	sort.Sort(commandsbyen{cmds})
	c.children["en"] = cmds
	return cmds
}

// GetDefineText returns the byte value of a define section in the commands xml
func (c *Commands) GetDefineText(section string) []byte {
	if t, ok := c.defines[section]; ok {
		return t.Text
	}
	return []byte("")
}

// Commands returns a list of all commands sorted by name.
func (c *Commands) Commands() []*Command {
	return c.commandsSortedEn
}

// Commands is the root structure of all Commands
type Commands struct {
	commandsEn       map[string]*Command
	commandsSortedEn []*Command
	defines          map[string]*define
}

// sorting (de, en)
type sortcommands []*Command
type sortattributes []*Attribute

func (s sortcommands) Len() int        { return len(s) }
func (s sortcommands) Swap(i, j int)   { s[i], s[j] = s[j], s[i] }
func (s sortattributes) Len() int      { return len(s) }
func (s sortattributes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type commandsbyen struct{ sortcommands }
type attributesbyen struct{ sortattributes }

func (s commandsbyen) Less(i, j int) bool { return s.sortcommands[i].Name < s.sortcommands[j].Name }
func (s attributesbyen) Less(i, j int) bool {
	return s.sortattributes[i].Name < s.sortattributes[j].Name
}

// ReadCommandsFile reads from the reader. It must be in the format of a commands file.
func ReadCommandsFile(r io.Reader) (*Commands, error) {
	commands := &Commands{}
	commands.defines = make(map[string]*define)
	commands.commandsEn = make(map[string]*Command)
	dec := xml.NewDecoder(r)
	for {
		tok, err := dec.Token()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		switch v := tok.(type) {
		case xml.StartElement:

			switch v.Name.Local {
			// case "commands":
			// 	// OK, root element
			case "define":
				d := &define{}
				err = dec.DecodeElement(d, &v)
				if err != nil {
					return nil, err
				}
				commands.defines[d.Name] = d
			case "command":
				c := &Command{}
				c.commands = commands
				c.children = make(map[string][]*Command)
				c.parentelements = make(map[*Command]bool)

				err = dec.DecodeElement(c, &v)
				if err != nil {
					return nil, err
				}
				commands.commandsSortedEn = append(commands.commandsSortedEn, c)

				for _, attribute := range v.Attr {
					if attribute.Name.Local == "en" {
						commands.commandsEn[attribute.Value] = c
						c.Name = attribute.Value
					}
					if attribute.Name.Local == "css" {
						c.CSS = attribute.Value
					}
					if attribute.Name.Local == "since" {
						c.Since = attribute.Value
					}
					if attribute.Name.Local == "pro" {
						c.Pro = attribute.Value == "yes"
					}
					if attribute.Name.Local == "deprecated" {
						c.Deprecated = attribute.Value == "yes"
					}
				}
			}
		}
	}
	sort.Sort(commandsbyen{commands.commandsSortedEn})
	// to get the full list of parent elements, the child element of each command
	// have to be called at least once. I know this sucks...
	for _, v := range commands.commandsEn {
		v.Childelements()
	}
	return commands, nil
}

// LoadCommandsFile opens the doc/commands.xml/commands.xml in the given base dir
func LoadCommandsFile(basedir string) (*Commands, error) {
	r, err := os.Open(filepath.Join(basedir, "doc", "commands-xml", "commands.xml"))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ReadCommandsFile(r)
}
