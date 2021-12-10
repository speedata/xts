package core

import (
	"fmt"
	"strconv"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath/xpath"
)

// Get the values from the child elements of B, Paragraph and its ilk and fill
// the provided typesetting element to get a recursive data structure.
func getTextvalues(te *document.TypesettingElement, seq xpath.Sequence, cmdname string) {
	for _, itm := range seq {
		switch t := itm.(type) {
		case *goxml.Element:
			te.Items = append(te.Items, t.Stringvalue())
		case float64:
			te.Items = append(te.Items, strconv.FormatFloat(t, 'f', -1, 64))
		case goxml.CharData:
			te.Items = append(te.Items, string(t))
		case string:
			te.Items = append(te.Items, t)
		case *document.TypesettingElement:
			te.Items = append(te.Items, t)
		default:
			bag.Logger.DPanicf("%s: unknown type %T", cmdname, t)
		}
	}
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
			bag.Logger.Errorf("Layout line %d: attribute %s on element %s not found", element.Line, name, element.Name)
			return "", fmt.Errorf("line %d: attribute %s on element %s not found", element.Line, name, element.Name)
		}
		value = dflt
	}

	value = attributeValueRE.ReplaceAllStringFunc(value, func(a string) string {
		// strip curly braces
		seq, err := xp.Evaluate(a[1 : len(a)-1])
		if err != nil {
			bag.Logger.Errorf("Layout line %d: %s", element.Line, err)
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
	if wd, err := strconv.Atoi(val); err == nil {
		return xd.currentGrid.width(wd), nil

	}
	return bag.Sp(val)
}
