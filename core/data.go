package core

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/speedata/bagme/document"
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/csshtml"
	"github.com/speedata/boxesandglue/frontend"
	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
)

// A seqfunc is used to defer the execution of a function.
type seqfunc func() (xpath.Sequence, error)

func returnEvalBodyLater(layoutelt *goxml.Element, xd *xtsDocument, ctx *xpath.Context) seqfunc {
	oldCtx := xpath.CopyContext(ctx)
	return func() (xpath.Sequence, error) {
		eval, err := dispatch(xd, layoutelt, xd.data)
		xd.data.Ctx = oldCtx
		return eval, err
	}
}

// applyLayoutStylesheet creates an HTML fragment, applies CSS and reads the
// attributes from the fragment. This is handy when styling layout elements with
// CSS.
func (xd *xtsDocument) applyLayoutStylesheet(classname string, id string, style string, eltnames ...string) (map[string]string, error) {
	htmlstrings := []string{}
	for i, eltname := range eltnames {
		if i == len(eltnames)-1 {
			htmlstrings = append(htmlstrings, "<", eltname)
			if classname != "" {
				htmlstrings = append(htmlstrings, fmt.Sprintf(" class=%q ", classname))
			}
			if id != "" {
				htmlstrings = append(htmlstrings, fmt.Sprintf(" id=%q", id))
			}
			if style != "" {
				htmlstrings = append(htmlstrings, fmt.Sprintf(" style=%q", style))
			}
			htmlstrings = append(htmlstrings, ">")
		} else {
			htmlstrings = append(htmlstrings, "<"+eltname+">")

		}
	}
	for i, j := 0, len(eltnames)-1; i < j; i, j = i+1, j-1 {
		eltnames[i], eltnames[j] = eltnames[j], eltnames[i]
	}

	for _, eltname := range eltnames {
		htmlstrings = append(htmlstrings, "</"+eltname+">")
	}

	htmlstring := strings.Join(htmlstrings, "")
	doc, err := xd.layoutcss.ParseHTMLFragment(htmlstring)
	if err != nil {
		return nil, err
	}
	eltname := eltnames[0]
	sel := doc.Find(eltname)
	a := sel.Nodes[0]
	attrs, _, _ := csshtml.ResolveAttributes(a.Attr)
	return attrs, nil
}

// decodeHTMLFromHTMLNode takes a parsed HTML structure and return a function
// that formats the the input to a VList.
func (xd *xtsDocument) decodeHTMLFromHTMLNode(input *html.Node) (frontend.FormatToVList, error) {
	ftv := func(wd bag.ScaledPoint) (*node.VList, error) {
		d := document.NewWithFrontend(xd.document, xd.layoutcss)
		te, err := d.ParseHTMLFromNode(input)
		if err != nil {
			return nil, err
		}
		vl, err := d.CreateVlist(te, wd)
		if err != nil {
			return nil, newTypesettingError(err)
		}
		return vl, nil
	}

	return ftv, nil
}

// decodeHTML takes a simple text and return a function that formats the the
// input to a VList.
func (xd *xtsDocument) decodeHTML(input string) (frontend.FormatToVList, error) {
	ftv := func(wd bag.ScaledPoint) (*node.VList, error) {
		d := document.NewWithFrontend(xd.document, xd.layoutcss)
		te, err := d.HTMLToText(input)
		if err != nil {
			return nil, err
		}
		vl, err := d.CreateVlist(te, wd)
		if err != nil {
			return nil, err
		}

		return vl, err
	}

	return ftv, nil
}

// parseHTMLText takes well formed XML input and interprets this as HTML.
func (xd *xtsDocument) parseHTMLText(input string) (*html.Node, error) {
	s := strings.NewReader(input)
	return html.Parse(s)
}

// getFontSizeLeading returns the font size and ghe leading size from a single
// input. The input can have the format "10pt/12pt" where the first length is a
// font size and the second length is the leading.
func (xd *xtsDocument) getFontSizeLeading(size string) (fontsize bag.ScaledPoint, leading bag.ScaledPoint, err error) {
	if sp := strings.Split(size, "/"); len(sp) == 2 {
		if fontsize, err = bag.Sp(sp[0]); err != nil {
			return
		}
		if leading, err = bag.Sp(sp[1]); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("unknown font size %s", size)
		return
	}
	return
}

func debugFrontendText(fe *frontend.Text, level int) {
	for i, itm := range fe.Items {
		switch t := itm.(type) {
		case *frontend.Text:
			debugFrontendText(t, level+1)
		case string:
			fmt.Println(strings.Repeat("-", level), i, t)
		default:
			fmt.Printf("itm %T\n", itm)
		}
	}
}

// Get the values from the child elements of B, Paragraph and its ilk and return
// the equivalent HTML structure.
func (xd *xtsDocument) getTextvalues(tagname string, seq xpath.Sequence, attributes map[string]string, cmdname string, line int) (*html.Node, error) {
	n := &html.Node{}
	n.Data = tagname
	n.Type = html.ElementNode
	for k, v := range attributes {
		n.Attr = append(n.Attr, html.Attribute{Key: k, Val: v})
	}

	if len(seq) == 0 && tagname == "p" {
		cd := &html.Node{}
		cd.Type = html.TextNode
		cd.Data = "\u200B"
		n.AppendChild(cd)
		return n, nil
	}
	for _, itm := range seq {
		switch t := itm.(type) {
		case string:
			cd := &html.Node{}
			cd.Type = html.TextNode
			cd.Data = t
			n.AppendChild(cd)
		case *html.Node:
			n.AppendChild(t)
		case *goxml.Element:
			xmltext, err := xd.parseHTMLText(t.ToXML())
			if err != nil {
				return nil, err
			}
			elt := xmltext.FirstChild.LastChild.FirstChild
			elt.Parent = nil
			n.AppendChild(elt)
		case *goxml.Attribute:
			cd := &html.Node{}
			cd.Type = html.TextNode
			cd.Data = t.Stringvalue()
			n.AppendChild(cd)
		case float64:
			cd := &html.Node{}
			cd.Type = html.TextNode
			cd.Data = strconv.FormatFloat(t, 'f', -1, 64)
			n.AppendChild(cd)
		case int:
			cd := &html.Node{}
			cd.Type = html.TextNode
			cd.Data = fmt.Sprintf("%d", t)
			n.AppendChild(cd)
		default:
			slog.Error(fmt.Sprintf("%s (line %d): unknown type %T (getTextvalues)", cmdname, line, t))
		}
	}
	return n, nil
}

func getStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

var (
	dummyBool           bool
	dummyStr            string
	dummyInt            int
	dummySP             bag.ScaledPoint
	boolType            = reflect.TypeOf(true)
	boolPtrType         = reflect.TypeOf(&dummyBool)
	stringType          = reflect.TypeOf("")
	stringPtrType       = reflect.TypeOf(&dummyStr)
	intType             = reflect.TypeOf(0)
	intPtrType          = reflect.TypeOf(&dummyInt)
	scaledPointsType    = reflect.TypeOf(dummySP)
	scaledPointsPtrType = reflect.TypeOf(&dummySP)
)

// getXMLAttributes fills the struct at v with the attribute values of the
// current element.
func getXMLAttributes(xd *xtsDocument, layoutelt *goxml.Element, v any) error {
	attributes := make(map[string]string)

	// Activate this code to get attributes from <Attributes><Attribute>
	// elements:
	//  for _, v := range layoutelt.Children() {
	//  if elt, ok := v.(*goxml.Element); ok {
	//      if elt.Name == "Attributes" {
	//          seq, err := dispatch(xd, elt, xd.data)
	//          if err != nil {
	//              return err
	//          }
	//          for _, itm := range seq {
	//              if attr, ok := itm.(goxml.Attribute); ok {
	//                  attributes[attr.Name] = attr.Value
	//              }
	//          }
	//      }
	//  }
	// }

	for _, attrib := range layoutelt.Attributes() {
		name := strings.ReplaceAll(attrib.Name, "-", "")
		attributes[name] = attrib.Value
	}

	val := reflect.ValueOf(v)

	// If it's an interface or a pointer, unwrap it.
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		return fmt.Errorf("s must be a struct")
	}

	var valNumFields = val.NumField()

	var mustexist bool
	var dflt string
	var allowXPath bool
	var attValue string

	for i := 0; i < valNumFields; i++ {
		mustexist = false
		allowXPath = true
		dflt = ""

		field := val.Field(i)
		structField := val.Type().Field(i)
		fieldName := strings.ToLower(structField.Name)
		for _, tag := range strings.Split(getStructTag(structField, "sdxml"), ",") {
			if suffix, ok := strings.CutPrefix(tag, "default:"); ok {
				dflt = suffix
			} else if suffix, ok := strings.CutPrefix(tag, "attr:"); ok {
				fieldName = suffix
			} else if tag == "mustexist" {
				mustexist = true
			} else if tag == "noescape" {
				allowXPath = false
			}
		}
		hasAttribute := false
		if a, ok := attributes[fieldName]; ok {
			hasAttribute = true
			if allowXPath {
				attValue = attributeValueRE.ReplaceAllStringFunc(a, func(a string) string {
					// strip curly braces
					seq, err := evaluateXPath(xd, layoutelt.Namespaces, a[1:len(a)-1])
					if err != nil {
						slog.Error(fmt.Sprintf("Layout line %d: %s", layoutelt.Line, err))
						return ""
					}
					return seq.Stringvalue()
				})

			} else {
				attValue = a
			}
		} else {
			if mustexist {
				slog.Error(fmt.Sprintf("Layout line %d: attribute %s on element %s not) found", layoutelt.Line, fieldName, layoutelt.Name))
				return fmt.Errorf("line %d: attribute %s on element %s not found", layoutelt.Line, fieldName, layoutelt.Name)
			}
			if dflt != "" {
				attValue = dflt
				hasAttribute = true
			}
		}
		if hasAttribute {
			switch field.Type() {
			case intType:
				attInt, ok := getInt(attValue)
				if !ok {
					return fmt.Errorf("Could not get int from %s", attValue)
				}
				field.SetInt(int64(attInt))
			case intPtrType:
				attInt, err := strconv.Atoi(attValue)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(&attInt))
			case stringType:
				field.SetString(attValue)
			case stringPtrType:
				a := attValue
				field.Set(reflect.ValueOf(&a))
			case boolPtrType:
				b := attValue == "yes"
				field.Set(reflect.ValueOf(&b))
			case boolType:
				field.SetBool(attValue == "yes")
			case scaledPointsType:
				var wd bag.ScaledPoint
				if cols, err := strconv.Atoi(attValue); err == nil {
					if strings.Contains(fieldName, "width") || strings.HasSuffix(fieldName, "x") {
						wd = xd.currentGrid.width(coord(cols))
					} else if strings.Contains(fieldName, "height") || strings.HasSuffix(fieldName, "y") {
						wd = xd.currentGrid.height(coord(cols))
					}
				} else {
					wd, err = bag.Sp(attValue)
					if err != nil {
						return err
					}
				}
				field.Set(reflect.ValueOf(wd).Convert(scaledPointsType))
			case scaledPointsPtrType:
				var wd bag.ScaledPoint
				if cols, err := strconv.Atoi(attValue); err == nil {
					if strings.Contains(fieldName, "width") || strings.HasSuffix(fieldName, "x") {
						wd = xd.currentGrid.width(coord(cols))
					} else if strings.Contains(fieldName, "height") || strings.HasSuffix(fieldName, "y") {
						wd = xd.currentGrid.height(coord(cols))
					}
				} else {
					wd, err = bag.Sp(attValue)
					if err != nil {
						return err
					}
				}
				field.Set(reflect.ValueOf(&wd).Convert(scaledPointsPtrType))
			}
		}
	}
	return nil
}

func findAttribute(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string, xp *xpath.Parser) (string, error) {
	var value string
	var found bool
	for _, attrib := range element.Attributes() {
		if attrib.Name == name {
			found = true
			value = attrib.Value
			break
		}
	}
	if !found {
		if mustexist {
			slog.Error(fmt.Sprintf("Layout line %d: attribute %s on element %s not) found", element.Line, name, element.Name))
			return "", fmt.Errorf("line %d: attribute %s on element %s not found", element.Line, name, element.Name)
		}
		value = dflt
	}

	value = attributeValueRE.ReplaceAllStringFunc(value, func(a string) string {
		// strip curly braces
		seq, err := xp.Evaluate(a[1 : len(a)-1])
		if err != nil {
			slog.Error(fmt.Sprintf("Layout line %d: %s", element.Line, err))
			return ""
		}
		return seq.Stringvalue()
	})
	return value, nil
}

func (xd *xtsDocument) getAttributeBool(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (bool, error) {
	attr, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if attr == "yes" {
		return true, err
	}

	return false, err
}

func (xd *xtsDocument) getAttributeString(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (string, error) {
	return findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
}

func (xd *xtsDocument) getAttributeInt(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (int, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// getAttributeSize returns the provided width in scaled points.
func (xd *xtsDocument) getAttributeSize(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (bag.ScaledPoint, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	return bag.Sp(val)
}

// getAttributeWidth returns the width which ich provided either by grid cells or a length value.
func (xd *xtsDocument) getAttributeWidth(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (bag.ScaledPoint, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	if cols, err := strconv.Atoi(val); err == nil {
		return xd.currentGrid.width(coord(cols)), nil

	}
	return bag.Sp(val)
}

// getAttributeHeight returns the width which ich provided either by grid cells or a length value.
func (xd *xtsDocument) getAttributeHeight(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (bag.ScaledPoint, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	if cols, err := strconv.Atoi(val); err == nil {
		return xd.currentGrid.height(coord(cols)), nil

	}
	return bag.Sp(val)
}

// evaluateXPath runs an XPath expression. It saves and restores the current
// context.
func evaluateXPath(xd *xtsDocument, namespaces map[string]string, xpath string) (xpath.Sequence, error) {
	oldContext := xd.data.Ctx.GetContextSequence()
	xd.data.Ctx.Namespaces = namespaces
	seq, err := xd.data.Evaluate(xpath)
	for _, itm := range seq {
		if f, ok := itm.(seqfunc); ok {
			// we assume that f() re-sets the old context
			seq, err = f()
			return seq, err
		}
	}

	xd.data.Ctx.SetContextSequence(oldContext)
	return seq, err
}

func getInt(in string) (int, bool) {
	f, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0, false
	}
	return int(f), true
}

func getFourValues(str string) map[string]string {
	fields := strings.Fields(str)
	fourvalues := make(map[string]string)
	switch len(fields) {
	case 1:
		fourvalues["top"] = fields[0]
		fourvalues["bottom"] = fields[0]
		fourvalues["left"] = fields[0]
		fourvalues["right"] = fields[0]
	case 2:
		fourvalues["top"] = fields[0]
		fourvalues["bottom"] = fields[0]
		fourvalues["left"] = fields[1]
		fourvalues["right"] = fields[1]
	case 3:
		fourvalues["top"] = fields[0]
		fourvalues["left"] = fields[1]
		fourvalues["right"] = fields[1]
		fourvalues["bottom"] = fields[2]
	case 4:
		fourvalues["top"] = fields[0]
		fourvalues["right"] = fields[1]
		fourvalues["bottom"] = fields[2]
		fourvalues["left"] = fields[3]
	}

	return fourvalues
}

func getFourValuesSP(str string) (map[string]bag.ScaledPoint, error) {
	var err error
	fields := strings.Fields(str)
	fieldsSP := make([]bag.ScaledPoint, len(fields))
	for i, f := range fields {
		if fieldsSP[i], err = bag.Sp(f); err != nil {
			return nil, err
		}
	}

	fourvalues := make(map[string]bag.ScaledPoint)
	switch len(fields) {
	case 1:
		fourvalues["top"] = fieldsSP[0]
		fourvalues["bottom"] = fieldsSP[0]
		fourvalues["left"] = fieldsSP[0]
		fourvalues["right"] = fieldsSP[0]
	case 2:
		fourvalues["top"] = fieldsSP[0]
		fourvalues["bottom"] = fieldsSP[0]
		fourvalues["left"] = fieldsSP[1]
		fourvalues["right"] = fieldsSP[1]
	case 3:
		fourvalues["top"] = fieldsSP[0]
		fourvalues["left"] = fieldsSP[1]
		fourvalues["right"] = fieldsSP[1]
		fourvalues["bottom"] = fieldsSP[2]
	case 4:
		fourvalues["top"] = fieldsSP[0]
		fourvalues["right"] = fieldsSP[1]
		fourvalues["bottom"] = fieldsSP[2]
		fourvalues["left"] = fieldsSP[3]
	}
	return fourvalues, nil
}
